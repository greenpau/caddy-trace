package debug

import (
	"net/http"
	"net/url"
	"testing"

	"github.com/caddyserver/caddy/v2/caddytest"
)

func TestCaddyfile(t *testing.T) {
	baseURL := "http://localhost:9080"
	localhost, _ := url.Parse(baseURL)
	tester := caddytest.NewTester(t)
	tester.InitServer(`
    {
      http_port     9080
      https_port    9443
    }

    localhost:9080 {
      route {
        request_debug disabled=yes
        request_debug disabled=no tag="foo"
        request_debug disabled=no tag="bar"
        respond /version 200 {
          body "1.0.0"
        }
        request_debug tag="marvel"
		respond /whoami* 200 {
          body "greenpau"
        }
      }
    }
    `, "caddyfile")

	tester.AssertGetResponse(baseURL+"/version", 200, "1.0.0")
	tester.AssertGetResponse(baseURL+"/whoami", 200, "greenpau")

	cookies := []*http.Cookie{}
	cookie := &http.Cookie{
		Name:  "access_code",
		Value: "anonymous",
	}

	cookies = append(cookies, cookie)
	tester.Client.Jar.SetCookies(localhost, cookies)
	tester.AssertGetResponse(baseURL+"/whoami", 200, "greenpau")
	tester.AssertGetResponse(baseURL+"/whoami?user=greenpau&mode=raw", 200, "greenpau")
}
