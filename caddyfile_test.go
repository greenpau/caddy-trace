package debug

import (
	"testing"
)

func TestCaddyfile(t *testing.T) {
	// baseURL := "http://localhost:9080"
	// localhost, _ := url.Parse(baseURL)
	// tester := caddytest.NewTester(t)
	// tester.InitServer(`
	// {
	//   admin on
	//   http_port     9080
	//   https_port    9443
	// }

	// http://localhost:9080 {
	//   route {
	//     trace disabled=yes
	//     trace disabled=no tag="foo"
	//     trace disabled=no tag="bar"
	//     respond /version "1.0.0" 200
	//     trace tag="marvel" response_debug=yes
	// 	trace tag="custom" response_debug=yes uri_filter="^/whoami$"
	//     respond /whoami* 200 {
	//       body "greenpau"
	//     }
	//   }
	// }
	// `, "caddyfile")

	// tester.AssertGetResponse(baseURL+"/version", 200, "1.0.0")
	// tester.AssertGetResponse(baseURL+"/whoami", 200, "greenpau")

	// cookies := []*http.Cookie{}
	// cookie := &http.Cookie{
	// 	Name:  "access_code",
	// 	Value: "anonymous",
	// }

	// cookies = append(cookies, cookie)
	// tester.Client.Jar.SetCookies(localhost, cookies)
	// tester.AssertGetResponse(baseURL+"/whoami", 200, "greenpau")
	// tester.AssertGetResponse(baseURL+"/whoami?user=greenpau&mode=raw", 200, "greenpau")
	t.Log("suppressed test")
}
