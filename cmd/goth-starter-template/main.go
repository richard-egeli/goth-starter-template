package main

import (
	"fmt"
	"net/http"
	"os"

	"goth-starter-template/pkg/router"
	"goth-starter-template/views/layout"
	"goth-starter-template/views/pages"

	"github.com/joho/godotenv"
)

func defaultMW(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if r.URL.Path == "/" {
				next.ServeHTTP(w, r)
			} else if r.URL.Path != "/404" {
				fmt.Println(r.URL.Path)
				http.Redirect(w, r, "/404", http.StatusFound)
			}
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println(err.Error())
	}

	port := os.Getenv("PORT")
	runtime := os.Getenv("RUNTIME")

	if port == "" {
		port = "8080"
	}

	layoutData := layout.BaseLayoutData{
		Title: "go-starter-template",
	}

	base := router.New()
	base.Use(router.GzipMiddleware)
	base.Use(router.CSRFMiddleware)

	if runtime != "production" {
		base.SetupBrowserRefreshEvent()
		base.Dir("/scripts/", "./web/src", router.TypescriptTranspilationMiddleware)
	} else {
		// Expects in production that typescript should be transpiled into a javascript file
		base.Dir("/scripts/", "./scripts", router.GzipMiddleware)
		base.Use(router.HTTPSMiddleware)
	}

	base.Dir("/static/", "./static", router.GzipMiddleware)
	base.Get("/", defaultMW, router.Page(layout.BaseLayout, &layoutData))
	base.Get("/404", router.Page(pages.NotFound, nil))

	fmt.Println("Listening on port:", port)
	base.Listen(port)
}
