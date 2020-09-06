package debug

import (
	"testing"

	"github.com/caddyserver/caddy/v2/caddytest"
)

func TestCaddyfile(t *testing.T) {
	tester := caddytest.NewTester(t)
	tester.InitServer(`
    {
      http_port     9080
      https_port    9443
    }

    localhost:9080 {
      route /* {
	    request_debug
	    request_debug
  	    respond /version 200 {
	      body "1.0.0"
	    }
	  }
    }
    `, "caddyfile")

	tester.AssertGetResponse("http://localhost:9080/version", 200, "1.0.0")
}
