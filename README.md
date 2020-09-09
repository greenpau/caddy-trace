# caddy-request-debug

<a href="https://github.com/greenpau/caddy-request-debug/actions/" target="_blank"><img src="https://github.com/greenpau/caddy-request-debug/workflows/build/badge.svg?branch=main"></a>
<a href="https://pkg.go.dev/github.com/greenpau/caddy-request-debug" target="_blank"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"></a>
<a href="https://caddy.community" target="_blank"><img src="https://img.shields.io/badge/community-forum-ff69b4.svg"></a>

Request Debugging Middleware Plugin for [Caddy v2](https://github.com/caddyserver/caddy).

<!-- begin-markdown-toc -->
## Table of Contents

* [Overview](#overview)
* [Getting Started](#getting-started)

<!-- end-markdown-toc -->

## Overview

The plugin is a middleware which displays the content of the request it
handles. It helps troubleshooting web requests by exposing headers
(e.g. cookies), URL parameters, etc.

For background, the idea for the creation of plugin came during a
development of another plugin which rewrites headers of web requests.
There was a need to compare "before and after" content of the request.

The `request_debug` directive gets inserted prior to the plugin which
modifies a request and immediately after it. The log with the content
of the request show up twice and it is easy to compare the two.

## Getting Started

Add `request_debug` handler to enable this plugin.

The `disabled=yes` argument disables the operation of the plugin.

The `tag` argument injects the value in the log output. This way, one can have
multiple handlers and there is a way to deferentiate between them.

The `response_debug` argument instructs the plugin to buffer
responses and log response related metadata, i.e. status codes, length, etc.

When a request arrives for `/version`, the plugin will be triggered two (2)
times. The first handler is disables. The two (2) other handlers will trigger
with different tags. The `respond` handler is terminal and it means the handler
with `marvel` tag will not trigger.

When a request arrives for `/whoami`, the plugin will be triggered three (2)
times because `respond /version` will not terminate the handling of the plugin.
Notably, the plugin will output response metadata due to the presence of
`response_debug` argument.

```
{
  http_port     9080
  https_port    9443
}

localhost:9080 {
  route {
    request_debug disabled=yes
    request_debug disabled=no tag="foo"
    request_debug disabled=no tag="bar"
    respond /version "1.0.0" 200
    request_debug tag="marvel" response_debug=yes
    respond /whoami 200 {
      body "greenpau"
    }
  }
}
```

The same JSON configuration:

```json
{
  "apps": {
    "http": {
      "http_port": 9080,
      "https_port": 9443,
      "servers": {
        "srv0": {
          "listen": [
            ":9080"
          ],
          "routes": [
            {
              "handle": [
                {
                  "handler": "subroute",
                  "routes": [
                    {
                      "handle": [
                        {
                          "handler": "subroute",
                          "routes": [
                            {
                              "handle": [
                                {
                                  "disabled": true,
                                  "handler": "request_debug"
                                }
                              ]
                            },
                            {
                              "handle": [
                                {
                                  "handler": "request_debug",
                                  "tag": "foo"
                                }
                              ]
                            },
                            {
                              "handle": [
                                {
                                  "handler": "request_debug",
                                  "tag": "bar"
                                }
                              ]
                            },
                            {
                              "handle": [
                                {
                                  "body": "1.0.0",
                                  "handler": "static_response",
                                  "status_code": 200
                                }
                              ],
                              "match": [
                                {
                                  "path": [
                                    "/version"
                                  ]
                                }
                              ]
                            },
                            {
                              "handle": [
                                {
                                  "handler": "request_debug",
                                  "response_debug_enabled": true,
                                  "tag": "marvel"
                                }
                              ]
                            },
                            {
                              "handle": [
                                {
                                  "body": "greenpau",
                                  "handler": "static_response",
                                  "status_code": 200
                                }
                              ],
                              "match": [
                                {
                                  "path": [
                                    "/whoami*"
                                  ]
                                }
                              ]
                            }
                          ]
                        }
                      ]
                    }
                  ]
                }
              ],
              "match": [
                {
                  "host": [
                    "localhost"
                  ]
                }
              ],
              "terminal": true
            }
          ]
        }
      }
    }
  }
}
```
