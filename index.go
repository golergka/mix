package mix

type Index struct {
	Bytes [2]byte
	Sign  Sign
}

func (i *Index) GetAdr() SignAdr {
	return MakeSignAdr(&i.Sign, i.Bytes[:])
}

func (i *Index) FromWord(w *Word) {
	i.Sign = w.Sign
	copy(i.Bytes[:], w.Bytes[3:5])
}

func (i *Index) ToWord() Word {
	r := Word{Sign: i.Sign}
	copy(r.Bytes[3:5], i.Bytes[:])
	return r
}
