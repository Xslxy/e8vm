package ir

import (
	"e8vm.io/e8vm/link8"
)

// BuildPkg builds a package and returns the built lib
func BuildPkg(p *Pkg) *link8.Pkg {
	p.strPool.declare(p.lib)

	for _, v := range p.vars {
		var align uint32 = regSize
		if v.size <= 1 {
			align = 1
		}
		obj := link8.NewVar(align)
		obj.Pad(uint32(v.size))
		p.lib.DefineVar(v.sym, obj)
	}

	for _, f := range p.funcs {
		genFunc(p.g, f)
		writeFunc(p, f)
	}

	if p.tests != nil {
		v := link8.NewVar(regSize)
		for _, f := range p.tests.funcs {
			if err := v.WriteLink(0, f.index); err != nil {
				panic(err)
			}
		}
		p.lib.DefineVar(p.tests.sym, v)
	}

	return p.lib
}