package main

import (
	"github.com/rk1165/feedcreator/internal/assert"
	"net/http"
	"testing"
)

func TestFeedView(t *testing.T) {
	app := newTestApplication(t)
	ts := newTestServer(t, app.routes())
	defer ts.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/feed/view?id=1",
			wantCode: http.StatusOK,
			wantBody: "Feed for example",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/feed/view?id=2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/feed/view?id=-1",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Decimal ID",
			urlPath:  "/feed/view?id=1.23",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "String ID",
			urlPath:  "/feed/view?id=foo",
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Empty ID",
			urlPath:  "/feed/view?id=",
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			code, _, body := ts.get(t, tc.urlPath)
			assert.Equal(t, code, tc.wantCode)
			if tc.wantBody != "" {
				assert.StringContains(t, body, tc.wantBody)
			}
		})
	}
}
