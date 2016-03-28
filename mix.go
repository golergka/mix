package mix

import "fmt"
import _ "encoding/binary"

type Sign bool

func (s* Sign) String() string {
	if *s {
		return "+"
	} else {
		return "-"
	}
}

type Word struct {
	Bytes	[5]byte // 0 is MSB, 4 is LSB
	Sign	Sign
}

func (w* Word) String() string {
	return fmt.Sprintf("%v %02d %02d %02d %02d %02d", 
		w.Sign.String(), 
		w.Bytes[0], 
		w.Bytes[1],
		w.Bytes[2],
		w.Bytes[3],
		w.Bytes[4])
}

func (w* Word) Field(mod byte) Word {
	r := Word{Sign:true}
	l := mod / 8
	h := mod % 8
	if l == 0 {
		r.Sign = w.Sign
		l++
	}
	o := 5 - h
	for i := l - 1; i < h; i++ {
		r.Bytes[i + o] = w.Bytes[i]
	}
	return r
}

type Index struct {
	Bytes	[2]byte
	Sign	Sign
}

type Jump struct {
	Bytes	[2]byte
}

type Comparison int
const (
	Less Comparison = iota
	Equal
	Greater
)

func Adr(raw []byte) int16 {
	r := int16(0)
	l := len(raw)
	for i, v := range raw {
		a := int16(v) << byte((l - i - 1) * 6)
		r += int16(a)
	}
	return r
}

func SignedAdr(s* Sign, raw []byte) int16 {
	r := Adr(raw)
	if !*s {
		r = -r
	}
	return r
}

func (w *Word) SignedAdr() int16 {
	return SignedAdr(&w.Sign, w.Bytes[0:2])
}

func (i *Index) SignedAdr() int16 {
	return SignedAdr(&i.Sign, i.Bytes[:])
}

type Registers struct {
	RA	Word
	RX	Word
	RI	[6]Index

	O	bool
	C	Comparison
}

type Mix struct {
	Registers
	Memory [4000]Word
}

func (m *Mix) EffectiveAdr(i *Word) int16 {
	return m.RI[i.Bytes[2] - 1].SignedAdr() + i.SignedAdr()
}

type Op byte

const (
	OP_LDA Op = 8
	OP_LD1 Op = 9
	OP_LD2 Op = 10
	OP_LD3 Op = 11
	OP_LD4 Op = 12
	OP_LD5 Op = 13
	OP_LD6 Op = 14
	OP_LDX Op = 15
)

func (w *Word) Opcode() Op {
	return Op(w.Bytes[4])
}

func (m *Mix) Do(i *Word) {
	switch i.Opcode() {
	case OP_LDA:
		a := m.EffectiveAdr(i)
		m.RA = m.Memory[a].Field(i.Bytes[3])
	}
}
