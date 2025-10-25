package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/mzulfanw/gateway-service/config"
	"github.com/mzulfanw/gateway-service/internal/handler"
	"github.com/mzulfanw/gateway-service/internal/middleware"
	"github.com/mzulfanw/gateway-service/internal/proxy"
	"github.com/sirupsen/logrus"
)

func main() {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})

	loadConfig, err := config.LoadConfig(".env")
	if err != nil {
		log.Errorf("failed to load config: %v", err)
	}

	p := proxy.NewProxy(loadConfig.ProductServiceUrl, loadConfig.OrderServiceUrl)
	r := mux.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(log))

	productHandler := handler.NewProductHandler(p)
	orderHandler := handler.NewOrderHandler(p)

	// registering route of product service
	r.HandleFunc("/products", productHandler.CreateProduct).Methods(http.MethodPost)
	r.HandleFunc("/products/{id}", productHandler.GetProductDetail).Methods(http.MethodGet)

	// registering route of order service
	r.HandleFunc("/orders", orderHandler.CreateOrder).Methods(http.MethodPost)
	r.HandleFunc("/orders/product/{productId}", orderHandler.GetByProductID).Methods(http.MethodGet)

	srv := &http.Server{
		Addr:         ":3000",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Infof("BFF Service running on port %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
	} else {
		log.Println("Server exiting")
	}
}
