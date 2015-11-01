package ir

import (
	"fmt"

	"e8vm.io/e8vm/link8"
)

// Pkg is a package in its intermediate representation.
type Pkg struct {
	lib *link8.Pkg

	path string

	funcs   []*Func
	vars    []*HeapSym
	tests   *testList
	strPool *strPool

	// helper functions required for generating
	g *gener
}

// NewPkg creates a package with a particular path name.
func NewPkg(path string) *Pkg {
	ret := new(Pkg)
	ret.path = path
	ret.lib = link8.NewPkg(path)
	ret.strPool = newStrPool()

	ret.g = new(gener)

	return ret
}

// NewFunc creates a new function for the package.
func (p *Pkg) NewFunc(name string, sig *FuncSig) *Func {
	ret := newFunc(name, sig)
	ret.index = p.lib.DeclareFunc(ret.name)
	p.funcs = append(p.funcs, ret)
	return ret
}

// NewMethod creates a new method function for the package.
func (p *Pkg) NewMethod(name string, sig *FuncSig) *Func {
	ret := newMethod(name, sig)
	ret.index = p.lib.DeclareFunc(ret.name)
	p.funcs = append(p.funcs, ret)
	return ret
}

// NewGlobalVar creates a new global variable reference.
func (p *Pkg) NewGlobalVar(
	size int32, name string, u8, regSizeAlign bool,
) Ref {
	ret := newHeapSym(size, name, u8, regSizeAlign)
	ret.sym = p.lib.DeclareVar(ret.name)
	p.vars = append(p.vars, ret)
	return ret
}

// NewTestList creates a global variable of a list of function symbols.
func (p *Pkg) NewTestList(name string, funcs []*Func) Ref {
	if len(funcs) > 1000000 {
		panic("too many test cases")
	}
	if p.tests != nil {
		panic("tests already built")
	}

	ret := newTestList(name, funcs)
	ret.sym = p.lib.DeclareVar(ret.name)
	p.tests = ret

	return ret
}

// Require imports a linkable package.
func (p *Pkg) Require(pkg *link8.Pkg) uint32 { return p.lib.Require(pkg) }

// RequireBuiltin imports the builtin package that provides neccessary
// builtin functions.
func (p *Pkg) RequireBuiltin(pkg *link8.Pkg) (uint32, error) {
	ret := p.Require(pkg)

	var err error
	se := func(e error) {
		if err != nil {
			err = e
		}
	}

	o := func(f string) *FuncSym {
		sym, index := pkg.SymbolByName(f)
		if sym == nil {
			se(fmt.Errorf("%s missing in builtin", f))
		} else if sym.Type != link8.SymFunc {
			se(fmt.Errorf("%s is not a function", f))
		}

		return &FuncSym{pkg: ret, sym: index}
	}

	p.g.memCopy = o("MemCopy")
	p.g.memClear = o("MemClear")

	return ret, err
}

// NewString adds a new string constant and returns its reference.
func (p *Pkg) NewString(s string) Ref {
	return p.strPool.addString(s)
}