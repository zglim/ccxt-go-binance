package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseServerTime(t *testing.T) {
	t.Run("valid response", func(t *testing.T) {
		data := []byte(`{"serverTime": 1499827319559}`)
		st, err := ParseServerTime(data)
		assert.NoError(t, err)
		assert.EqualValues(t, 1499827319559, st)
	})

	t.Run("large timestamp", func(t *testing.T) {
		data := []byte(`{"serverTime": 1692387156596}`)
		st, err := ParseServerTime(data)
		assert.NoError(t, err)
		assert.EqualValues(t, 1692387156596, st)
	})

	t.Run("zero serverTime", func(t *testing.T) {
		data := []byte(`{"serverTime": 0}`)
		st, err := ParseServerTime(data)
		assert.NoError(t, err)
		assert.EqualValues(t, 0, st)
	})

	t.Run("missing serverTime field returns 0", func(t *testing.T) {
		data := []byte(`{}`)
		st, err := ParseServerTime(data)
		assert.NoError(t, err)
		assert.EqualValues(t, 0, st)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		data := []byte(``)
		_, err := ParseServerTime(data)
		assert.Error(t, err)
	})

	t.Run("non-object JSON", func(t *testing.T) {
		data := []byte(`1399827319559`)
		_, err := ParseServerTime(data)
		// simplejson.NewJson accepts scalars; .Get("serverTime") on a scalar returns zero
		assert.NoError(t, err)
	})
}

func TestComputeTimeOffset(t *testing.T) {
	t.Run("local ahead of server", func(t *testing.T) {
		offset := ComputeTimeOffset(1500000000000, 1499999999000)
		assert.EqualValues(t, 1000, offset)
	})

	t.Run("local behind server", func(t *testing.T) {
		offset := ComputeTimeOffset(1499999998000, 1499999999000)
		assert.EqualValues(t, -1000, offset)
	})

	t.Run("clocks in sync", func(t *testing.T) {
		offset := ComputeTimeOffset(1499999999000, 1499999999000)
		assert.EqualValues(t, 0, offset)
	})
}
