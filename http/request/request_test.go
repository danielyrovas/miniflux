// Copyright 2018 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package request // import "miniflux.app/http/request"

import (
	"net/http"
	"testing"
)

func TestRealIPWithoutHeaders(t *testing.T) {
	r := &http.Request{RemoteAddr: "192.168.0.1:4242"}
	if ip := RealIP(r); ip != "192.168.0.1" {
		t.Fatalf(`Unexpected result, got: %q`, ip)
	}

	r = &http.Request{RemoteAddr: "192.168.0.1"}
	if ip := RealIP(r); ip != "192.168.0.1" {
		t.Fatalf(`Unexpected result, got: %q`, ip)
	}
}

func TestRealIPWithXFFHeader(t *testing.T) {
	// Test with multiple IPv4 addresses.
	headers := http.Header{}
	headers.Set("X-Forwarded-For", "203.0.113.195, 70.41.3.18, 150.172.238.178")
	r := &http.Request{RemoteAddr: "192.168.0.1:4242", Header: headers}

	if ip := RealIP(r); ip != "203.0.113.195" {
		t.Fatalf(`Unexpected result, got: %q`, ip)
	}

	// Test with single IPv6 address.
	headers = http.Header{}
	headers.Set("X-Forwarded-For", "2001:db8:85a3:8d3:1319:8a2e:370:7348")
	r = &http.Request{RemoteAddr: "192.168.0.1:4242", Header: headers}

	if ip := RealIP(r); ip != "2001:db8:85a3:8d3:1319:8a2e:370:7348" {
		t.Fatalf(`Unexpected result, got: %q`, ip)
	}

	// Test with single IPv4 address.
	headers = http.Header{}
	headers.Set("X-Forwarded-For", "70.41.3.18")
	r = &http.Request{RemoteAddr: "192.168.0.1:4242", Header: headers}

	if ip := RealIP(r); ip != "70.41.3.18" {
		t.Fatalf(`Unexpected result, got: %q`, ip)
	}

	// Test with invalid IP address.
	headers = http.Header{}
	headers.Set("X-Forwarded-For", "fake IP")
	r = &http.Request{RemoteAddr: "192.168.0.1:4242", Header: headers}

	if ip := RealIP(r); ip != "192.168.0.1" {
		t.Fatalf(`Unexpected result, got: %q`, ip)
	}
}

func TestRealIPWithXRealIPHeader(t *testing.T) {
	headers := http.Header{}
	headers.Set("X-Real-Ip", "192.168.122.1")
	r := &http.Request{RemoteAddr: "192.168.0.1:4242", Header: headers}

	if ip := RealIP(r); ip != "192.168.122.1" {
		t.Fatalf(`Unexpected result, got: %q`, ip)
	}
}

func TestRealIPWithBothHeaders(t *testing.T) {
	headers := http.Header{}
	headers.Set("X-Forwarded-For", "203.0.113.195, 70.41.3.18, 150.172.238.178")
	headers.Set("X-Real-Ip", "192.168.122.1")

	r := &http.Request{RemoteAddr: "192.168.0.1:4242", Header: headers}

	if ip := RealIP(r); ip != "203.0.113.195" {
		t.Fatalf(`Unexpected result, got: %q`, ip)
	}
}
