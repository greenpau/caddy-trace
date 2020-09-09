// Copyright 2020 Paul Greenberg @greenpau
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package debug

import (
	"net/http"
	"testing"
)

func TestRequestDebugger(t *testing.T) {
	for i, tc := range []struct {
		input *http.Request
		dbg   RequestDebugger
	}{
		{
			dbg:   RequestDebugger{},
			input: newRequest(t, "GET", "/"),
		},
	} {
		tc.dbg.logger = initLogger()
		tc.dbg.debugRequest(tc.input)
		t.Logf("PASS: Test %d", i)
	}
}

func newRequest(t *testing.T, method, uri string) *http.Request {
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		t.Fatalf("error creating request: %v", err)
	}
	req.RequestURI = req.URL.RequestURI()
	return req
}
