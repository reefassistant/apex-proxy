package proxy

import (
	"fmt"
	"net"

	lru "github.com/hashicorp/golang-lru"
)

// NewAccessList creates a new IP access list for given CIDRs. A LRU cache is used
// when cacheSize is specified to speedup frequently looked up addresses.
func NewAccessList(cidr []string, cacheSize int) (*AccessList, error) {

	var al = &AccessList{}

	// Parse CIDRs into ipnets
	for _, in := range cidr {
		_, ipnet, err := net.ParseCIDR(in)
		if err != nil {
			return nil, err
		}
		al.ipnets = append(al.ipnets, ipnet)
	}

	// Initialize LRU cache
	if cacheSize > 0 {
		cache, err := lru.New(cacheSize)
		if err != nil {
			return nil, err
		}
		al.cache = cache
	}

	return al, nil
}

// AccessList represents an IP access list.
type AccessList struct {
	ipnets []*net.IPNet
	cache  *lru.Cache
}

// Contains verifies if given IP address is member of any defined access list nets.
func (a *AccessList) Contains(in string) bool {
	if cached, err := a.fromCache(in); err == nil {
		return cached
	}

	ip := net.ParseIP(in)
	for _, ipnet := range a.ipnets {
		if ipnet.Contains(ip) {
			return a.toCache(in, true)
		}
	}
	return a.toCache(in, false)
}

func (a *AccessList) toCache(in string, value bool) bool {
	if a.cache != nil {
		a.cache.Add(in, value)
	}
	return value
}

func (a *AccessList) fromCache(in string) (bool, error) {
	if a.cache != nil {
		if cached, ok := a.cache.Get(in); ok {
			return cached.(bool), nil
		}
	}
	return false, fmt.Errorf("missing cache value")
}
