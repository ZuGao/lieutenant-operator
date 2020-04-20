package vault

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/operator-framework/operator-sdk/pkg/log/zap"
)

func testGetHTTPServer(statusCode int, body []byte) *httptest.Server {
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(statusCode)
		_, _ = res.Write(body)
	}))

	return testServer
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name       string
		want       *VaultClient
		wantErr    bool
		statusCode int
		body       string
	}{
		{
			name:       "test create secret",
			wantErr:    false,
			statusCode: 200,
			body: `{
				"data": {
				  "created_time": "2018-03-22T02:36:43.986212308Z",
				  "deletion_time": "",
				  "destroyed": false,
				  "version": 1
				}
			  }`,
		},
		{
			name:       "test error",
			wantErr:    true,
			statusCode: 404,
			body:       "{}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			instanceClient = nil
			server := testGetHTTPServer(tt.statusCode, []byte(tt.body))

			os.Setenv(tokenEnvVar, "myroot")
			os.Setenv(vaultEnvVar, server.URL)

			defer server.Close()

			client, err := NewClient()
			if err != nil {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			err = client.SetToken("1234/6789", "test", zap.Logger())
			if (err != nil) != tt.wantErr {
				t.Errorf("Write() error = %v", err)
			}
		})
	}
}
