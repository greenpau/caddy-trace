# caddy-trace

<a href="https://github.com/greenpau/caddy-trace/actions/" target="_blank"><img src="https://github.com/greenpau/caddy-trace/workflows/build/badge.svg?branch=main"></a>
<a href="https://pkg.go.dev/github.com/greenpau/caddy-trace" target="_blank"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"></a>
<a href="https://caddy.community" target="_blank"><img src="https://img.shields.io/badge/community-forum-ff69b4.svg"></a>
<a href="https://caddyserver.com/docs/modules/http.handlers.trace" target="_blank"><img src="https://img.shields.io/badge/caddydocs-trace-green.svg"></a>

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

The `trace` directive gets inserted prior to the plugin which
modifies a request and immediately after it. The log with the content
of the request show up twice and it is easy to compare the two.

## Getting Started

Add `trace` handler to enable this plugin.

The `disabled=yes` argument disables the operation of the plugin.

The `tag` argument injects the value in the log output. This way, one can have
multiple handlers and there is a way to deferentiate between them.

The `response_debug` argument instructs the plugin to buffer
responses and log response related metadata, i.e. status codes, length, etc.

The `uri_filter` directive instructs the plugin to intercepts only
the requests with the URI matching the regular expression in the filter.

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
    trace disabled=yes
    trace disabled=no tag="foo"
    trace disabled=no tag="bar"
    respond /version "1.0.0" 200
    trace tag="marvel" response_debug=yes
    trace tag="custom" response_debug=yes uri_filter="^/whoami$"
    respond /whoami 200 {
      body "greenpau"
    }
  }
}
```
