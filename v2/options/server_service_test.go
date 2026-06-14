package options

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/adshao/go-binance/v2/common"
	"github.com/stretchr/testify/suite"
)

type ServerServiceTestSuite struct {
	baseTestSuite
}

func TestPingService(t *testing.T) {
	suite.Run(t, new(ServerServiceTestSuite))
}

func (s *ServerServiceTestSuite) TestPing() {
	data := []byte(`{}`)
	s.mockDo(data, nil)
	defer s.assertDo()

	err := s.client.NewPingService().Do(newContext())
	s.r().Equal(err, nil, "err != nil")
}

func (s *ServerServiceTestSuite) TestPingError() {
	s.mockDo([]byte("{}"), fmt.Errorf("dummy error"), http.StatusInternalServerError)
	defer s.assertDo()

	err := s.client.NewPingService().Do(newContext())
	s.r().Error(err)
	s.r().Contains(err.Error(), "dummy error")
}

func (s *ServerServiceTestSuite) TestPingBadRequest() {
	s.mockDo([]byte(`{
		"code": -1121,
		"msg": "Invalid symbol."
	}`), nil, http.StatusBadRequest)
	defer s.assertDo()

	err := s.client.NewPingService().Do(newContext())
	s.r().Error(err)
	s.r().True(common.IsAPIError(err))
}

func (s *ServerServiceTestSuite) TestServerTime() {
	data := []byte(`{
		"serverTime": 1592387156596
	}`)
	s.mockDo(data, nil)
	defer s.assertDo()

	st, err := s.client.NewServerTimeService().Do(newContext())
	var targetServerTime int64 = 1592387156596
	s.r().Equal(st, targetServerTime, "serverTime")
	s.r().Equal(err, nil, "err != nil")
}

func (s *ServerServiceTestSuite) TestServerTimeError() {
	s.mockDo([]byte("{}"), fmt.Errorf("dummy error"), http.StatusInternalServerError)
	defer s.assertDo()

	_, err := s.client.NewServerTimeService().Do(newContext())
	s.r().Error(err)
	s.r().Contains(err.Error(), "dummy error")
}

func (s *ServerServiceTestSuite) TestServerTimeBadRequest() {
	s.mockDo([]byte(`{
		"code": -1121,
		"msg": "Invalid symbol."
	}`), nil, http.StatusBadRequest)
	defer s.assertDo()

	_, err := s.client.NewServerTimeService().Do(newContext())
	s.r().Error(err)
	s.r().True(common.IsAPIError(err))
}

func (s *ServerServiceTestSuite) TestInvalidResponseBody() {
	s.mockDo([]byte(``), nil)
	defer s.assertDo()

	_, err := s.client.NewServerTimeService().Do(newContext())
	s.r().Error(err)
	s.r().False(common.IsAPIError(err))
}
