package util

import (
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestPassword(t *testing.T) {
	t.Run("correct password", func(t *testing.T) {
		password := RandomString(8)

		hashedPassword, err := HashPassword(password)
		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword)

		err = CheckPassword(password, hashedPassword)
		require.NoError(t, err)
	})

	t.Run("incorrect password", func(t *testing.T) {
		password := RandomString(8)
		falsePassword := RandomString(10)

		hashedPassword, err := HashPassword(password)
		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword)

		err = CheckPassword(falsePassword, hashedPassword)
		require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())
	})

	t.Run("all hashed password are different", func(t *testing.T) {
		password := RandomString(6)

		hashedPassword1, err := HashPassword(password)
		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword1)

		err = CheckPassword(password, hashedPassword1)
		require.NoError(t, err)

		wrongPassword := RandomString(6)
		err = CheckPassword(wrongPassword, hashedPassword1)
		require.EqualError(t, err, bcrypt.ErrMismatchedHashAndPassword.Error())

		hashedPassword2, err := HashPassword(password)
		require.NoError(t, err)
		require.NotEmpty(t, hashedPassword2)
		require.NotEqual(t, hashedPassword1, hashedPassword2)
	})
}
