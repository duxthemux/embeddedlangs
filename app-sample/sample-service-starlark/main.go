package main

import (
	"context"
	_ "embed"
	"net/http"
	"os"
	"os/signal"

	_ "github.com/joho/godotenv/autoload"

	"github.com/duxthemux/embeddedlangs/app-sample/sample-service-starlark/motor_calculo"
)

func run(ctx context.Context) error {
	root := os.Getenv("ROOT")

	svc := &motor_calculo.Service{Path: root}
	if err := svc.Init(); err != nil {
		return err
	}

	api := &motor_calculo.Api{Service: svc}

	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: api,
	}

	go func() {
		<-ctx.Done()
		_ = httpServer.Shutdown(context.TODO())
	}()

	return httpServer.ListenAndServe()
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	if err := run(ctx); err != nil {
		panic(err)
	}
}
