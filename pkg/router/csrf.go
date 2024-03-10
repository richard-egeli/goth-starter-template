package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
)

func CSRFMiddleware(next http.Handler) http.Handler {
	secret := os.Getenv("CSRF_SECRET")
	runtime := os.Getenv("RUNTIME")

	if secret == "" {
		if runtime == "production" {
			fmt.Println()
			fmt.Println("=========== ERROR ===========")
			fmt.Println()
			fmt.Println("CSRF_SECRET env variable not found")
			fmt.Println()
			fmt.Println("=========== ERROR ===========")
			fmt.Println()
			os.Exit(1)
		}

		fmt.Println()
		fmt.Println("=========== WARNING ===========")
		fmt.Println()
		fmt.Println("CSRF_SECRET env variable not found, using hardcoded substitute but it is not safe")
		fmt.Println()
		fmt.Println("=========== WARNING ===========")
		fmt.Println()
		secret = "y3qOQoP4UZpux3limEVSp4FjyP48AxTd%"
	}

	var CSRF func(http.Handler) http.Handler
	if runtime == "production" {
		// Only accept HTTPS in production
		CSRF = csrf.Protect([]byte(secret), csrf.Secure(true))
	} else {
		CSRF = csrf.Protect([]byte(secret), csrf.Secure(false))
	}

	return CSRF(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := csrf.Token(r)

		w.Header().Set("X-CSRF-Token", token)
		next.ServeHTTP(w, r)
	}))
}
