package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	mu sync.Mutex
	visitors map[string] int
	limit int
	resetTime time.Duration
}

func NewRateLimiter(limit int, resetTime time.Duration) *rateLimiter{
	rl := &rateLimiter{
		visitors: make(map[string]int),
		limit:limit,
		resetTime: resetTime,
	}
	// start the reset routine
	go rl.resetVisitorCount()
	return rl
}

func (rl *rateLimiter) resetVisitorCount(){
	for{
		time.Sleep(rl.resetTime)
		rl.mu.Lock()
		rl.visitors = make(map[string]int)
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) Middleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		rl.mu.Lock()
		defer rl.mu.Unlock()

		visitorIP := r.RemoteAddr // you might want to extract the IP in a more sophisticated way
		rl.visitors[visitorIP]++
		fmt.Printf("Visitor count from %v is %v\n",visitorIP,rl.visitors[visitorIP])

		if rl.visitors[visitorIP] > rl.limit {
			http.Error(w, "Too many requests", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w,r)
	})
}