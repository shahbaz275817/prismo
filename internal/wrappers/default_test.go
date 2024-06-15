package wrappers

import (
	"encoding/json"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type DefaultRouteHandlerSuite struct {
	suite.Suite
}

func TestDefaultRouteHandlerSuite(t *testing.T) {
	suite.Run(t, new(DefaultRouteHandlerSuite))
}

func (suite *DefaultRouteHandlerSuite) TestPanic() {
	handler := DefaultWrapper(func(http.ResponseWriter, *http.Request) error {
		panic("A sample panic")
	})
	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/", nil)
	suite.Assert().Panics(func() {
		handler.ServeHTTP(w, r)
	})

	suite.Assert().Equal(500, w.Result().StatusCode)

	rBody := w.Result().Body
	body, _ := ioutil.ReadAll(rBody)
	_ = rBody.Close()

	var b map[string]interface{}
	_ = json.Unmarshal(body, &b)

	suite.Assert().Equal(false, b["success"])
}
