package common

import (
	"github.com/bitly/go-simplejson"
)

// ParseServerTime extracts the "serverTime" field from a JSON response body.
// This is the common parsing step shared by all server time services
// (spot, futures, delivery, options).
func ParseServerTime(data []byte) (int64, error) {
	j, err := simplejson.NewJson(data)
	if err != nil {
		return 0, err
	}
	return j.Get("serverTime").MustInt64(), nil
}

// ComputeTimeOffset calculates the offset between local clock and server clock.
// nowMillis should be the current local timestamp in milliseconds (e.g. from currentTimestamp()).
// serverTime is the value returned by ParseServerTime.
// The result is positive when the local clock is ahead of the server.
func ComputeTimeOffset(nowMillis, serverTime int64) int64 {
	return nowMillis - serverTime
}
