package mix

import "testing"
import "fmt"
import "github.com/stretchr/testify/assert"

func TestWordString(t *testing.T) {
	var m Word
	assert.Equal(t, "- 00 00 00 00 00", m.String())
}

func TestWordFormat(t *testing.T) {
	{
		m := Word{}

		s := fmt.Sprintf("%v", m)

		assert.Equal(t, "- 00 00 00 00 00", s)
	}
	{
		m := Word{Sign: true, Bytes: [5]byte{18, 48, 0, 19, 24}}

		s := fmt.Sprintf("%i", m)

		assert.Equal(t, "STA 1200(2:3)", s)
	}
	{
		m := Word{Sign: false, Bytes: [5]byte{0, 32, 2, 11, 10}}

		s := fmt.Sprintf("%i", m)

		assert.Equal(t, "LD2 -32,2(1:3)", s)
	}
}

func TestWordScan(t *testing.T) {
	{
		var m Word
		var w interface{}

		w = &m

		if _, ok := w.(fmt.Scanner); !ok {
			t.Error("Word doesn't imlpement Scanner")
		}
	}
	{
		s := "- 00 00 00 00 00"
		var m Word

		_, err := fmt.Sscanf(s, "%v", &m)

		assert.Nil(t, err)
		assert.Equal(t, Word{}, m)
	}
	{
		s := "+ 1 2 3 4 5"
		var m Word

		_, err := fmt.Sscanf(s, "%v", &m)

		assert.Nil(t, err)
		assert.Equal(t, Word{Sign: true, Bytes: [5]byte{1, 2, 3, 4, 5}}, m)
	}
	{
		s := "LD2 -32,2(1:3)"
		var m Word

		_, err := fmt.Sscanf(s, "%i", &m)

		assert.Nil(t, err)
		assert.Equal(t, Word{Sign: false, Bytes: [5]byte{0, 32, 2, 11, 10}}, m)
	}
	{
		s := "STA 1200(2:3)"
		var m Word

		_, err := fmt.Scanf(s, "%i", &m)

		assert.Nil(t, err)
		assert.Equal(t, Word{Sign:true, Bytes:[5]byte{18, 48, 0, 19, 24}}, m)
	}
}

func TestWordField(t *testing.T) {
	assert.Equal(t,
		Word{Sign: true, Bytes: [5]byte{0, 0, 10, 11, 0}},
		(&Word{Sign: false, Bytes: [5]byte{10, 11, 0, 11, 22}}).GetField(11))
}

func TestWordSignAdr(t *testing.T) {
	assert.Equal(t, SignAdr(0), (&Word{Sign: false, Bytes: [5]byte{0, 0, 1, 2, 3}}).GetAdr())
	assert.Equal(t, SignAdr(-1), (&Word{Sign: false, Bytes: [5]byte{0, 1, 1, 2, 3}}).GetAdr())
	assert.Equal(t, SignAdr(64), (&Word{Sign: true, Bytes: [5]byte{1, 0, 1, 2, 3}}).GetAdr())
	assert.Equal(t, SignAdr(-5), (&Word{Sign: false, Bytes: [5]byte{0, 5, 1, 2, 3}}).GetAdr())
}

func TestWordGetOp(t *testing.T) {
	assert.Equal(t, LDA, (&Word{Bytes: [5]byte{0, 0, 0, 0, 8}}).GetOp())
	assert.Equal(t, LDA, (&Word{Bytes: [5]byte{1, 2, 4, 16, 8}}).GetOp())
}

func TestWordSetOp(t *testing.T) {
	{
		w := Word{Bytes: [5]byte{1, 2, 3, 4, 5}}

		w.SetOp(LDA)

		assert.Equal(t, byte(8), w.Bytes[4])
	}
}

