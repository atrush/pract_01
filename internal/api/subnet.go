package api

import (
	"net"
	"net/http"
	"strings"
)

type Subnet struct {
	AllowMasks []net.IPNet
}

func NewSubnet(cfg string) *Subnet {
	return &Subnet{
		AllowMasks: readMasks(cfg),
	}
}

// Middleware checks that request ip in allowed network
func (s *Subnet) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		strIP := r.Header.Get("X-Real-IP")
		ip := net.ParseIP(strIP)

		if ip != nil && len(s.AllowMasks) != 0 {
			for _, mask := range s.AllowMasks {
				if mask.Contains(ip) {

					next.ServeHTTP(w, r)
					return
				}
			}
		}

		http.Error(w, "not allowed", http.StatusForbidden)
	})
}

//  readMasks reads ip network masks from config string
func readMasks(cfg string) []net.IPNet {
	strMasks := strings.Split(cfg, ",")
	if len(strMasks) == 0 {
		return nil
	}

	masks := []net.IPNet{}
	for _, m := range strMasks {
		m = strings.ReplaceAll(m, " ", "")
		_, ipnet, err := net.ParseCIDR(m)
		if err == nil && ipnet != nil {
			masks = append(masks, *ipnet)
		}
	}

	return masks
}
