package password

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Password(t *testing.T) {
	pass := &Password{
		PlainTextPass: "123456789",
	}

	t.Run("should hash password", func(t *testing.T) {
		err := pass.HashPassword(pass.PlainTextPass)

		assert.NoError(t, err)

		assert.NotNil(t, pass.Hash)
	})

	t.Run("hashed pass should match with plain text pass", func(t *testing.T) {
		isMatch, err := pass.Matches(pass.PlainTextPass)

		assert.NoError(t, err)

		assert.True(t, isMatch)
	})

	t.Run("should always be different hash for same pass", func(t *testing.T) {
		pass1 := &Password{
			PlainTextPass: "123456789",
		}

		pass2 := &Password{
			PlainTextPass: "123456789",
		}

		pass1.HashPassword(pass1.PlainTextPass)
		pass2.HashPassword(pass2.PlainTextPass)

		assert.NotEqual(t, pass1.Hash, pass2.Hash)
	})
}
