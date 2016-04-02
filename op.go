package mix

import "fmt"
import "errors"

type Op byte

//go:generate stringer -type=Op
const (
	LDA Op = 8
	LD1 Op = 9
	LD2 Op = 10
	LD3 Op = 11
	LD4 Op = 12
	LD5 Op = 13
	LD6 Op = 14
	LDX Op = 15

	LDAN Op = 16
	LD1N Op = 17
	LD2N Op = 18
	LD3N Op = 19
	LD4N Op = 20
	LD5N Op = 21
	LD6N Op = 22
	LDXN Op = 23

	STA Op = 24
	ST1 Op = 25
	ST2 Op = 26
	ST3 Op = 27
	ST4 Op = 28
	ST5 Op = 29
	ST6 Op = 30
	STX Op = 31
)

func (b *Op) Scan(state fmt.ScanState, verb rune) error {
	t, err := state.Token(true, nil)
	if err != nil {
		return err
	}
	s := string(t)
	switch s {
	case "LDA":
		*b = LDA
		return nil
	case "LD1":
		*b = LD1
		return nil
	case "LD2":
		*b = LD2
		return nil
	case "LD3":
		*b = LD3
		return nil
	case "LD4":
		*b = LD4
		return nil
	case "LD5":
		*b = LD5
		return nil
	case "LD6":
		*b = LD6
		return nil
	case "LDX":
		*b = LDX
		return nil

	case "LDAN":
		*b = LDAN
		return nil
	case "LD1N":
		*b = LD1N
		return nil
	case "LD2N":
		*b = LD2N
		return nil
	case "LD3N":
		*b = LD3N
		return nil
	case "LD4N":
		*b = LD4N
		return nil
	case "LD5N":
		*b = LD5N
		return nil
	case "LD6N":
		*b = LD6N
		return nil
	case "LDXN":
		*b = LDXN
		return nil

	case "STA":
		*b = STA
		return nil
	case "ST1":
		*b = ST1
		return nil
	case "ST2":
		*b = ST2
		return nil
	case "ST3":
		*b = ST3
		return nil
	case "ST4":
		*b = ST4
		return nil
	case "ST5":
		*b = ST5
		return nil
	case "ST6":
		*b = ST6
		return nil
	case "STX":
		*b = STX
		return nil

	default:
		return errors.New("unknown instructtion token: " + s)
	}
}

