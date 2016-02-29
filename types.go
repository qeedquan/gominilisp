package main

type Obj interface {
	obj()
}

type Int int

type Cell struct {
	Car, Cdr Obj
}

type Symbol string

type Primitive func(env, args Obj) Obj

type FuncOrMacro struct {
	Params Obj
	Body   Obj
	Env    Obj
}

type Func FuncOrMacro
type Macro FuncOrMacro

type Env struct {
	Vars Obj
	Up   Obj
}

type Special int

var (
	Nil    Special = 1
	Dot    Special = 2
	Cparen Special = 3
	True   Special = 4
)

func (Int) obj()       {}
func (*Cell) obj()     {}
func (Symbol) obj()    {}
func (Func) obj()      {}
func (Macro) obj()     {}
func (Primitive) obj() {}
func (Special) obj()   {}
func (*Env) obj()      {}
