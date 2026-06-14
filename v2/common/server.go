package common

import (
	"time"

	"github.com/bitly/go-simplejson"
)

// ParseServerTime extracts the "serverTime" field (a Unix timestamp in
// milliseconds) from a raw /time endpoint response body. It returns an error
// only when the body cannot be parsed as JSON; a missing serverTime field
// yields 0, matching the behaviour all market clients have historically relied
// on.
func ParseServerTime(data []byte) (int64, error) {
	j, err := simplejson.NewJson(data)
	if err != nil {
		return 0, err
	}
	return j.Get("serverTime").MustInt64(), nil
}

// CalculateTimeOffset returns the difference, in milliseconds, between the
// local clock and the given server time (localNow - serverTime). Clients store
// this value as TimeOffset and apply it to outgoing signed-request timestamps
// so they stay within Binance's recvWindow.
func CalculateTimeOffset(serverTime int64) int64 {
	return time.Now().UnixNano()/int64(time.Millisecond) - serverTime
}
