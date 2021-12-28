package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContains(t *testing.T) {

	set, err := NewAccessList([]string{
		"1.0.0.0/8",
		"2.2.2.0/24",
		"3.3.3.3/32",
		"2002:0:0:1200::/56",
		"2604:ca00:108:b846::e61:f746/128",
	}, 100)
	require.NoError(t, err)

	for _, d := range []struct {
		ip  string
		exp bool
	}{
		{
			ip:  "1.2.3.4",
			exp: true,
		},
		{
			ip:  "2.2.2.2",
			exp: true,
		},
		{
			ip:  "3.3.3.3",
			exp: true,
		},
		{
			ip:  "4.4.4.4",
			exp: false,
		},
		{
			ip:  "2002::1234:abcd:ffff:c0a8:101",
			exp: true,
		},
		{
			ip:  "2002::1334:abcd:ffff:c0a8:101",
			exp: false,
		},
		{
			ip:  "2604:ca00:108:b846::e61:f746",
			exp: true,
		},
		{
			ip:  "2604:ca00:108:b846::e61:f747",
			exp: false,
		},
	} {
		assert.Equal(t, d.exp, set.Contains(d.ip))

		cached, err := set.fromCache(d.ip)
		assert.NoError(t, err)
		assert.Equal(t, d.exp, cached)
	}
}
