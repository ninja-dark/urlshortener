package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"urlsortener/internal/infrastructure/api/router"
	"urlsortener/internal/infrastructure/db/pgstore"
	"urlsortener/internal/usecases/app"
	"urlsortener/internal/infrastructure/api/server"
)


func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	ust, err := pgstore.NewStore("postgress://test:1234@localhost:5432/urlshortener")
	if err != nil{
		log.Fatalln(err)
	}
	app := app.NewUrls(ust)
	rt := router.NewRouter(app)
	srv := server.NewServer(":8000", rt)
	srv.Start()
	log.Print("Start")

	<-ctx.Done()
	cancel()
	srv.Stop()
	log.Print("Exit")
}