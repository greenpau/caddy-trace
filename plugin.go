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
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"os"
	"strconv"
)

func init() {
	caddy.RegisterModule(RequestDebugger{})
}

// RequestDebugger is a middleware which displays the content of the request it
// handles. It helps troubleshooting web requests by exposing headers
// (e.g. cookies), URL parameters, etc.
type RequestDebugger struct {
	// Enables or disables the plugin.
	Disabled bool `json:"disabled,omitempty"`
	// Adds a tag to a log message
	Tag    string `json:"tag,omitempty"`
	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (RequestDebugger) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.request_debug",
		New: func() caddy.Module { return new(RequestDebugger) },
	}
}

// Provision sets up RequestDebugger.
func (dbg *RequestDebugger) Provision(ctx caddy.Context) error {
	// dbg.logger = ctx.Logger(dbg)
	if dbg.logger == nil {
		dbg.logger = initLogger()
	}
	return nil
}

func (dbg RequestDebugger) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	if dbg.Disabled {
		return next.ServeHTTP(w, r)
	}
	dbg.debug(r, "incoming")
	return next.ServeHTTP(w, r)
}

func (dbg *RequestDebugger) debug(r *http.Request, reqDirection string) {
	var requestID string
	rawRequestID := caddyhttp.GetVar(r.Context(), "request_id")
	if rawRequestID == nil {
		requestID = uuid.NewV4().String()
		caddyhttp.SetVar(r.Context(), "request_id", requestID)
	} else {
		requestID = rawRequestID.(string)
	}

	cookies := r.Cookies()

	var remotePort int
	remoteAddr, remotePortStr, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		dbg.logger.Error(
			"request debugging: encountered source ip parsing error",
			zap.Any("request_id", requestID),
			zap.String("direction", reqDirection),
			zap.String("tag", dbg.Tag),
			zap.String("error", err.Error()),
		)
	}

	if remotePortStr != "" {
		remotePort, err = strconv.Atoi(remotePortStr)
		if err != nil {
			dbg.logger.Error(
				"request debugging: encountered source port parsing error",
				zap.Any("request_id", requestID),
				zap.String("direction", reqDirection),
				zap.String("tag", dbg.Tag),
				zap.String("error", err.Error()),
			)
		}
	}

	// Extract query parameters
	queryParams := make(map[string]interface{})
	queryValues := r.URL.Query()
	for k, v := range queryValues {
		if len(v) == 1 {
			queryParams[k] = v[0]
		} else {
			queryParams[k] = v
		}
	}

	// Extract headers
	reqHeaders := make(map[string]interface{})
	if r.Header != nil {
		for k, v := range r.Header {
			if k == "Cookie" || k == "Set-Cookie" {
				continue
			}
			if len(v) == 1 {
				reqHeaders[k] = v[0]
			} else {
				reqHeaders[k] = v
			}
		}
	}

	dbg.logger.Debug(
		"debugging request",
		zap.Any("request_id", requestID),
		zap.String("direction", reqDirection),
		zap.String("tag", dbg.Tag),
		zap.String("method", r.Method),
		zap.String("proto", r.Proto),
		zap.String("host", r.Host),
		zap.String("uri", r.RequestURI),
		zap.String("remote_addr_port", r.RemoteAddr),
		zap.String("remote_addr", remoteAddr),
		zap.Int("remote_port", remotePort),
		zap.Int64("content_length", r.ContentLength),
		zap.Int("cookie_count", len(cookies)),
		zap.String("user_agent", r.UserAgent()),
		zap.String("referer", r.Referer()),
		zap.Any("cookies", cookies),
		zap.Any("query_params", queryParams),
		zap.Any("headers", reqHeaders),
	)

}

func initLogger() *zap.Logger {
	logAtom := zap.NewAtomicLevel()
	logAtom.SetLevel(zapcore.DebugLevel)
	logEncoderConfig := zap.NewProductionEncoderConfig()
	logEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	logEncoderConfig.TimeKey = "time"
	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(logEncoderConfig),
		zapcore.Lock(os.Stdout),
		logAtom,
	))
	return logger

}

// Interface guard
var _ caddyhttp.MiddlewareHandler = (*RequestDebugger)(nil)
