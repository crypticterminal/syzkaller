// Copyright 2015 syzkaller project authors. All rights reserved.
// Use of this source code is governed by Apache 2 LICENSE that can be found in the LICENSE file.

// Conservative resource-related analysis of programs.
// The analysis figures out what files descriptors are [potentially] opened
// at a particular point in program, what pages are [potentially] mapped,
// what files were already referenced in calls, etc.

package prog

import (
	"fmt"
)

type state struct {
	target    *Target
	ct        *ChoiceTable
	files     map[string]bool
	resources map[string][]Arg
	strings   map[string]bool
	ma        *memAlloc
	va        *vmaAlloc
}

// analyze analyzes the program p up to but not including call c.
func analyze(ct *ChoiceTable, p *Prog, c *Call) *state {
	s := newState(p.Target, ct)
	for _, c1 := range p.Calls {
		if c1 == c {
			break
		}
		s.analyze(c1)
	}
	return s
}

func newState(target *Target, ct *ChoiceTable) *state {
	s := &state{
		target:    target,
		ct:        ct,
		files:     make(map[string]bool),
		resources: make(map[string][]Arg),
		strings:   make(map[string]bool),
		ma:        newMemAlloc(target.NumPages * target.PageSize),
		va:        newVmaAlloc(target.NumPages),
	}
	return s
}

func (s *state) analyze(c *Call) {
	ForeachArg(c, func(arg Arg, _ *ArgCtx) {
		switch a := arg.(type) {
		case *PointerArg:
			switch {
			case a.IsNull():
			case a.VmaSize != 0:
				s.va.noteAlloc(a.Address/s.target.PageSize, a.VmaSize/s.target.PageSize)
			default:
				s.ma.noteAlloc(a.Address, a.Res.Size())
			}
		}
		switch typ := arg.Type().(type) {
		case *ResourceType:
			if typ.Dir() != DirIn {
				s.resources[typ.Desc.Name] = append(s.resources[typ.Desc.Name], arg)
				// TODO: negative PIDs and add them as well (that's process groups).
			}
		case *BufferType:
			a := arg.(*DataArg)
			if typ.Dir() != DirOut && len(a.Data()) != 0 {
				switch typ.Kind {
				case BufferString:
					s.strings[string(a.Data())] = true
				case BufferFilename:
					s.files[string(a.Data())] = true
				}
			}
		}
	})
}

type ArgCtx struct {
	Parent *[]Arg      // GroupArg.Inner (for structs) or Call.Args containing this arg
	Base   *PointerArg // pointer to the base of the heap object containing this arg
	Offset uint64      // offset of this arg from the base
	Stop   bool        // if set by the callback, subargs of this arg are not visited
}

func ForeachSubArg(arg Arg, f func(Arg, *ArgCtx)) {
	foreachArgImpl(arg, ArgCtx{}, f)
}

func ForeachArg(c *Call, f func(Arg, *ArgCtx)) {
	ctx := ArgCtx{}
	if c.Ret != nil {
		foreachArgImpl(c.Ret, ctx, f)
	}
	ctx.Parent = &c.Args
	for _, arg := range c.Args {
		foreachArgImpl(arg, ctx, f)
	}
}

func foreachArgImpl(arg Arg, ctx ArgCtx, f func(Arg, *ArgCtx)) {
	f(arg, &ctx)
	if ctx.Stop {
		return
	}
	switch a := arg.(type) {
	case *GroupArg:
		if _, ok := a.Type().(*StructType); ok {
			ctx.Parent = &a.Inner
		}
		var totalSize uint64
		for _, arg1 := range a.Inner {
			foreachArgImpl(arg1, ctx, f)
			if !arg1.Type().BitfieldMiddle() {
				size := arg1.Size()
				ctx.Offset += size
				totalSize += size
			}
		}
		claimedSize := a.Size()
		varlen := a.Type().Varlen()
		if varlen && totalSize > claimedSize || !varlen && totalSize != claimedSize {
			panic(fmt.Sprintf("bad group arg size %v, should be <= %v for %#v type %#v",
				totalSize, claimedSize, a, a.Type()))
		}
	case *PointerArg:
		if a.Res != nil {
			ctx.Base = a
			ctx.Offset = 0
			foreachArgImpl(a.Res, ctx, f)
		}
	case *UnionArg:
		foreachArgImpl(a.Option, ctx, f)
	}
}

func RequiredFeatures(p *Prog) (bitmasks, csums bool) {
	for _, c := range p.Calls {
		ForeachArg(c, func(arg Arg, _ *ArgCtx) {
			if a, ok := arg.(*ConstArg); ok {
				if a.Type().BitfieldOffset() != 0 || a.Type().BitfieldLength() != 0 {
					bitmasks = true
				}
			}
			if _, ok := arg.Type().(*CsumType); ok {
				csums = true
			}
		})
	}
	return
}
