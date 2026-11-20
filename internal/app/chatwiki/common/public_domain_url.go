// Copyright © 2016- 2025 Wuhan Sesame Small Customer Service Network Technology Co., Ltd.

package common

import (
	"chatwiki/internal/app/chatwiki/i18n"
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	neturl "net/url"
	"strconv"
	"strings"
	"time"
)

const publicDomainURLDNSLookupTimeout = time.Second
const publicDomainURLDefaultTimeout = 60 * time.Second

var nonPublicDomainURLNetworks = func() []*net.IPNet {
	cidrs := []string{
		`0.0.0.0/8`, `10.0.0.0/8`, `100.64.0.0/10`, `127.0.0.0/8`,
		`169.254.0.0/16`, `172.16.0.0/12`, `192.0.0.0/24`, `192.0.2.0/24`,
		`192.88.99.0/24`, `192.168.0.0/16`, `198.18.0.0/15`, `198.51.100.0/24`, `203.0.113.0/24`,
		`224.0.0.0/4`, `240.0.0.0/4`,
		// Exclude special protocol assignments, transition mechanisms and documentation ranges.
		`2001::/23`, `2001:db8::/32`, `2002::/16`, `3fff::/20`,
	}
	networks := make([]*net.IPNet, 0, len(cidrs))
	for _, cidr := range cidrs {
		_, network, _ := net.ParseCIDR(cidr)
		networks = append(networks, network)
	}
	return networks
}()

var lookupPublicDomainURLIPAddr = func(ctx context.Context, host string) ([]net.IPAddr, error) {
	return net.DefaultResolver.LookupIPAddr(ctx, host)
}

// ValidatePublicDomainURL ensures a URL uses a domain that resolves only to public IP addresses.
func ValidatePublicDomainURL(lang, rawURL string) error {
	parsed, err := neturl.ParseRequestURI(strings.TrimSpace(rawURL))
	if err != nil || parsed == nil {
		return errors.New(i18n.Show(lang, `request_link_empty_or_error`))
	}
	if (parsed.Scheme != `http` && parsed.Scheme != `https`) || parsed.Host == `` || parsed.User != nil {
		return errors.New(i18n.Show(lang, `request_link_domain_only`))
	}
	host := parsed.Hostname()
	if port := parsed.Port(); port != `` {
		number, err := strconv.Atoi(port)
		if err != nil || number < 1 || number > 65535 {
			return errors.New(i18n.Show(lang, `request_link_empty_or_error`))
		}
	}
	if _, err := lookupPublicDomainURLIPs(context.Background(), lang, host); err != nil {
		return fmt.Errorf(`%s:%w`, rawURL, err)
	}
	return nil
}

func isNonPublicDomainURLIP(ip net.IP) bool {
	if ip == nil || !ip.IsGlobalUnicast() {
		return true
	}
	// Only native IPv6 global unicast (2000::/3) is allowed. This also
	// excludes NAT64 prefixes that could translate to private IPv4 addresses.
	if ip.To4() == nil && ip[0]&0xe0 != 0x20 {
		return true
	}
	for _, network := range nonPublicDomainURLNetworks {
		if network.Contains(ip) {
			return true
		}
	}
	return false
}

func lookupPublicDomainURLIPs(ctx context.Context, lang, host string) ([]net.IPAddr, error) {
	if host == `` || net.ParseIP(host) != nil {
		return nil, errors.New(i18n.Show(lang, `request_link_domain_only`))
	}
	ctx, cancel := context.WithTimeout(ctx, publicDomainURLDNSLookupTimeout)
	defer cancel()
	addresses, err := lookupPublicDomainURLIPAddr(ctx, host)
	if err != nil {
		return nil, errors.New(i18n.Show(lang, `request_link_dns_lookup_failed`))
	}
	if len(addresses) == 0 {
		return nil, errors.New(i18n.Show(lang, `request_link_dns_no_addresses`))
	}
	for _, address := range addresses {
		if address.Zone != `` || isNonPublicDomainURLIP(address.IP) {
			return nil, errors.New(i18n.Show(lang, `request_link_non_public_address`, address.IP.String()))
		}
	}
	return addresses, nil
}

// NewPublicDomainTransport creates a transport that validates every DNS result
// immediately before dialing and connects to the validated IP directly.
func NewPublicDomainTransport(lang string, timeout time.Duration, proxyURL string) http.RoundTripper {
	if timeout <= 0 {
		timeout = publicDomainURLDefaultTimeout
	}
	dialer := &net.Dialer{Timeout: timeout, KeepAlive: 2 * timeout}
	proxyAddress := publicDomainURLProxyAddress(proxyURL)
	return &http.Transport{
		Proxy: publicDomainURLProxy(proxyURL),
		DialContext: func(ctx context.Context, network, address string) (net.Conn, error) {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			if proxyAddress != `` && address == proxyAddress {
				return dialer.DialContext(ctx, network, address)
			}
			return dialPublicDomainURL(ctx, lang, network, address, dialer.DialContext)
		},
		ResponseHeaderTimeout: timeout,
		TLSHandshakeTimeout:   timeout,
		MaxIdleConnsPerHost:   100,
		DisableKeepAlives:     true,
	}
}

func publicDomainURLProxy(rawProxyURL string) func(*http.Request) (*neturl.URL, error) {
	proxyURL, err := parsePublicDomainURLProxy(rawProxyURL)
	if err != nil {
		return func(*http.Request) (*neturl.URL, error) {
			return nil, fmt.Errorf(`invalid workflow HTTP proxy URL`)
		}
	}
	if proxyURL == nil {
		return nil
	}
	return http.ProxyURL(proxyURL)
}

func parsePublicDomainURLProxy(rawProxyURL string) (*neturl.URL, error) {
	rawProxyURL = strings.TrimSpace(rawProxyURL)
	if rawProxyURL == `` {
		return nil, nil
	}
	proxyURL, err := neturl.Parse(rawProxyURL)
	if err != nil || (proxyURL.Scheme != `http` && proxyURL.Scheme != `https`) || proxyURL.Host == `` {
		return nil, errors.New(`invalid workflow HTTP proxy URL`)
	}
	return proxyURL, nil
}

func publicDomainURLProxyAddress(rawProxyURL string) string {
	proxyURL, err := parsePublicDomainURLProxy(rawProxyURL)
	if err != nil || proxyURL == nil {
		return ``
	}
	port := proxyURL.Port()
	if port == `` {
		port = `80`
		if proxyURL.Scheme == `https` {
			port = `443`
		}
	}
	return net.JoinHostPort(proxyURL.Hostname(), port)
}

func dialPublicDomainURL(ctx context.Context, lang, network, address string, dial func(context.Context, string, string) (net.Conn, error)) (net.Conn, error) {
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	addresses, err := lookupPublicDomainURLIPs(ctx, lang, host)
	if err != nil {
		return nil, err
	}
	var lastErr error
	for index, resolved := range addresses {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		dialCtx := ctx
		cancel := func() {}
		if deadline, ok := ctx.Deadline(); ok {
			remaining := time.Until(deadline)
			if remaining <= 0 {
				return nil, context.DeadlineExceeded
			}
			dialCtx, cancel = context.WithTimeout(ctx, remaining/time.Duration(len(addresses)-index))
		}
		// Pass an IP literal to the system dialer so it cannot resolve the domain again.
		conn, err := dial(dialCtx, network, net.JoinHostPort(resolved.IP.String(), port))
		cancel()
		if err == nil {
			return conn, nil
		}
		lastErr = err
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return nil, lastErr
}
