package proxy

import (
	"github.com/valyala/fasthttp"
)

func newRetryFunc(c *Config) fasthttp.RetryIfErrFunc {
	return func(req *fasthttp.Request, attempt int, err error) (bool, bool) {
		if c.Retry > attempt {
			return true, true
		}

		return false, false
	}
}