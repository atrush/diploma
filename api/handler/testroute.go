package handler

import (
	"bytes"
	"github.com/atrush/diploma.git/services/auth"
	"github.com/atrush/diploma.git/services/order"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http/httptest"
	"testing"
)

//  HandlerTest stores test data and expected response params.
type TestRoute struct {
	svcAuth  auth.Authenticator //  authentication service
	svcOrder order.OrderManager //  orders service

	name    string            //  test name
	method  string            //  http method
	url     string            //  request url
	body    string            //  request body content
	headers map[string]string //  request headers

	expectedBody    string            //  expected response body
	expectedHeaders map[string]string //  expected response headers
	expectedCode    int               // expected response http status
}

//  CheckTest runs handler, builds request and checks response values
func (tt *TestRoute) CheckTest(t *testing.T) {
	//  new handler with mock services
	h, err := NewHandler(tt.svcAuth, tt.svcOrder)
	require.NoError(t, err)

	//  new router with handler
	r := NewRouter(h)
	request := httptest.NewRequest(tt.method, tt.url, bytes.NewBuffer([]byte(tt.body)))

	//  set headers
	if len(tt.headers) > 0 {
		for header, value := range tt.headers {
			request.Header.Set(header, value)
		}
	}

	//  make request
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)

	res := w.Result()
	resBody, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	defer res.Body.Close()

	strBody := string(resBody)

	//  check response headers
	if len(tt.expectedHeaders) > 0 {
		for expHeader, expValue := range tt.expectedHeaders {
			resHeader := res.Header.Get(expHeader)
			assert.Equal(t, expValue, resHeader, "ожидался заголовок %v, значение %v, получено значение %v",
				expHeader, expValue, resHeader)
		}
	}

	//  check response code
	if tt.expectedCode != 0 {
		assert.True(t, res.StatusCode == tt.expectedCode, "Ожидался код ответа %d, получен %d", tt.expectedCode, w.Code)
	}

	//  check response body
	if tt.expectedBody != "" {
		assert.Equal(t, strBody, tt.expectedBody, "Ожидался ответа %v, получен %v", tt.expectedBody, strBody)
	}
}
