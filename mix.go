package mix

import "fmt"
import _ "encoding/binary"

type Sign bool

type Word struct {
	Bytes	[5]byte // 0 is MSB, 4 is LSB
	Sign	Sign
}

func (s* Sign) String() string {
	if *s {
		return "+"
	} else {
		return "-"
	}
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

type Registers struct {
	RA	Word
	RX	Word
	RI1	Index
	RI2	Index
	RI3	Index
	RI4	Index
	RI5	Index
	RI6	Index

	O	bool
	C	Comparison
}

type Mix struct {
	Registers
	Memory [4000]byte
}

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
