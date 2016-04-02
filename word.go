package mix

import "fmt"
import "errors"

type Word struct {
	Bytes [5]byte // 0 is MSB, 4 is LSB
	Sign  Sign
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
			fmt.Fprintf(f, "(%v:%v)", l, h)
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

func scanIndex(state fmt.ScanState) (byte, error) {
	b := byte(0) // Default value
	r, _, err := state.ReadRune()
	if err != nil {
		state.UnreadRune()
		return b, err
	}
	if r != ',' {
		state.UnreadRune()
	} else if _, err := fmt.Fscanf(state, "%v", &b); err != nil {
		return b, err
	}
	return b, nil
}

func scanMod(state fmt.ScanState) (l byte, h byte, e error) {
	l, h = 0, 0 // Default value
	r, _, err := state.ReadRune()
	if err != nil {
		state.UnreadRune()
		return l, h, err
	}
	if r != '(' {
		state.UnreadRune()
	} else if _, err := fmt.Fscanf(state, "%v:%v)", &l, &h); err != nil {
		return l, h, err
	}
	return l, h, nil
}

func (w *Word) Scan(state fmt.ScanState, verb rune) error {
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
			if i, err := scanIndex(state); err != nil {
				return err
			} else {
				w.SetIndex(i)
			}
			if l, h, err := scanMod(state); err != nil {
				return err
			} else {
				w.SetMod(PackMod(l, h))
			}
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
	r := Word{Sign: true}
	r.SetField(mod, w)
	return r
}

func (w Word) GetMod() byte {
	return w.Bytes[3]
}

func (w *Word) SetMod(m byte) {
	w.Bytes[3] = m
}

func (w *Word) SetField(mod byte, v Word) {
	l, h := UnpackMod(mod)
	if l == 0 {
		w.Sign = v.Sign
		l++
	}
	o := 5 - h
	for i := l - 1; i < h; i++ {
		w.Bytes[i+o] = v.Bytes[i]
	}
}

func (w Word) GetAdr() SignAdr {
	return MakeSignAdr(&w.Sign, w.Bytes[0:2])
}

func (w *Word) SetAdr(a SignAdr) {
	s, b := a.ToBytes()
	w.Sign = s
	copy(w.Bytes[0:2], b[:])
}

func (w Word) GetIndex() byte {
	return w.Bytes[2]
}

func (w *Word) SetIndex(i byte) {
	w.Bytes[2] = i
}

func (w Word) GetOp() Op {
	return Op(w.Bytes[4])
}

func (w *Word) SetOp(o Op) {
	w.Bytes[4] = byte(o)
}

