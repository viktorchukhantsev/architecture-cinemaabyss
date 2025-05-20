package strangler

import (
	"log"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type config struct {
	percentage int
}

type StranglerRouter struct {
	oldServiceProxy *httputil.ReverseProxy
	newServiceProxy *httputil.ReverseProxy
	config          config
}

func NewStranglerRouter(oldURL, newURL string, percentage int) *StranglerRouter {
	old, err := url.Parse(oldURL)
	if err != nil {
		log.Fatalf("unable to parse URL: %s", oldURL)
	}

	new, err := url.Parse(newURL)
	if err != nil {
		log.Fatalf("unable to parse URL: %s", newURL)
	}

	proxyMonolith := httputil.NewSingleHostReverseProxy(old)
	proxyMovies := httputil.NewSingleHostReverseProxy(new)

	return &StranglerRouter{
		oldServiceProxy: proxyMonolith,
		newServiceProxy: proxyMovies,
		config:          config{percentage: percentage},
	}
}

func (p *StranglerRouter) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("received request: %s %s %s\n", r.URL.Scheme, r.Host, r.URL.String())
	if shouldRouteToNewSystem(p.config.percentage) {
		p.newServiceProxy.ServeHTTP(w, r)
	} else {
		p.oldServiceProxy.ServeHTTP(w, r)
	}
}

func shouldRouteToNewSystem(percentage int) bool {
	if percentage >= 100 {
		return true
	}
	if percentage <= 0 {
		return false
	}
	return randInt(1, 100) <= percentage
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}
