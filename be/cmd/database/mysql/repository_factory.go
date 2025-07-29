package mysql

import (
	"homestay-be/cmd/database/repo"

	"github.com/jmoiron/sqlx"
)

// RepositoryFactory chứa tất cả các repository
type RepositoryFactory struct {
	UserRepo             repo.UserRepository
	HomestayRepo         repo.HomestayRepository
	RoomRepo             repo.RoomRepository
	RoomAvailabilityRepo repo.RoomAvailabilityRepository
	BookingRequestRepo   repo.BookingRequestRepository
	BookingRepo          repo.BookingRepository
	PaymentRepo          repo.PaymentRepository
	ReviewRepo           repo.ReviewRepository
}

// NewRepositoryFactory tạo instance mới của RepositoryFactory
func NewRepositoryFactory(db *sqlx.DB) *RepositoryFactory {
	return &RepositoryFactory{
		UserRepo:             NewUserRepository(db),
		HomestayRepo:         NewHomestayRepo(db),
		RoomRepo:             NewRoomRepository(db),
		RoomAvailabilityRepo: NewRoomAvailabilityRepository(db),
		BookingRequestRepo:   NewBookingRequestRepository(db),
		BookingRepo:          NewBookingRepository(db),
		PaymentRepo:          NewPaymentRepository(db),
		ReviewRepo:           NewReviewRepository(db),
	}
}
