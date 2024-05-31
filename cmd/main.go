package main

import (
	"log"
	"net/http"
	"shopTestTask/cfg"
	"shopTestTask/handlers"
	"shopTestTask/middleware"
)

func main() {
	mux := http.NewServeMux()
	h := handlers.New()

	mux.Handle("POST /order", handlers.CtxHandler(h.PlaceOrder))
	mux.Handle("GET /order/{id}", handlers.CtxHandler(h.QueryOrder))
	mux.Handle("GET /order/{id}/products", handlers.CtxHandler(h.QueryOrderProducts))
	mux.Handle("DELETE /order/{id}", handlers.CtxHandler(h.CancelOrder))
	mux.Handle("PATCH /order/{id}", handlers.CtxHandler(h.ChangeOrder))

	supaware := middleware.Combine(
		middleware.Auth,
		middleware.Logger,
	)(mux)
	if err := http.ListenAndServe(cfg.Address, supaware); err != nil {
		log.Fatal(err)
	}
}
