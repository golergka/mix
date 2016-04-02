package mix

func UnpackMod(mod byte) (byte, byte) {
	return (mod / 8), (mod % 8)
}

func PackMod(l byte, h byte) byte {
	return l*8 + h
}

type Comparison int

const (
	Less Comparison = iota
	Equal
	Greater
)

type Adr uint16
type SignAdr int16

func MakeAdr(raw []byte) Adr {
	r := Adr(0)
	l := len(raw)
	for i, v := range raw {
		a := int16(v) << byte((l-i-1)*6)
		r += Adr(a)
	}
	return r
}

func MakeSignAdr(s *Sign, raw []byte) SignAdr {
	r := SignAdr(MakeAdr(raw))
	if !*s {
		r = -r
	}
	return r
}

func (a SignAdr) ToBytes() (Sign, [2]byte) {
	var s Sign
	if a < 0 {
		s = false
		a = -a
	} else {
		s = true
	}
	var b [2]byte
	b[0] = byte(a / 64)
	b[1] = byte(a % 64)

	return s, b
}

type Registers struct {
	RA Word
	RX Word
	RI [6]Index

	O bool
	C Comparison
}

type Mix struct {
	Registers
	Memory [4000]Word
}

func (m *Mix) EffectiveAdr(w *Word) Adr {
	var indexDelta SignAdr
	i := w.GetIndex()
	if i > 0 {
		indexDelta = m.RI[i-1].GetAdr()
	}
	r := indexDelta + w.GetAdr()
	return Adr(r)
}

func (m *Mix) GetMem(i *Word) Word {
	a := m.EffectiveAdr(i)
	return m.Memory[a].GetField(i.GetMod())
}

func (m *Mix) Do(i *Word) {
	switch o := i.GetOp(); o {
	case LDA:
		v := m.GetMem(i)
		m.RA = v
	case LDX:
		v := m.GetMem(i)
		m.RX = v
		break
	case LD1, LD2, LD3, LD4, LD5, LD6:
		v := m.GetMem(i)
		m.RI[o-LD1].FromWord(&v)
		break
	case LDAN:
		v := m.GetMem(i)
		v.Sign = !v.Sign
		m.RA = v
	case LDXN:
		v := m.GetMem(i)
		v.Sign = !v.Sign
		m.RX = v
		break
	case LD1N, LD2N, LD3N, LD4N, LD5N, LD6N:
		v := m.GetMem(i)
		v.Sign = !v.Sign
		m.RI[o-LD1N].FromWord(&v)
		break
	case STA:
		a := m.EffectiveAdr(i)
		m.Memory[a].SetField(i.GetMod(), m.RA)
		break
	case ST1, ST2, ST3, ST4, ST5, ST6:
		a := m.EffectiveAdr(i)
		w := m.RI[o-ST1].ToWord()
		m.Memory[a].SetField(i.GetMod(), w)
		break
	case STX:
		a := m.EffectiveAdr(i)
		m.Memory[a].SetField(i.GetMod(), m.RX)
		break
	}
}
