package proxy

// ProxyOption allows configuring the Apex proxy.
type ProxyOption func(*Apex)

// WithAllowPrivateIPs allows requests from all private IP ranges.
func WithAllowPrivateIPs() ProxyOption {
	return func(a *Apex) {
		a.allowedIPs = append(a.allowedIPs, privateRanges...)
	}
}

// WithAllowIPs allows requests from given list of IP ranges.
func WithAllowIPs(ips []string) ProxyOption {
	return func(a *Apex) {
		a.allowedIPs = append(a.allowedIPs, ips...)
	}
}

// WithAllowXMLEndpoints allows GET requests for cgi-bin XML endpoints.
func WithAllowXMLEndpoints() ProxyOption {
	return func(a *Apex) {
		a.addURLs(XMLEndpoints)
	}
}

// WithAllowJSONEndpoints allows GET requests for cgi-bin JSON endpoints.
func WithAllowJSONEndpoints() ProxyOption {
	return func(a *Apex) {
		a.addURLs(JSONEndpoints)
	}
}

// WithAllowedURLs allows GET requests to given list of URLs.
func WithAllowedURLs(urls []string) ProxyOption {
	return func(a *Apex) {
		a.addURLs(urls)
	}
}

// WithAccessListCacheSize configures the LRU IP matcher cache size.
func WithAccessListCacheSize(size int) ProxyOption {
	return func(a *Apex) {
		a.acesssListCacheSize = size
	}
}
