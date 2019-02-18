package webserver

import (
	"net/http"
	"time"

	"github.com/zekroTJA/timedmap"
	"golang.org/x/time/rate"
)

const keyExpireTime = 10 * time.Minute

type LimiterOpts struct {
	rate  int
	burst int
}

type RateLimiter struct {
	limiters        *timedmap.TimedMap
	limits          map[string]*LimiterOpts
	stdLimit        *LimiterOpts
	disallowHandler func(w http.ResponseWriter, req *http.Request)
}

func NewRateLimiter(stdLimiterOpts *LimiterOpts, disallowHandler func(w http.ResponseWriter, req *http.Request)) *RateLimiter {
	return &RateLimiter{
		limiters:        timedmap.New(1 * time.Minute),
		limits:          make(map[string]*LimiterOpts),
		stdLimit:        stdLimiterOpts,
		disallowHandler: disallowHandler,
	}
}

func (r *RateLimiter) Register(ident string, rate, burst int) *RateLimiter {
	r.limits[ident] = &LimiterOpts{rate, burst}
	return r
}

func (r *RateLimiter) Check(ident string, w http.ResponseWriter, req *http.Request) bool {
	var limiter *rate.Limiter

	_limiter := r.limiters.GetValue(req.RemoteAddr)
	if _limiter == nil {
		limiter = r.createLimiter(ident)
		r.limiters.Set(req.RemoteAddr, limiter, keyExpireTime)
	} else {
		var ok bool
		limiter, ok = _limiter.(*rate.Limiter)
		if !ok {
			limiter = r.createLimiter(ident)
			r.limiters.Set(req.RemoteAddr, limiter, keyExpireTime)
		}
	}

	allowed := limiter.Allow()
	if !allowed {
		if r.disallowHandler != nil {
			r.disallowHandler(w, req)
		}
	}

	r.limiters.Refresh(req.RemoteAddr, keyExpireTime)

	return allowed
}

func (r *RateLimiter) createLimiter(ident string) *rate.Limiter {
	limOpt, ok := r.limits[ident]
	if !ok {
		limOpt = r.stdLimit
	}
	return rate.NewLimiter(rate.Limit(limOpt.rate), limOpt.burst)
}
