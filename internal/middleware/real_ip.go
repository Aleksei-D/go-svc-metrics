package middleware

import (
	"net"
	"net/http"
)

type RealIPMiddleware struct {
	network *net.IPNet
}

func NewRealIPMiddleware(trustSubnet string) (RealIPMiddleware, error) {
	realIPMiddleware := RealIPMiddleware{}
	_, network, err := net.ParseCIDR(trustSubnet)
	if err != nil {
		return realIPMiddleware, err
	}
	realIPMiddleware.network = network

	return RealIPMiddleware{network: network}, nil
}

func (m *RealIPMiddleware) GetRealIPMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if m.network == nil {
			next.ServeHTTP(w, r)
		}

		ipStr := r.Header.Get("X-Real-IP")
		if ipStr != "" {
			ip := net.ParseIP(ipStr)
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
