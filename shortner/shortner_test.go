package base62shortner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_base10ToBase62(t *testing.T) {
	t.Run("get base 62 coded string", func(t *testing.T) {
		res := base10ToBase62(100000000000)
		assert.Equal(t, "1L9zO9O", res)
	})
}
