package http

import (
	"testing"
)

func Test_httpServer_ListenAndServe(t *testing.T) {
	hs := NewHttpServer(nil, nil, nil, nil)
	go hs.ListenAndServe(":8000")
}
