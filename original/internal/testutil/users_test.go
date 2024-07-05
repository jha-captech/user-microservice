package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUsers(t *testing.T) {
	count := 3
	IDRangeStart := 3

	users := NewUsers(
		count,
		WithIDStartRange(IDRangeStart),
	)

	for i := IDRangeStart; i < IDRangeStart+count; i++ {
		assert.Equal(t, i, int(users[i-IDRangeStart].ID), "IDs dont match")
	}
}
