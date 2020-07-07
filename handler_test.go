package gotesting

import (
	"errors"
	"io"
	"net/http/httptest"
	"strings"
	"testing"
)

type errorReader struct{}

func (errorReader) Read(p []byte) (n int, err error) {
	return 0, errors.New("mock body error")
}

func TestLoginHandler(t *testing.T) {

	testCases := []struct {
		Name string
		Code int
		Body interface{}
	}{
		{"ok", 200, `{"code":"a@example.com", "password":"password"}`},
		{"read body error", 500, new(errorReader)},
		{"invalid format", 400, `{"code":1, "password":"password"}`},
		{"invalid code", 400, `{"code":"a@example.com1", "password":"password"}`},
		{"invalid password", 400, `{"code":"a@example.com", "password":"password1"}`},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {

			var body io.Reader
			if stringBody, ok := tc.Body.(string); ok {
				body = strings.NewReader(stringBody)
			} else {
				body = tc.Body.(io.Reader)
			}

			req := httptest.NewRequest("POST", "http://example.com/foo", body)
			w := httptest.NewRecorder()

			LoginHandler(w, req)

			resp := w.Result()
			if resp.StatusCode != tc.Code {
				t.Errorf("response code is invalid, expect=%d but got=%d",
					tc.Code, resp.StatusCode)
			}
		})
	}
}
