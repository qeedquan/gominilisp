package main

import "log"

func addVariable(env, sym, val Obj) {
	e := env.(*Env)
	e.Vars = acons(sym, val, e.Vars)
}

func pushEnv(env, vars, values Obj) Obj {
	if listLength(vars) != listLength(values) {
		log.Fatalln("cannot apply function: number of argument does not match")
	}

	m := Obj(Nil)
	p, q := vars, values
	for !isNil(p) {
		pc, qc := p.(*Cell), q.(*Cell)
		sym := pc.Car
		val := qc.Car
		m = acons(sym, val, m)

		p, q = pc.Cdr, qc.Cdr
	}
	return &Env{m, env}
}

func progn(env, list Obj) Obj {
	var r Obj
	for lp := list; !isNil(lp); lp = lp.(*Cell).Cdr {
		r = eval(env, lp.(*Cell).Car)
	}
	return r
}

func evalList(env, list Obj) Obj {
	var head, tail Obj
	for lp := list; !isNil(lp); lp = lp.(*Cell).Cdr {
		tmp := eval(env, lp.(*Cell).Car)
		if head == nil {
			head = cons(tmp, Nil)
			tail = head
		} else {
			tc := tail.(*Cell)
			tc.Cdr = cons(tmp, Nil)
			tail = tc.Cdr
		}
	}
	if head == nil {
		return Nil
	}
	return head
}

func apply(env, fn, args Obj) Obj {
	if !isList(args) {
		log.Fatalln("argument must be a list")
	}

	switch fn := fn.(type) {
	case Primitive:
		return fn(env, args)
	case Func:
		body := fn.Body
		params := fn.Params
		eargs := evalList(env, args)
		newenv := pushEnv(fn.Env, params, eargs)
		return progn(newenv, body)
	}
	log.Fatalln("not supported")
	return nil
}

func find(env, sym Obj) Obj {
	for p := env; p != nil; p = p.(*Env).Up {
		for c := p.(*Env).Vars; !isNil(c); c = c.(*Cell).Cdr {
			bind := c.(*Cell).Car
			if sym == bind.(*Cell).Car {
				return bind
			}
		}
	}
	return nil
}

func macroExpand(env, obj Obj) Obj {
	if !isCell(obj) {
		return obj
	}

	cell := obj.(*Cell)
	if !isSymbol(cell.Car) {
		return obj
	}

	bind := find(env, cell.Car)
	if bind == nil || !isMacro(bind.(*Cell).Cdr) {
		return obj
	}

	args := cell.Cdr
	macro := bind.(*Cell).Cdr.(Macro)
	body := macro.Body
	params := macro.Params
	newenv := pushEnv(env, params, args)
	return progn(newenv, body)
}

func eval(env, obj Obj) Obj {
	switch obj := obj.(type) {
	case Int, Primitive, Func, Special:
		return obj
	case Symbol:
		bind := find(env, obj)
		if bind == nil {
			log.Fatalln("undefined symbol:", obj)
		}
		return bind.(*Cell).Cdr
	case *Cell:
		expanded := macroExpand(env, obj)
		if expanded != obj {
			return eval(env, expanded)
		}
		fn := eval(env, obj.Car)
		args := obj.Cdr
		if !isPrimitive(fn) && !isFunc(fn) {
			log.Fatalf("the head of a list must be a function, got: %T\n", fn)
		}
		return apply(env, fn, args)
	}
	log.Fatalf("bug: eval: unknown tag type: %T\n", obj)
	return nil
}
