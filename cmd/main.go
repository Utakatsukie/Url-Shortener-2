package main

import (
	"Url-Shortener-2/configs"
	"Url-Shortener-2/internal/auth"
	"Url-Shortener-2/internal/link"
	"Url-Shortener-2/middleware"
	"Url-Shortener-2/pkg/db"
	"fmt"
	"net/http"
)

func main() {
	conf := configs.LoadConfig()
	database := db.NewDb(conf)
	router := http.NewServeMux()

	// Repositories
	linkRepository := link.NewLinkRepository(database)

	// Handler
	//hello.NewHelloHandler(router)
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config: conf,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
	})

	server := http.Server{
		Addr:    ":8081",
		Handler: middleware.CORS(middleware.Logging(router)),
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
