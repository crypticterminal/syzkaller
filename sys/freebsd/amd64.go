// AUTOGENERATED FILE
package freebsd

import . "github.com/google/syzkaller/prog"

func init() {
	RegisterTarget(&Target{OS: "freebsd", Arch: "amd64", Revision: revision_amd64, PtrSize: 8, Syscalls: syscalls_amd64, Resources: resources_amd64, Structs: structDescs_amd64, Consts: consts_amd64}, initTarget)
}

var resources_amd64 = []*ResourceDesc{
	{Name: "fd", Type: &IntType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "int32", TypeSize: 4}}}, Kind: []string{"fd"}, Values: []uint64{18446744073709551615}},
}

var structDescs_amd64 = []*KeyedStruct{
	{Key: StructKey{Name: "pipefd", Dir: 1}, Desc: &StructDesc{TypeCommon: TypeCommon{TypeName: "pipefd", TypeSize: 8, ArgDir: 1}, Fields: []Type{
		&ResourceType{TypeCommon: TypeCommon{TypeName: "fd", FldName: "rfd", TypeSize: 4, ArgDir: 1}},
		&ResourceType{TypeCommon: TypeCommon{TypeName: "fd", FldName: "wfd", TypeSize: 4, ArgDir: 1}},
	}}},
}

var syscalls_amd64 = []*Syscall{
	{NR: 477, Name: "mmap", CallName: "mmap", Args: []Type{
		&VmaType{TypeCommon: TypeCommon{TypeName: "vma", FldName: "addr", TypeSize: 8}},
		&LenType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "len", FldName: "len", TypeSize: 8}}, Buf: "addr"},
		&FlagsType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "mmap_prot", FldName: "prot", TypeSize: 8}}, Vals: []uint64{1, 2}},
		&FlagsType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "mmap_flags", FldName: "flags", TypeSize: 8}}, Vals: []uint64{2, 4096, 16}},
		&ConstType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "const", FldName: "fd", TypeSize: 8}}, Val: 18446744073709551615},
		&ConstType{IntTypeCommon: IntTypeCommon{TypeCommon: TypeCommon{TypeName: "const", FldName: "offset", TypeSize: 8}}},
	}, Ret: &VmaType{TypeCommon: TypeCommon{TypeName: "vma", FldName: "ret", TypeSize: 8, ArgDir: 1}}},
	{ID: 1, NR: 42, Name: "pipe", CallName: "pipe", Args: []Type{
		&PtrType{TypeCommon: TypeCommon{TypeName: "ptr", FldName: "pipefd", TypeSize: 8}, Type: &StructType{Key: StructKey{Name: "pipefd", Dir: 1}}},
	}},
}

var consts_amd64 = []ConstValue{
	{Name: "MAP_ANONYMOUS", Value: 4096},
	{Name: "MAP_FIXED", Value: 16},
	{Name: "MAP_PRIVATE", Value: 2},
	{Name: "PROT_READ", Value: 1},
	{Name: "PROT_WRITE", Value: 2},
	{Name: "SYS_mmap", Value: 477},
	{Name: "SYS_pipe", Value: 42},
}

const revision_amd64 = "7c737d486a33a6a0817ce924247b4b67428f7a07"