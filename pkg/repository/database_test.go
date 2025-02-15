package repository

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewDatabase tests the database factory function
func TestNewDatabase(t *testing.T) {
	t.Run("unsupported driver", func(t *testing.T) {
		_, err := NewDatabase()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "database driver not supported")
	})
}
