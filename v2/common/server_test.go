package common

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseServerTime(t *testing.T) {
	assert := assert.New(t)

	t.Run("valid serverTime", func(t *testing.T) {
		serverTime, err := ParseServerTime([]byte(`{"serverTime": 1499827319559}`))
		assert.NoError(err)
		assert.EqualValues(1499827319559, serverTime)
	})

	t.Run("missing serverTime field", func(t *testing.T) {
		serverTime, err := ParseServerTime([]byte(`{}`))
		assert.NoError(err)
		assert.EqualValues(0, serverTime)
	})

	t.Run("invalid json", func(t *testing.T) {
		serverTime, err := ParseServerTime([]byte(``))
		assert.Error(err)
		assert.EqualValues(0, serverTime)
	})
}

func TestCalculateTimeOffset(t *testing.T) {
	assert := assert.New(t)

	now := time.Now().UnixNano() / int64(time.Millisecond)

	// A server in sync with the local clock yields an offset close to zero.
	// Allow a generous window to absorb test-execution latency.
	offset := CalculateTimeOffset(now)
	assert.LessOrEqual(offset, int64(1000))
	assert.GreaterOrEqual(offset, int64(-1000))

	// A server reporting a time 5s in the past yields a positive offset of ~5000ms.
	offset = CalculateTimeOffset(now - 5000)
	assert.GreaterOrEqual(offset, int64(4000))
}
