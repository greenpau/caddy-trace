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
	//"encoding/json"
	//"strconv"
	//"strings"

	//"github.com/caddyserver/caddy/v2"
	//"github.com/caddyserver/caddy/v2/caddyconfig"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
)

func init() {
	//httpcaddyfile.RegisterHandlerDirective("uri", parseCaddyfileRequestDebugger)
	httpcaddyfile.RegisterHandlerDirective("request_debug", parseCaddyfileRequestDebugger)
}

// parseCaddyfileRequestDebugger sets up a handler for debugging http requests. Syntax:
//
//     request_debug [log_level <debug|info|warn|error>] [disabled <yes|no>] [enable_uuid <yes|no>]
//     uri [<matcher>] strip_prefix|strip_suffix|replace <target> [<replacement> [<limit>]]
//
//
// The disabled is being set to true, there will be no output from the plugin.
// If enable_uuid is being set to true, then the plugin adds request_id field to
// its output.
func parseCaddyfileRequestDebugger(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var dbg RequestDebugger
	for h.Next() {
		args := h.RemainingArgs()
		if len(args) == 0 {
			dbg.LogLevel = "error"
			dbg.Disabled = false
			dbg.EnableUUID = true
			return dbg, nil
		}

		if len(args) > 0 {
			return nil, fmt.Errorf("the number of elements is %d", len(args))

		}

		/*
			switch args[0] {
			case "strip_prefix":
				if len(args) > 2 {
					return nil, h.ArgErr()
				}
				dbg.StripPathPrefix = args[1]
				if !strings.HasPrefix(dbg.StripPathPrefix, "/") {
					dbg.StripPathPrefix = "/" + dbg.StripPathPrefix
				}
			case "strip_suffix":
				if len(args) > 2 {
					return nil, h.ArgErr()
				}
				dbg.StripPathSuffix = args[1]
			case "replace":
				var find, replace, lim string
				switch len(args) {
				case 4:
					lim = args[3]
					fallthrough
				case 3:
					find = args[1]
					replace = args[2]
				default:
					return nil, h.ArgErr()
				}

				var limInt int
				if lim != "" {
					var err error
					limInt, err = strconv.Atoi(lim)
					if err != nil {
						return nil, h.Errf("limit must be an integer; invalid: %v", err)
					}
				}

				dbg.URISubstring = append(dbg.URISubstring, replacer{
					Find:    find,
					Replace: replace,
					Limit:   limInt,
				})
			default:
				return nil, h.Errf("unrecognized URI manipulation '%s'", args[0])
			}
		*/
	}
	return dbg, nil
}
