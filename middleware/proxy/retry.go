package proxy

import "github.com/valyala/fasthttp"

type RetryIf func(req *fasthttp.Request, res *fasthttp.Response, err error) bool