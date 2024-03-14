package main

import (
	"io"
	"strconv"

	"github.com/k0l1br1/converter/bins"
)

const (
	printLimit  = 36
	defaultPrec = 4
)

var (
	redC        = []byte{0x1b, '[', '3', '1', 'm'}
	greenC      = []byte{0x1b, '[', '3', '2', 'm'}
	resetC      = []byte{0x1b, '[', '0', 'm'}
	cursorUp    = []byte{0x1b, '[', '3', '6', 'A'}
	eraseScreen = []byte{0x1b, '[', '2', 'J'}
	eraseUp     = []byte{0x1b, '[', '1', 'J'}
	eraseDown   = []byte{0x1b, '[', 'J'}
)

type Printer struct {
	out io.Writer
	buf []byte
	ts  []float64
	vs  []float64
}

func NewPrinter(out io.Writer) *Printer {
	return &Printer{
		out: out,
		buf: make([]byte, 0, 1<<12),
		ts:  make([]float64, printLimit),
		vs:  make([]float64, printLimit),
	}
}

func (p *Printer) Flush() {
	p.out.Write(p.buf)
	p.buf = p.buf[:0]
}

func (p *Printer) BreakLine() {
	p.buf = append(p.buf, '\r', '\n')
	p.Flush()
}

func (p *Printer) EraseScreen() {
	p.buf = append(p.buf, cursorUp...)
	p.buf = append(p.buf, eraseDown...)
}

func (p *Printer) Reprint(bs []bins.Bin, off int) {
    p.EraseScreen()
    p.Print(bs, off)
}

func (p *Printer) Print(bs []bins.Bin, off int) {
	tmpBs := bs[off : printLimit+off]
	split(tmpBs, p.ts, p.vs)
	NormalizeZ(p.ts)
	NormalizeZ(p.vs)
	for i, v := range tmpBs {
		if v.IsUp == 1 {
			p.buf = append(p.buf, greenC...)
		} else {
			p.buf = append(p.buf, redC...)
		}
		if p.ts[i] > 0 {
			p.buf = append(p.buf, 0x20) // add space for non negative
		}
		p.buf = strconv.AppendFloat(p.buf, p.ts[i], 'f', defaultPrec, 64)

		p.buf = append(p.buf, 0x20) // space delim

		if p.vs[i] > 0 {
			p.buf = append(p.buf, 0x20)
		}
		p.buf = strconv.AppendFloat(p.buf, p.vs[i], 'f', defaultPrec, 64)
		p.buf = append(p.buf, resetC...)
		p.buf = append(p.buf, '\r', '\n')
	}
	p.Flush()
}
