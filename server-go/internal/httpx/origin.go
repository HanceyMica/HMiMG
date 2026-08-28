package httpx

import "net/http"

func Origin(r *http.Request, trustProxy bool) string {
	proto := "http"
	host := r.Host
	if trustProxy {
		if v := r.Header.Get("X-Forwarded-Proto"); v != "" {
			proto = v
		}
		if v := r.Header.Get("X-Forwarded-Host"); v != "" {
			host = v
		}
	}
	if proto == "http" && r.TLS != nil {
		proto = "https"
	}
	return proto + "://" + host
}
