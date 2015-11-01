package ir

import (
	"fmt"
)

// Func is an IR function. It consists of a bunch of named
// or unamed local variables and also a set of basic blocks.
// it can generate a linkable function.
type Func struct {
	name string

	sig *FuncSig

	savedRegs []*varRef
	locals    []*varRef // local variables
	retAddr   *varRef   // saved return address register

	prologue *Block
	epilogue *Block
	nblock   int

	nvar      int
	frameSize int32

	index    uint32 // the index in the lib
	isMain   bool
	isMethod bool
}

func newFunc(name string, sig *FuncSig) *Func {
	ret := new(Func)
	ret.name = name
	ret.sig = sig

	ret.prologue = ret.newBlock(nil)
	ret.epilogue = ret.newBlock(ret.prologue)

	return ret
}

func newMethod(name string, sig *FuncSig) *Func {
	ret := newFunc(name, sig)
	ret.isMethod = true
	return ret
}

// ThisRef returns the ref of the first arg.
func (f *Func) ThisRef() Ref {
	return f.sig.args[0]
}

// ArgRefs returns the refs to the arguments of the function
func (f *Func) ArgRefs() []Ref {
	var ret []Ref
	for _, arg := range f.sig.args {
		ret = append(ret, arg)
	}
	return ret
}

// RetRefs returns the refs to the return values of the function
func (f *Func) RetRefs() []Ref {
	var ret []Ref
	for _, arg := range f.sig.rets {
		ret = append(ret, arg)
	}
	return ret
}

// NewLocal creates a new named local variable of size n on the stack.
func (f *Func) NewLocal(n int32, name string, u8, regSizeAlign bool) Ref {
	ret := newVar(n, name, u8, regSizeAlign)
	f.locals = append(f.locals, ret)
	return ret
}

func (f *Func) newTempName() string {
	ret := fmt.Sprintf("<%d>", f.nvar)
	f.nvar++
	return ret
}

// NewTemp creates a new temp variable of size n on the stack.
func (f *Func) NewTemp(n int32, u8, regSizeAlign bool) Ref {
	return f.NewLocal(n, f.newTempName(), u8, regSizeAlign)
}

func (f *Func) newBlock(after *Block) *Block {
	ret := new(Block)
	ret.id = f.nblock
	ret.frameSize = &f.frameSize

	f.nblock++

	if after != nil {
		ret.next = after.next
		ret.jump = after.jump

		after.next = ret
		after.jump = nil // jump to natural next, which is ret
	}

	return ret
}

// End returns the ending block of the function (the epilogue).
func (f *Func) End() *Block { return f.epilogue }

// NewBlock creates a new basic block for the function
func (f *Func) NewBlock(after *Block) *Block {
	if after == nil {
		after = f.prologue
	}
	ret := f.newBlock(after)
	return ret
}

// SetAsMain marks the function as main function
// and this function will have a bare metal prologue and epilogue.
func (f *Func) SetAsMain() { f.isMain = true }

func (f *Func) String() string { return f.name }

// Size returns the size of a function pointer.
func (f *Func) Size() int32 { return regSize }

// RegSizeAlign returns true. A function pointer is always word aligned.
func (f *Func) RegSizeAlign() bool { return true }

// Index returns the symbol index of this function in the package.
func (f *Func) Index() uint32 { return f.index }

// Import returns the function symbol for the function where the package
// is imported as pindex.
func (f *Func) Import(pindex uint32) *FuncSym {
	return &FuncSym{name: f.name, pkg: pindex, sym: f.index, sig: f.sig}
}