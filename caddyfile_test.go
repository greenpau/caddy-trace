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
        request_debug disabled=yes
        request_debug disabled=no tag="foo"
        request_debug disabled=no tag="bar"
        respond /version 200 {
          body "1.0.0"
        }
        request_debug tag="marvel"
		respond /whoami 200 {
          body "greenpau"
        }
      }
    }
    `, "caddyfile")

	tester.AssertGetResponse("http://localhost:9080/version", 200, "1.0.0")
	tester.AssertGetResponse("http://localhost:9080/whoami", 200, "greenpau")
}
