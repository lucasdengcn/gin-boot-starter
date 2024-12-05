package exception

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthErrorIs(t *testing.T) {
	err := &AuthError{Errno: "AUTH_401", CorrelationId: "test", Message: "Invalid token"}
	assert.True(t, errors.Is(err, &AuthError{}))
}
