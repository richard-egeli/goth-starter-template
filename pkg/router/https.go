package router

import "net/http"

func HTTPSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.TLS == nil && r.Header.Get("X-Forwarded-Proto") != "https" {
			redirectURL := "https://" + r.Host + r.RequestURI
			http.Redirect(w, r, redirectURL, http.StatusMovedPermanently)
			return
		}

		next.ServeHTTP(w, r)
	})
}
