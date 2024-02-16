package delete

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/wlcmtunknwndth/REST_API/internal/http-server/handlers/delete/mocks"
	"github.com/wlcmtunknwndth/REST_API/internal/http-server/handlers/url/save"
	"github.com/wlcmtunknwndth/REST_API/internal/lib/logger/handlers/slogdiscard"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeleteFunc(t *testing.T) {
	cases := []struct {
		name      string
		alias     string
		url       string
		respError string
		mockError error
	}{
		{
			name:      "Success",
			alias:     "test_alias",
			respError: "invalid request",
			//url:   "https://google.com",
		},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			urlDeleteMock := mocks.NewURLDelete(t)

			if tc.respError == "" || tc.mockError != nil {
				urlDeleteMock.On("Delete", tc.alias,
					mock.AnythingOfType("string"),
				).Return(tc.mockError).Once()
			}

			var handler = New(slogdiscard.NewDiscardLogger(), urlDeleteMock)
			input := fmt.Sprintf(`{"alias": "%s"}`, tc.alias)

			req, err := http.NewRequest(http.MethodDelete, "/delete", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

		})
	}

}
