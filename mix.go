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

func (w* Word) GetField(mod byte) Word {
	r := Word{Sign:true}
	r.SetField(mod, w)
	return r
}

func (w* Word) SetField(mod byte, v *Word) {
	l := mod / 8
	h := mod % 8
	if l == 0 {
		w.Sign = v.Sign
		l++
	}
	o := 5 - h
	for i := l - 1; i < h; i ++ {
		w.Bytes[i + o] = v.Bytes[i]
	}
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

func (m *Mix) GetMem(i *Word) Word {
	a := m.EffectiveAdr(i)
	return m.Memory[a].GetField(i.Bytes[3])
}

func (i *Index) FromWord(w *Word) {
	i.Sign = w.Sign
	copy(i.Bytes[:], w.Bytes[3:5])
}

func (i *Index) ToWord() Word {
	r := Word{Sign:i.Sign}
	copy(r.Bytes[3:5], i.Bytes[:])
	return r
}

type Op byte

const (
	OP_LDA	Op = 8
	OP_LD1 	Op = 9
	OP_LD2 	Op = 10
	OP_LD3 	Op = 11
	OP_LD4 	Op = 12
	OP_LD5 	Op = 13
	OP_LD6 	Op = 14
	OP_LDX 	Op = 15

	OP_LDAN Op = 16
	OP_LD1N	Op = 17
	OP_LD2N	Op = 18
	OP_LD3N	Op = 19
	OP_LD4N	Op = 20
	OP_LD5N	Op = 21
	OP_LD6N	Op = 22
	OP_LDXN Op = 23

	OP_STA  Op = 24
	OP_ST1 	Op = 25
	OP_ST2 	Op = 26
	OP_ST3 	Op = 27
	OP_ST4 	Op = 28
	OP_ST5 	Op = 29
	OP_ST6 	Op = 30
	OP_STX 	Op = 31
)

func (w *Word) Opcode() Op {
	return Op(w.Bytes[4])
}

func (m *Mix) Do(i *Word) {
	switch  o := i.Opcode(); o {
	case OP_LDA:
		v := m.GetMem(i)
		m.RA = v
	case OP_LDX:
		v := m.GetMem(i)
		m.RX = v
		break
	case OP_LD1, OP_LD2, OP_LD3, OP_LD4, OP_LD5, OP_LD6:
		v := m.GetMem(i)
		m.RI[o - OP_LD1].FromWord(&v)
		break
	case OP_LDAN:
		v := m.GetMem(i)
		v.Sign = !v.Sign
		m.RA = v
	case OP_LDXN:
		v := m.GetMem(i)
		v.Sign = !v.Sign
		m.RX = v
		break
	case OP_LD1N, OP_LD2N, OP_LD3N, OP_LD4N, OP_LD5N, OP_LD6N:
		v := m.GetMem(i)
		v.Sign = !v.Sign
		m.RI[o - OP_LD1N].FromWord(&v)
		break
	case OP_STA:
		a := m.EffectiveAdr(i)
		m.Memory[a].SetField(i.Bytes[3], &m.RA)
		break
	case OP_ST1, OP_ST2, OP_ST3, OP_ST4, OP_ST5, OP_ST6:
		a := m.EffectiveAdr(i)
		w := m.RI[o - OP_ST1].ToWord()
		m.Memory[a].SetField(i.Bytes[3], &w)
		break
	case OP_STX:
		a := m.EffectiveAdr(i)
		m.Memory[a].SetField(i.Bytes[3], &m.RX)
		break
	}
}
