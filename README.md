# caddy-request-debug

<a href="https://github.com/greenpau/caddy-request-debug/actions/" target="_blank"><img src="https://github.com/greenpau/caddy-request-debug/workflows/build/badge.svg?branch=main"></a>
<a href="https://pkg.go.dev/github.com/greenpau/caddy-request-debug" target="_blank"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"></a>
<a href="https://caddy.community" target="_blank"><img src="https://img.shields.io/badge/community-forum-ff69b4.svg"></a>

Request Debugging Middleware Plugin for [Caddy v2](https://github.com/caddyserver/caddy).

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

TBD.
