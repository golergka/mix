package mix

import "testing"
import "fmt"
import "github.com/stretchr/testify/assert"

func TestOpScan(t *testing.T) {
	{
		var o Op
		s := "LDA"

		_, err := fmt.Sscanf(s, "%v", &o)

		assert.Nil(t, err)
		assert.Equal(t, LDA, o)
	}
	{
		var o Op
		s := "asdasd"

		_, err := fmt.Sscanf(s, "%v", &o)

		assert.NotNil(t, err)
	}
}
