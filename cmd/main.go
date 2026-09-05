package main

import (
	"Url-Shortener-2/configs"
	"Url-Shortener-2/internal/auth"
	"Url-Shortener-2/internal/link"
	"Url-Shortener-2/internal/stat"
	"Url-Shortener-2/internal/user"
	"Url-Shortener-2/middleware"
	"Url-Shortener-2/pkg/db"
	"Url-Shortener-2/pkg/event"
	"fmt"
	"log"
	"net/http"
)

func main() {
	conf, err := configs.LoadConfig()
	if err != nil {
		log.Fatal("Config initialisation failed: %v", err)
	}
	database := db.NewDb(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	// Repositories
	linkRepository := link.NewLinkRepository(database)
	userRepository := user.NewUserRepository(database)
	statRepository := stat.NewStatRepository(database)

	// Services
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepository,
	})

	// Handler
	//hello.NewHelloHandler(router)
	auth.NewAuthHandler(router, auth.AuthHandlerDeps{
		Config:      conf,
		AuthService: authService,
	})
	link.NewLinkHandler(router, link.LinkHandlerDeps{
		LinkRepository: linkRepository,
		EventBus:       eventBus,
		Config:         conf,
	})
	stat.NewStatHandler(router, stat.StatHandlerDeps{
		StatRepository: statRepository,
		Config:         conf,
	})

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)

	server := http.Server{
		Addr:    ":8081",
		Handler: stack(router),
	}

	go statService.AddClick()

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}
