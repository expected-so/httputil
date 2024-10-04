package httputil

import (
	"net"
	"net/http"
)

func GetIP(r *http.Request) string {
	if value := r.Header.Get("cf-connecting-ip"); value != "" {
		return FormatIP(value)
	}
	return FormatIP(r.RemoteAddr)
}

func FormatIP(raw string) string {
	if raw == "::1" {
		return "127.0.0.1"
	}

	host, _, err := net.SplitHostPort(raw)
	if err == nil {
		return host
	}
	return raw
}
