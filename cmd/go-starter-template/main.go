package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/richard-egeli/go-starter-template/pkg/router"
	"github.com/richard-egeli/go-starter-template/views/layout"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(err.Error())
	}

	layoutData := layout.BaseLayoutData{
		Title: "go-starter-template",
	}

	base := router.New()
	base.SetupBrowserRefreshEvent()
	base.Dir("/static/", "./static", nil)
	base.Dir("/scripts/", "./web/src", []router.Middleware{router.TypescriptTranspilationMiddleware})
	base.Get("/", nil, router.Page(layout.BaseLayout, &layoutData))
	base.Listen("8080")
}
