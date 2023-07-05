package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Faner201/sc_links/internal/auth"
	"github.com/Faner201/sc_links/internal/config"
	"github.com/Faner201/sc_links/internal/db"
	"github.com/Faner201/sc_links/internal/github"
	"github.com/Faner201/sc_links/internal/server"
	"github.com/Faner201/sc_links/internal/shorten"
	"github.com/Faner201/sc_links/internal/storage/shortening"
	"github.com/Faner201/sc_links/internal/storage/user"
)

func main() {
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer dbCancel()

	err := config.Init()
	if err != nil {
		log.Fatal(err)
	}

	mgoClient, err := db.Connect(dbCtx, config.Get().DB.URI)
	if err != nil {
		log.Fatal(err)
	}

	mgoDB := mgoClient.Client().Database(config.Get().DB.Database)

	var (
		shorteningStorage = shortening.NewMongoDB(mgoDB)
		userStorage       = user.NewMongoDB(mgoDB)
		shortener         = shorten.NewService(shorteningStorage)
		githubClient      = github.NewClien()
		authenticator     = auth.NewService(
			githubClient,
			userStorage,
			config.Get().Github.ClientID,
			config.Get().Github.ClientSecret,
		)
		srv = server.New(shortener, authenticator)
	)

	srv.AddCloser(mgoClient.Close)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := http.ListenAndServe(config.Get().ListenAddr(), srv); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error running server: %v", err)
		}
	}()

	log.Println("server started")
	<-quit

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("error closing server: %v", err)
	}

	log.Println("server stopped")

}
