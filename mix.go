package mix

import "fmt"
import "errors"

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

func (w Word) String() string {
	return fmt.Sprintf("%v %02d %02d %02d %02d %02d", 
		w.Sign.String(), 
		w.Bytes[0], 
		w.Bytes[1],
		w.Bytes[2],
		w.Bytes[3],
		w.Bytes[4])
}

func (w Word) Format(f fmt.State, c rune) {
	switch c {
	case 'i':
		o := w.GetOp()
		a := w.GetAdr()
		fmt.Fprintf(f, "%v %v", o, a)

		if i := w.GetIndex(); i != 0 {
			fmt.Fprintf(f, ",%v", i)
		}
		if l, h := UnpackMod(w.GetMod()); (l != 0) || (h != 0) {
			fmt.Fprintf(f, "(%v:%v)", l,h)
		}
		break
	default:
		fmt.Fprintf(f, w.String())
	}
}

func scanSign(state fmt.ScanState) (Sign, error) {
	sign, _, err := state.ReadRune()
	if err != nil {
		state.UnreadRune()
		return false, err
	}
	switch sign {
	case '+':
		return true, nil
	case '-':
		return false, nil
	default:
		state.UnreadRune()
		return false, errors.New(fmt.Sprintf("unexpected sign rune: ", sign))
	}
}

func (w* Word) Scan(state fmt.ScanState, verb rune) error {
	switch verb {
	case 'i': // OP ADDRESS,INDEX(MOD)
		// Scan OP & address
		{
			var o Op
			var a SignAdr
			if _, err := fmt.Fscanf(state, "%v %v", &o, &a); err != nil {
				return err
			}
			w.SetOp(o)
			w.SetAdr(a)
		}
		// Scan index (can be omitted)
		// Scan mod (can be omitted)
		return nil
	default:
		s, err := scanSign(state)
		w.Sign = s
		if err != nil {
			return err
		}
		if _, err := fmt.Fscanf(state, 
			"%v %v %v %v %v", 
			&w.Bytes[0], 
			&w.Bytes[1],
			&w.Bytes[2],
			&w.Bytes[3],
			&w.Bytes[4]); err != nil {
			return err
		}
		return nil
	}
}

func (w Word) GetField(mod byte) Word {
	r := Word{Sign:true}
	r.SetField(mod, w)
	return r
}

func (w Word) GetMod() byte {
	return w.Bytes[3]
}

func UnpackMod(mod byte) (byte, byte) {
	return (mod / 8), (mod % 8)
}

func (w* Word) SetField(mod byte, v Word) {
	l, h := UnpackMod(mod)
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

type Adr uint16
type SignAdr int16

func MakeAdr(raw []byte) Adr {
	r := Adr(0)
	l := len(raw)
	for i, v := range raw {
		a := int16(v) << byte((l - i - 1) * 6)
		r += Adr(a)
	}
	return r
}

func MakeSignAdr(s* Sign, raw []byte) SignAdr {
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

func (w Word) GetAdr() SignAdr {
	return MakeSignAdr(&w.Sign, w.Bytes[0:2])
}

func (w *Word) SetAdr(a SignAdr) {
	s, b := a.ToBytes()
	w.Sign = s
	copy(w.Bytes[0:2], b[:])
}

func (i *Index) GetAdr() SignAdr {
	return MakeSignAdr(&i.Sign, i.Bytes[:])
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

func (w Word) GetIndex() byte {
	return w.Bytes[2]
}

func (m *Mix) EffectiveAdr(i *Word) Adr {
	var indexDelta SignAdr
	if i.Bytes[2] > 0 {
		indexDelta = m.RI[i.Bytes[2] - 1].GetAdr()
	}
	r := indexDelta + i.GetAdr()
	return Adr(r)
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

//go:generate stringer -type=Op
const (
	LDA		Op = 8
	LD1 	Op = 9
	LD2 	Op = 10
	LD3 	Op = 11
	LD4 	Op = 12
	LD5 	Op = 13
	LD6 	Op = 14
	LDX 	Op = 15

	LDAN	Op = 16
	LD1N	Op = 17
	LD2N	Op = 18
	LD3N	Op = 19
	LD4N	Op = 20
	LD5N	Op = 21
	LD6N	Op = 22
	LDXN	Op = 23

	STA		Op = 24
	ST1 	Op = 25
	ST2 	Op = 26
	ST3 	Op = 27
	ST4 	Op = 28
	ST5 	Op = 29
	ST6 	Op = 30
	STX 	Op = 31
)

func (b *Op) Scan(state fmt.ScanState, verb rune) error {
	t, err := state.Token(true, nil)
	if err != nil {
		return err
	}
	s := string(t)
	switch s {
	case "LDA": *b = LDA; return nil
	case "LD1": *b = LD1; return nil
	case "LD2": *b = LD2; return nil
	case "LD3": *b = LD3; return nil
	case "LD4": *b = LD4; return nil
	case "LD5": *b = LD5; return nil
	case "LD6": *b = LD6; return nil
	case "LDX": *b = LDX; return nil

	case "LDAN": *b = LDAN; return nil
	case "LD1N": *b = LD1N; return nil
	case "LD2N": *b = LD2N; return nil
	case "LD3N": *b = LD3N; return nil
	case "LD4N": *b = LD4N; return nil
	case "LD5N": *b = LD5N; return nil
	case "LD6N": *b = LD6N; return nil
	case "LDXN": *b = LDXN; return nil

	case "STA": *b = STA; return nil
	case "ST1": *b = ST1; return nil
	case "ST2": *b = ST2; return nil
	case "ST3": *b = ST3; return nil
	case "ST4": *b = ST4; return nil
	case "ST5": *b = ST5; return nil
	case "ST6": *b = ST6; return nil
	case "STX": *b = STX; return nil

	default:
		return errors.New("unknown instructtion token: " + s)
	}
}

func (w Word) GetOp() Op {
	return Op(w.Bytes[4])
}

func (w *Word) SetOp(o Op) {
	w.Bytes[4] = byte(o)
}

func (m *Mix) Do(i *Word) {
	switch  o := i.GetOp(); o {
	case LDA:
		v := m.GetMem(i)
		m.RA = v
	case LDX:
		v := m.GetMem(i)
		m.RX = v
		break
	case LD1, LD2, LD3, LD4, LD5, LD6:
		v := m.GetMem(i)
		m.RI[o - LD1].FromWord(&v)
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
		m.RI[o - LD1N].FromWord(&v)
		break
	case STA:
		a := m.EffectiveAdr(i)
		m.Memory[a].SetField(i.GetMod(), m.RA)
		break
	case ST1, ST2, ST3, ST4, ST5, ST6:
		a := m.EffectiveAdr(i)
		w := m.RI[o - ST1].ToWord()
		m.Memory[a].SetField(i.GetMod(), w)
		break
	case STX:
		a := m.EffectiveAdr(i)
		m.Memory[a].SetField(i.GetMod(), m.RX)
		break
	}
}
