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
	"context"
	"github.com/satori/go.uuid"
	"net/http"

	//"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	//"net/http/httputil"
	"os"
)

func init() {
	caddy.RegisterModule(RequestDebugger{})
}

// RequestDebugger is a middleware which displays the content of the request it
// handles. It helps troubleshooting web requests by exposing headers
// (e.g. cookies), URL parameters, etc.
type RequestDebugger struct {
	// Sets logging level. The default level is Error.
	LogLevel string `json:"log_level,omitempty"`
	// Enables or disables the plugin.
	Disabled bool `json:"disabled,omitempty"`
	// Generate UUIDs for requests
	EnableUUID bool `json:"enable_uuid,omitempty"`
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
	dbg.debug(r)
	return next.ServeHTTP(w, r)
}

func (dbg *RequestDebugger) debug(r *http.Request) {
	var requestID string
	var varMap map[string]interface{}
	var exists bool
	varMap, exists = r.Context().Value("vars").(map[string]interface{})
	dbg.logger.Debug(
		"vars pre",
		zap.Any("vars", varMap),
		zap.Any("exists", exists),
	)

	if !exists {
		ctx := context.WithValue(r.Context(), "vars", make(map[string]interface{}))
		r = r.WithContext(ctx)
		dbg.logger.Debug("xxx")
		varMap = ctx.Value("vars").(map[string]interface{})
	}

	dbg.logger.Debug(
		"vars state",
		zap.Any("vars", varMap),
	)

	if v, exists := varMap["request_id"]; exists {
		requestID = v.(string)
	} else {
		requestID = uuid.NewV4().String()
		varMap["request_id"] = requestID
	}

	/*

		varsCtx := r.Context().Value("vars")
		if varsCtx == nil {
			ctx := context.WithValue(r.Context(), "vars", make(map[string]interface{}))
			r = r.WithContext(ctx)
		}
		vars := r.Context().Value("vars").(map[string]interface{})

		if v, exists := vars["request_id"]; exists {
			requestID = v.(string)
		} else {
			requestID = uuid.NewV4().String()
			vars["request_id"] = requestID
			// TODO: inject request ID into the context
			//ctx := context.WithValue(r.Context(), "vars", vars)
			// r = r.WithContext(ctx)
		}
	*/

	dbg.logger.Debug("request debugging",
		zap.Any("request_id", requestID),
		zap.String("tag", dbg.Tag),
		zap.String("method", r.Method),
		zap.String("uri", r.RequestURI),
		zap.Any("vars", varMap),
	)

	/*
		if reqDump, err := httputil.DumpRequest(r, true); err == nil {
			dbg.logger.Debug(fmt.Sprintf("request: %s", reqDump))
		}
	*/
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
