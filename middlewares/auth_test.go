package middlewares

import (
	"github.com/jrmarcco/go-backend-demo/api"
	"github.com/jrmarcco/go-backend-demo/util"
	"testing"
	"time"
)

func TestAuthMiddleware(t *testing.T) {
	tcs := []struct {
		name string
	}{
		{
			name: "Normal Case",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {

			// build test server
			_ = api.NewServer(util.ServerCfg{TokenDuration: time.Minute}, nil)

			// register middlewares

			// setup authorization

		})
	}
}
