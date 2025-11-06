package middleware

import (
	"net"
	"net/http"
)

type RealIPMiddleware struct {
	network *net.IPNet
}

func NewRealIPMiddleware(network *net.IPNet) RealIPMiddleware {
	return RealIPMiddleware{network: network}
}

func (m *RealIPMiddleware) GetRealIPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.network == nil {
			next.ServeHTTP(w, r)
		}

		ipSTR := r.Header.Get("X-Real-IP")
		if ipSTR != "" {
			ip := net.ParseIP(ipSTR)
			if ip == nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if !m.network.Contains(ip) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}
