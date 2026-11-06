package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/PreetiToppo/movie-booking/internal/adapters/redis"
	"github.com/PreetiToppo/movie-booking/internal/booking"
	"github.com/PreetiToppo/movie-booking/internal/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /movies", listMovies)

	mux.Handle("GET /", http.FileServer(http.Dir("static")))

	store := booking.NewRedisStore(redis.NewClient("localhost:6379"))
	svc := booking.NewService(store)

	bookingHandler := booking.NewHandler(svc)

	mux.HandleFunc("GET /movies/{movieID}/seats", bookingHandler.ListSeats)
	mux.HandleFunc("POST /movies/{movieID}/seats/{seatID}/hold", bookingHandler.HoldSeat)

	mux.HandleFunc("PUT /sessions/{sessionID}/confirm", bookingHandler.ConfirmSession)
	mux.HandleFunc("DELETE /sessions/{sessionID}", bookingHandler.ReleaseSession)

	mux.Handle("GET /metrics", promhttp.Handler())

	// 1. Wrap your mux in a custom http.Server
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// 2. Start the server in a background goroutine
	go func() {
		log.Println("Starting server on http://localhost:8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 3. Wait for an interrupt signal (Ctrl+C)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// 4. Shutdown with a 5-second timeout to let active requests finish
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting gracefully")
}

var movies = []movieResponse{
	{ID: "inception", Title: "Inception", Rows: 5, SeatsPerRow: 8},
	{ID: "dune", Title: "Dune: Part Two", Rows: 4, SeatsPerRow: 6},
}

func listMovies(w http.ResponseWriter, r *http.Request) {
	utils.WriteJSON(w, http.StatusOK, movies)
}

type movieResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}