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
	"fmt"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"regexp"
	"strings"
)

func init() {
	httpcaddyfile.RegisterHandlerDirective("trace", parseCaddyfileRequestDebugger)
}

// parseCaddyfileRequestDebugger sets up a handler for debugging http requests. Syntax:
//
//     trace [log_level <debug|info|warn|error>] [disabled <yes|no>] [enable_uuid <yes|no>]
//
// The disabled is being set to true, there will be no output from the plugin.
// If enable_uuid is being set to true, then the plugin adds request_id field to
// its output.
func parseCaddyfileRequestDebugger(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var dbg RequestDebugger
	for h.Next() {
		args := h.RemainingArgs()
		if len(args) == 0 {
			dbg.Disabled = false
			return dbg, nil
		}
		for _, arg := range args {
			parts := strings.SplitN(arg, "=", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("unsupported argument: %s", arg)
			}
			k := parts[0]
			v := strings.Trim(parts[1], "\"")
			switch k {
			case "tag":
				dbg.Tag = v
			case "disabled":
				if !isSwitchArg(v) {
					return nil, fmt.Errorf("%s argument value of %s is unsupported", k, v)
				}
				if isEnabledArg(v) {
					dbg.Disabled = true
				}
			case "response_debug":
				if !isSwitchArg(v) {
					return nil, fmt.Errorf("%s argument value of %s is unsupported", k, v)
				}
				if isEnabledArg(v) {
					dbg.ResponseDebugEnabled = true
				}
			case "uri_filter":
				dbg.URIFilter = v
				if _, err := regexp.CompilePOSIX(dbg.URIFilter); err != nil {
					return nil, fmt.Errorf("%s directive value of %s fails to compile: %s", k, v, err)
				}
			default:
				return nil, fmt.Errorf("unsupported argument: %s", arg)
			}
		}
	}
	return dbg, nil
}

func isEnabledArg(s string) bool {
	if s == "yes" || s == "true" || s == "on" {
		return true
	}
	return false
}

func isSwitchArg(s string) bool {
	if s == "yes" || s == "true" || s == "on" {
		return true
	}
	if s == "no" || s == "false" || s == "off" {
		return true
	}
	return false
}
