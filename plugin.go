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
	"bytes"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"net"
	"net/http"
	"os"
	"regexp"
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
	Tag string `json:"tag,omitempty"`
	// Adds response buffering and debugging
	ResponseDebugEnabled bool `json:"response_debug_enabled,omitempty"`
	// If URIFilter is not empty, then only the plugin
	// intercepts only the pages matching the regular expression
	// in the filter
	URIFilter        string `json:"uri_filter,omitempty"`
	filterURIPattern *regexp.Regexp
	logger           *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (RequestDebugger) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.trace",
		New: func() caddy.Module { return new(RequestDebugger) },
	}
}

// Provision sets up RequestDebugger.
func (dbg *RequestDebugger) Provision(ctx caddy.Context) error {
	// dbg.logger = ctx.Logger(dbg)
	if dbg.logger == nil {
		dbg.logger = initLogger()
	}

	if dbg.URIFilter != "" {
		p, err := regexp.CompilePOSIX(dbg.URIFilter)
		if err != nil {
			return fmt.Errorf("failed to compile uri_filter %s: %s", dbg.URIFilter, err)
		}
		dbg.filterURIPattern = p
	}
	return nil
}

func (dbg RequestDebugger) ServeHTTP(resp http.ResponseWriter, req *http.Request, next caddyhttp.Handler) error {
	if dbg.Disabled {
		return next.ServeHTTP(resp, req)
	}
	dbg.debugRequest(req)
	if dbg.ResponseDebugEnabled {
		buf := new(bytes.Buffer)
		shouldBuffer := func(status int, header http.Header) bool {
			return true
		}
		wrapResp := caddyhttp.NewResponseRecorder(resp, buf, shouldBuffer)
		err := next.ServeHTTP(wrapResp, req)
		if err != nil {
			return err
		}
		if !wrapResp.Buffered() {
			return nil
		}
		dbg.debugResponse(req, wrapResp)
		wrapResp.WriteResponse()
		return nil
	}
	return next.ServeHTTP(resp, req)
}

func (dbg *RequestDebugger) debugResponse(req *http.Request, wrapResp caddyhttp.ResponseRecorder) {
	if dbg.filterURIPattern != nil {
		if !dbg.filterURIPattern.MatchString(req.RequestURI) {
			return
		}
	}
	var requestID string
	direction := "outgoing"
	rawRequestID := caddyhttp.GetVar(req.Context(), "request_id")
	if rawRequestID != nil {
		requestID = rawRequestID.(string)
	}
	bufferSize := 0
	if wrapResp.Buffer() != nil {
		bufferSize = wrapResp.Buffer().Len()
	}

	dbg.logger.Debug(
		"debugging response",
		zap.Any("request_id", requestID),
		zap.String("direction", direction),
		zap.String("tag", dbg.Tag),
		zap.Int("status_code", wrapResp.Status()),
		zap.Int("response_size", wrapResp.Size()),
		zap.Int("buffer_size", bufferSize),
		zap.Any("response_headers", wrapResp.Header()),
	)
}

func (dbg *RequestDebugger) debugRequest(r *http.Request) {
	if dbg.filterURIPattern != nil {
		if !dbg.filterURIPattern.MatchString(r.RequestURI) {
			return
		}
	}
	var requestID string
	reqDirection := "incoming"
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

	// Extract Form
	form := make(map[string]interface{})
	if r.Header != nil {
		if r.Method == "POST" &&
			r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" &&
			r.ContentLength < 1000 {
			r.ParseForm()
			for k, v := range r.Form {
				form[k] = v
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
		zap.Any("form", form),
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
