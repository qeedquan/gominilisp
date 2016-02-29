package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

type Interp struct {
	symbols Obj
	r       *bufio.Reader
}

func NewInterp() *Interp {
	return &Interp{}
}

func (ip *Interp) addPrimitive(env Obj, name string, fn Primitive) {
	sym := ip.intern(name)
	addVariable(env, sym, fn)
}

func (ip *Interp) defineConsts(env Obj) {
	sym := ip.intern("t")
	addVariable(env, sym, True)
}

func (ip *Interp) definePrimitives(env Obj) {
	prims := []struct {
		name      string
		primitive Primitive
	}{
		{"quote", primQuote},
		{"list", primList},
		{"setq", primSetq},
		{"+", primPlus},
		{"define", primDefine},
		{"defun", primDefun},
		{"defmacro", primDefmacro},
		{"macroexpand", primMacroExpand},
		{"lambda", primLambda},
		{"if", primIf},
		{"=", primNumEq},
		{"println", primPrintln},
		{"exit", primExit},
	}

	for _, p := range prims {
		ip.addPrimitive(env, p.name, p.primitive)
	}
}

func (ip *Interp) intern(name string) Obj {
	for p := ip.symbols; !isNil(p); p = p.(*Cell).Cdr {
		if p.(*Cell).Car.(Symbol) == Symbol(name) {
			return p.(*Cell).Car
		}
	}
	sym := Symbol(name)
	ip.symbols = cons(sym, ip.symbols)
	return sym
}

func (ip *Interp) Repl() {
	env := &Env{Nil, nil}
	ip.r = bufio.NewReader(os.Stdin)
	ip.symbols = Nil
	ip.defineConsts(env)
	ip.definePrimitives(env)

	for {
		expr := ip.read()
		if expr == nil {
			break
		}
		if isCparen(expr) {
			log.Fatalln("stray close parenthesis")
		} else if isDot(expr) {
			log.Fatalln(expr)
		}

		printObj(eval(env, expr))
		fmt.Println()
	}
}

const eof = -1

func (ip *Interp) getchar() rune {
	c, _, err := ip.r.ReadRune()
	if err == io.EOF {
		return eof
	}
	if err != nil {
		log.Fatalln(err)
	}
	return c
}

func (ip *Interp) peek() rune {
	c := ip.getchar()
	err := ip.r.UnreadRune()
	if err != nil {
		log.Fatalln(err)
	}
	return c
}

func (ip *Interp) read() Obj {
	for {
		c := ip.getchar()
		switch {
		case c == eof:
			return nil
		case c == ' ' || c == '\n' || c == '\r' || c == '\t':
			continue
		case c == ';':
			ip.skipLine()
			continue
		case c == '(':
			return ip.readList()
		case c == ')':
			return Cparen
		case c == '.':
			return Dot
		case c == '\'':
			return ip.readQuote()
		case isDigit(c):
			return Int(ip.readNumber(int(c - '0')))
		case c == '-':
			return Int(-ip.readNumber(0))
		case isAlpha(c) || isPunct(c):
			return ip.readSymbol(c)
		default:
			log.Fatalln("don't know how to handle:", string(c))
		}
	}
	return nil
}

func (ip *Interp) skipLine() {
	for {
		c := ip.getchar()
		switch c {
		case eof, '\n':
			return
		case '\r':
			if ip.peek() == '\n' {
				ip.getchar()
			}
			return
		}
	}
}

func (ip *Interp) readList() Obj {
	obj := ip.read()
	if obj == nil {
		log.Fatalln("unclosed parenthesis")
	}
	if isDot(obj) {
		log.Fatalln("stray dot")
	}
	if isCparen(obj) {
		return Nil
	}

	head := cons(obj, Nil)
	tail := head

	for {
		obj := ip.read()
		if obj == nil {
			log.Fatalln("unclosed parenthesis")
		}
		if isCparen(obj) {
			return head
		}
		if isDot(obj) {
			t := tail.(*Cell)
			t.Cdr = ip.read()
			if ip.read() != Cparen {
				log.Fatalln("unclosed parenthesis")
			}
			return head
		}

		t := tail.(*Cell)
		t.Cdr = cons(obj, Nil)
		tail = t.Cdr
	}
}

func (ip *Interp) readQuote() Obj {
	sym := ip.intern("quote")
	return cons(sym, cons(ip.read(), Nil))
}

func (ip *Interp) readNumber(val int) int {
	for isDigit(ip.peek()) {
		val = val*10 + (int(ip.getchar()) - '0')
	}
	return val
}

func (ip *Interp) readSymbol(c rune) Obj {
	buf := string(c)
	for isAlnum(ip.peek()) || ip.peek() == '-' {
		buf += string(ip.getchar())
	}
	return ip.intern(buf)
}
