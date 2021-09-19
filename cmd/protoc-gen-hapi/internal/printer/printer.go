package printer

import (
	"fmt"
	"strings"
)

type Printer struct {
	indent      string
	indentLevel int
	buf         strings.Builder
}

func New(indent string) *Printer {
	return &Printer{
		indent: indent,
	}
}

func (p *Printer) IncreaseIndentLevel() {
	p.indentLevel++
}

func (p *Printer) DecreaseIndentLevel() {
	if p.indentLevel <= 0 {
		panic(fmt.Errorf("negative indent count"))
	}
	p.indentLevel--
}

func (p *Printer) PrintNewLine() {
	p.buf.WriteByte('\n')
}

func (p *Printer) printIndents() {
	p.buf.WriteString(strings.Repeat(p.indent, p.indentLevel))
}

func (p *Printer) PrintlnWithIndents(s string, more ...string) {
	p.printIndents()
	p.buf.WriteString(s)
	for _, v := range more {
		p.buf.WriteString(v)
	}
	p.PrintNewLine()
}

func (p *Printer) Content() string {
	return p.buf.String()
}
