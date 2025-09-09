package http

import (
	"net/http"
	"testing"

	"github.com/Jashanveer-Singh/todo-go/test/mocks"
	"github.com/golang/mock/gomock"
)

func TestAuthMiddleware_isAuthenticatedMiddleware(t *testing.T) {
	tests := []struct {
		name     string // description of this test case
		setupMTP func(mtp *mocks.MockTokenProvider)
		// Named input parameters for target function.
		next http.HandlerFunc
		want http.HandlerFunc
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockTokenProvider := mocks.NewMockTokenProvider(ctrl)
			tt.setupMTP(mockTokenProvider)

			req := http.NewRequest(http, url string, body io.Reader)

			am := NewAuthMiddleware(mockTokenProvider)
			am.isAuthenticatedMiddleware(tt.next).ServeHTTP(w http.ResponseWriter, r *http.Request)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("isAuthenticatedMiddleware() = %v, want %v", got, tt.want)
			}
		})
	}
}
