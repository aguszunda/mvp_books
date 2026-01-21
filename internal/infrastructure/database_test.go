package infrastructure

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDSN(t *testing.T) {
	// Setup env vars
	os.Setenv("DB_USER", "user")
	os.Setenv("DB_PASSWORD", "pass")
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_NAME", "db")

	t.Run("Default Port", func(t *testing.T) {
		os.Unsetenv("DB_PORT")
		dsn := GetDSN()
		assert.Contains(t, dsn, "tcp(localhost:3306)")
	})

	t.Run("Custom Port", func(t *testing.T) {
		os.Setenv("DB_PORT", "3307")
		dsn := GetDSN()
		assert.Contains(t, dsn, "tcp(localhost:3307)")
	})
}
