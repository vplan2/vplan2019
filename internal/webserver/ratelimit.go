package webserver

import (
	"net/http"
	"time"

	"github.com/zekroTJA/timedmap"
	"golang.org/x/time/rate"
)

const keyExpireTime = 10 * time.Minute

// LimiterOpts contains the rate, which is the ammount
// of tokens regenerated per second and burst, which is
// the initial and total ammount of tokens available
// in a limiters bucket
type LimiterOpts struct {
	rate  float64
	burst int
}

// RateLimiter maintains rate limiting for general requests
// containing a timed map with a limiter for each RemoteAddress,
// the set limit options for each identifier, the standard
// limit and the handler for disallowed
type RateLimiter struct {
	limiters        *timedmap.TimedMap
	limits          map[string]*LimiterOpts
	stdLimit        *LimiterOpts
	disallowHandler func(w http.ResponseWriter, req *http.Request)
}

// NewRateLimiter returns a new instance of RateLimiter
//   stdLimiterOpts  : standard limiter options used when no options
//                     were set to the passed identifier
//   disallowHandler : Handler which will be called when the request
//                     was disallowed because of rate limit transgression
func NewRateLimiter(stdLimiterOpts *LimiterOpts, disallowHandler func(w http.ResponseWriter, req *http.Request)) *RateLimiter {
	return &RateLimiter{
		limiters:        timedmap.New(1 * time.Minute),
		limits:          make(map[string]*LimiterOpts),
		stdLimit:        stdLimiterOpts,
		disallowHandler: disallowHandler,
	}
}

// Register registers new limit options for an ident specified endpoint
//   ident : an identifier for the limit, i.e. the name of the endpoint
//   rate  : the ammount of tokens which will be generated per second
//   burst : the maximum and initial ammount of tokens available in a bucket
func (r *RateLimiter) Register(ident string, rate float64, burst int) *RateLimiter {
	r.limits[ident] = &LimiterOpts{rate, burst}
	return r
}

// Check checks if the request can be allowed, which consumes a token from
// the limietrs token bucket. If the bukket is empty, the reuqest will be
// disallowed, the disallowedHandler will be called and this function
// returns false.
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

// createLimiter creates a new instance of limiter with options
// taken by the matched ident. If the ident was not registered,
// the default options will be passed to the new limiter.
func (r *RateLimiter) createLimiter(ident string) *rate.Limiter {
	limOpt, ok := r.limits[ident]
	if !ok {
		limOpt = r.stdLimit
	}
	return rate.NewLimiter(rate.Limit(limOpt.rate), limOpt.burst)
}
