package main

import (
	"fmt"
	"log"
)

func cons(car, cdr Obj) Obj {
	return &Cell{car, cdr}
}

func acons(x, y, a Obj) Obj {
	return cons(cons(x, y), a)
}

func isNil(obj Obj) bool {
	return obj == Nil
}

func isDot(obj Obj) bool {
	return obj == Dot
}

func isCparen(obj Obj) bool {
	return obj == Cparen
}

func isInt(obj Obj) bool {
	_, ok := obj.(Int)
	return ok
}

func isCell(obj Obj) bool {
	_, ok := obj.(*Cell)
	return ok
}

func isMacro(obj Obj) bool {
	_, ok := obj.(Macro)
	return ok
}

func isFunc(obj Obj) bool {
	_, ok := obj.(Func)
	return ok
}

func isPrimitive(obj Obj) bool {
	_, ok := obj.(Primitive)
	return ok
}

func isList(obj Obj) bool {
	return isNil(obj) || isCell(obj)
}

func isSymbol(obj Obj) bool {
	_, ok := obj.(Symbol)
	return ok
}

func listLength(list Obj) int {
	length := 0
	for !isNil(list) {
		if !isCell(list) {
			log.Fatalln("length: cannot handle dotted list")
		}
		list = list.(*Cell).Cdr
		length++
	}
	return length
}

func isAlpha(c rune) bool {
	return ('a' <= c && c <= 'z') ||
		('A' <= c && c <= 'Z')
}

func isDigit(c rune) bool {
	return '0' <= c && c <= '9'
}

func isAlnum(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func isPunct(c rune) bool {
	switch c {
	case '+', '=', '!', '@', '#', '$', '%', '^', '&', '*':
		return true
	}
	return false
}

func printObj(obj Obj) {
	switch obj := obj.(type) {
	case Int, Symbol:
		fmt.Print(obj)
	case Primitive:
		fmt.Print("<primitive>")
	case Func:
		fmt.Print("<function>")
	case Macro:
		fmt.Print("<macro>")
	case Special:
		switch obj {
		case Nil:
			fmt.Print("()")
		case True:
			fmt.Print("t")
		default:
			log.Fatalln("bug: print: unknown subtype:", obj)
		}
	case *Cell:
		fmt.Print("(")
		t := Obj(obj)
		for {
			p := t.(*Cell)
			printObj(p.Car)
			if isNil(p.Cdr) {
				break
			}
			if !isCell(p.Cdr) {
				fmt.Print(" . ")
				printObj(p.Cdr)
				break
			}
			fmt.Print(" ")
			t = p.Cdr
		}
		fmt.Print(")")
	default:
		log.Fatalf("bug: print: unknown tag: %T\n", obj)
	}
}
