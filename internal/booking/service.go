package booking

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Service struct {
	store  BookingStore
	pgPool *pgxpool.Pool // Added Postgres Pool
}

func NewService(store BookingStore, pgPool *pgxpool.Pool) *Service {
	return &Service{
		store:  store,
		pgPool: pgPool,
	}
}

func (s *Service) Book(b Booking) (Booking, error) {
	return s.store.Book(b)
}

func (s *Service) ListBookings(movieID string) []Booking {
	return s.store.ListBookings(movieID)
}

func (s *Service) ConfirmSeat(ctx context.Context, sessionID string, userID string) (Booking, error) {
	// 1. Confirm in Redis (removes TTL)
	session, err := s.store.Confirm(ctx, sessionID, userID)
	if err != nil {
		return Booking{}, err
	}

	// 2. Insert into Postgres (Permanent Storage)
	query := `INSERT INTO confirmed_bookings (id, movie_id, seat_id, user_id) VALUES ($1, $2, $3, $4)`
	_, err = s.pgPool.Exec(ctx, query, session.ID, session.MovieID, session.SeatID, session.UserID)
	if err != nil {
		log.Printf("Failed to save to Postgres: %v", err)
		// Even if Postgres fails, the booking is still in Redis. 
		// In a real system, you'd use the Saga pattern here to roll back.
	} else {
		log.Printf("Booking %s permanently saved to Postgres", session.ID)
	}

	return session, nil
}

func (s *Service) ReleaseSeat(ctx context.Context, sessionID string, userID string) error {
	return s.store.Release(ctx, sessionID, userID)
}