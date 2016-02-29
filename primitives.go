package main

import (
	"fmt"
	"log"
	"os"
)

func primQuote(_, list Obj) Obj {
	if listLength(list) != 1 {
		log.Fatalln("malformed quote")
	}
	return list.(*Cell).Car
}

func primList(env, list Obj) Obj {
	return evalList(env, list)
}

func primSetq(env, list Obj) Obj {
	if listLength(list) != 2 || !isSymbol(list.(*Cell).Car) {
		log.Fatalln("malformed setq")
	}

	listCell := list.(*Cell)
	bind := find(env, listCell.Car)
	if bind == nil {
		log.Fatalln("unbound variable", listCell.Car.(Symbol))
	}
	value := eval(env, listCell.Cdr.(*Cell).Car)

	bindCell := bind.(*Cell)
	bindCell.Cdr = value
	return value
}

func primPlus(env, list Obj) Obj {
	sum := Int(0)
	for args := evalList(env, list); !isNil(args); args = args.(*Cell).Cdr {
		if !isInt(args.(*Cell).Car) {
			log.Fatalf("+ takes only numbers, got %T", args)
		}
		sum += args.(*Cell).Car.(Int)
	}
	return sum
}

func handleFunc(env, list Obj) FuncOrMacro {
	if !isCell(list) {
		log.Fatalln("malformed lambda")
	}

	listCell := list.(*Cell)
	if !isList(listCell.Car) || !isCell(listCell.Cdr) {
		log.Fatalln("malformed lambda")
	}

	for p := listCell.Car; !isNil(p); p = p.(*Cell).Cdr {
		if !isSymbol(p.(*Cell).Car) {
			log.Fatalln("parameter must be a symbol")
		}
		if !isList(p.(*Cell).Cdr) {
			log.Fatalln("parameter list is not a flat list")
		}
	}

	return FuncOrMacro{listCell.Car, listCell.Cdr, env}
}

func primLambda(env, list Obj) Obj {
	return Func(handleFunc(env, list))
}

func handleDefun(env, list Obj, kind int) Obj {
	car, cdr := list.(*Cell).Car, list.(*Cell).Cdr
	if !isSymbol(car) || !isCell(cdr) {
		log.Fatalln("malformed defun")
	}
	sym := car
	rest := cdr
	fn := handleFunc(env, rest)
	f, m := Func(fn), Macro(fn)
	switch kind {
	case 0:
		addVariable(env, sym, f)
		return f
	case 1:
		addVariable(env, sym, m)
		return m
	default:
		panic("unreachable")
	}
}

func primDefun(env, list Obj) Obj {
	return handleDefun(env, list, 0)
}

func primDefine(env, list Obj) Obj {
	if listLength(list) != 2 || !isSymbol(list.(*Cell).Car) {
		log.Fatalln("malformed setq")
	}

	sym := list.(*Cell).Car
	value := eval(env, list.(*Cell).Cdr.(*Cell).Car)
	addVariable(env, sym, value)
	return value
}

func primDefmacro(env, list Obj) Obj {
	return handleDefun(env, list, 1)
}

func primMacroExpand(env, list Obj) Obj {
	if listLength(list) != 1 {
		log.Fatalln("malformed macroexpand")
	}
	body := list.(*Cell).Car
	return macroExpand(env, body)
}

func primPrintln(env, list Obj) Obj {
	printObj(eval(env, list.(*Cell).Car))
	fmt.Println("")
	return Nil
}

func primIf(env, list Obj) Obj {
	if listLength(list) < 2 {
		log.Fatalln("malformed if")
	}
	cond := eval(env, list.(*Cell).Car)
	if !isNil(cond) {
		then := list.(*Cell).Cdr.(*Cell).Car
		return eval(env, then)
	}
	els := list.(*Cell).Cdr.(*Cell).Cdr
	if isNil(els) {
		return Nil
	}
	return progn(env, els)
}

func primNumEq(env, list Obj) Obj {
	if listLength(list) != 2 {
		log.Fatalln("malformed =")
	}
	values := evalList(env, list)
	x := values.(*Cell).Car
	y := values.(*Cell).Cdr.(*Cell).Car
	if !isInt(x) || !isInt(y) {
		log.Fatalf("= only takes numbers, got %T %T", x, y)
	}
	if x.(Int) == y.(Int) {
		return True
	}
	return Nil
}

func primExit(_, _ Obj) Obj {
	os.Exit(0)
	return nil
}
