package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"urlshortener/pkg/urlshortener"
	grpcApi "urlshortener/pkg/urlshortener/api/grpc"
	httpApi "urlshortener/pkg/urlshortener/api/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {
	cfg := readConfig()
	app, err := urlshortener.MakeUrlShortener(cfg)
	if err != nil {
		log.Fatalln(err.Error())
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT)

	httpHandlers := httpApi.MakeUrlShortenerHttpHandlers(&app)
	httpEngine := gin.Default()
	httpEngine.GET("/urls/:url", httpHandlers.GetHandler)
	httpEngine.POST("/urls", httpHandlers.PostHandler)
	httpServer := http.Server{Addr: ":8000", Handler: httpEngine}

	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			log.Fatalln(err.Error())
		}
	}()

	grpcHandlers := grpcApi.MakeUrlShortenerGrpcHandlers(&app)
	grpcServer := grpc.NewServer()
	grpcApi.RegisterUrlShortenerServer(grpcServer, &grpcHandlers)
	grpcListener, err := net.Listen("tcp", ":9111")
	if err != nil {
		log.Fatalln(err.Error())
	}

	go func() {
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalln(err.Error())
		}
	}()

	<-shutdown
	ctx, _ := context.WithTimeout(context.Background(), time.Minute)
	grpcServer.GracefulStop()
	if err = httpServer.Shutdown(ctx); err != nil {
		log.Fatalln(err.Error())
	}
}
