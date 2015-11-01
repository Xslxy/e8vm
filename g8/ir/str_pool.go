package ir

import (
	"fmt"

	"e8vm.io/e8vm/link8"
)

type strConst struct {
	id  int
	str string
	sym uint32
}

func newStrConst(id int, s string) *strConst {
	ret := new(strConst)
	ret.id = id
	ret.str = s
	return ret
}

func (s *strConst) String() string {
	return fmt.Sprintf("<str %d>", s.id)
}

func (s *strConst) RegSizeAlign() bool { return true }

// A string constant is basically a byte slice.  It contains two register
// size fields: a pointer to the start of the string, and the size of the
// string.
func (s *strConst) Size() int32 {
	return regSize * 2
}

type strPool struct {
	strs   []*strConst
	strMap map[string]*strConst
}

func newStrPool() *strPool {
	ret := new(strPool)
	ret.strMap = make(map[string]*strConst)
	return ret
}

func (p *strPool) addString(s string) *strConst {
	exist := p.strMap[s]
	if exist != nil {
		return exist
	}

	n := len(p.strs)
	ret := newStrConst(n, s)
	p.strs = append(p.strs, ret)
	p.strMap[s] = ret

	return ret
}

func countDigit(n int) int {
	ret := 1
	for n > 9 {
		n /= 10
		ret++
	}
	return ret
}

func (p *strPool) declare(lib *link8.Pkg) {
	ndigit := countDigit(len(p.strs))
	nfmt := fmt.Sprintf(":str_%%0%dd", ndigit)

	for i, s := range p.strs {
		name := fmt.Sprintf(nfmt, i)
		v := link8.NewVar(0)
		v.Write([]byte(s.str))

		s.sym = lib.DeclareVar(name)
		lib.DefineVar(s.sym, v)
	}
}