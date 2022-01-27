package proxy

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.skymeyer.dev/app"
	"go.uber.org/zap"

	"go.reefassistant.com/apex-proxy/pkg/logger"
)

var (
	XMLEndpoints = []string{
		"/cgi-bin/datalog.xml",
		"/cgi-bin/outlog.xml",
		"/cgi-bin/status.xml",
	}
	JSONEndpoints = []string{
		"/cgi-bin/datalog.json",
		"/cgi-bin/outlog.json",
		"/cgi-bin/status.json",
	}
	localhostRanges = []string{
		"127.0.0.0/8",
		"::1/128",
	}
	privateRanges = []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
		"fc00::/7",  // RFC-4193 Unique Local Addresses
		"fec0::/10", // RFC-3879 Site Local Addresses (deprecated)
	}
)

// New creates a new Apex proxy handler.
func New(endpoint string, opts ...ProxyOption) (*Apex, error) {

	// Default proxy setup
	apex := &Apex{
		acesssListCacheSize: 1000,
		allowedURLs:         make(map[string]struct{}),
		allowedIPs:          localhostRanges,
	}

	// Apply options
	for _, opt := range opts {
		opt(apex)
	}

	// Parse Apex endpoint URL
	target, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	// Setup IP access list
	accessList, err := NewAccessList(apex.allowedIPs, apex.acesssListCacheSize)
	if err != nil {
		return nil, err
	}
	apex.accessList = accessList

	// Setup reverse proxy
	apex.proxy = httputil.NewSingleHostReverseProxy(target)
	apex.proxy.ModifyResponse = apex.modifyResponse
	apex.proxy.ErrorHandler = apex.errorHandler

	return apex, nil
}

// Apex proxy handler.
type Apex struct {
	proxy               *httputil.ReverseProxy
	accessList          *AccessList
	acesssListCacheSize int
	allowedIPs          []string
	allowedURLs         map[string]struct{}
}

// ServeHTTP implements http.Handler.
func (a *Apex) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log := logger.Context(r.Context())
	w.Header().Set("server", fmt.Sprintf("apex-proxy %s", app.Version))

	// Derive source IP from reported remote address, no proxy support
	source, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		http.Error(w, "forbidden", http.StatusForbidden)
		log.Error("no remote ip found")
		return
	}

	// Is source IP is allowed ?
	if !a.accessList.Contains(source) {
		http.Error(w, "forbidden", http.StatusForbidden)
		log.Warn("ip blocked", zap.String("source", source))
		return
	}
	log.Debug("ip allowed", zap.String("source", source))

	// Is method allowed ?
	if r.Method != "GET" {
		http.Error(w, "forbidden", http.StatusForbidden)
		log.Warn("method blocked", zap.String("method", r.Method))
		return
	}

	// Is URL allowed ?
	if _, ok := a.allowedURLs[r.URL.Path]; !ok {
		http.Error(w, "forbidden", http.StatusForbidden)
		log.Warn("url blocked", zap.String("url", r.URL.Path))
		return
	}
	log.Debug("url allowed", zap.String("url", r.URL.Path))

	// Handle the request through the reverse proxy
	a.proxy.ServeHTTP(w, r)
}

// httputil.ReverseProxy.ErrorHandler implementation
func (a *Apex) errorHandler(w http.ResponseWriter, r *http.Request, err error) {
	logger.Context(r.Context()).Error(
		"proxy error",
		zap.String("error", err.Error()),
	)
	rid, _ := logger.RequestIDFrom(r.Context())
	http.Error(w, fmt.Sprintf("request-id=%s", rid), http.StatusServiceUnavailable)
}

// httputil.ReverseProxy.ModifyResponse implementation
func (a *Apex) modifyResponse(resp *http.Response) error {
	// TODO: implement caching and/or add rate limiting per source ip
	// TODO: add data obfuscation like hiding serial number and software version
	return nil
}

// addURLs is a helper to initialize the allowed URLs into a map.
func (a *Apex) addURLs(urls []string) {
	for _, u := range urls {
		if u, err := url.Parse(u); err == nil {
			a.allowedURLs[u.Path] = struct{}{}
		}
	}
}
