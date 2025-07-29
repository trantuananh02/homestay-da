package types

import "time"

// Types Review
type Review struct {
	ID         int       `json:"id" db:"id"`
	Rating     int       `json:"rating" db:"rating"`
	Comment    string    `json:"comment" db:"comment"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
	BookingID  int       `json:"bookingId" db:"booking_id"`
	GuestID    int       `json:"guestId" db:"guest_id"`
	GuestName  string    `json:"guestName" db:"guest_name"`
	HomestayID int       `json:"homestayId" db:"homestay_id"`
}

// Homestay types
type Homestay struct {
	ID           int       `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Description  string    `json:"description" db:"description"`
	Address      string    `json:"address" db:"address"`
	City         string    `json:"city" db:"city"`
	District     string    `json:"district" db:"district"`
	Ward         string    `json:"ward" db:"ward"`
	Latitude     float64   `json:"latitude" db:"latitude"`
	Longitude    float64   `json:"longitude" db:"longitude"`
	HostID       int       `json:"hostId" db:"host_id"`
	Status       string    `json:"status" db:"status"`              // active, inactive, pending
	Rate         float64   `json:"rate" db:"rate"`                  // Average rating, optional
	Rooms        []Room    `json:"rooms"`                           // Optional, populated in detail response
	Reviews      []Review  `json:"reviews"`                         // Optional, populated in detail response
	Rating       float64   `json:"rating" db:"rating"`              // Average rating, optional
	TotalReviews int       `json:"totalReviews" db:"total_reviews"` // Total number of reviews, optional
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

// CreateHomestayRequest - Request to create a new homestay
type CreateHomestayRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Description string  `json:"description" validate:"required,min=10,max=1000"`
	Address     string  `json:"address" validate:"required,min=5,max=200"`
	City        string  `json:"city" validate:"required,min=2,max=50"`
	District    string  `json:"district" validate:"required,min=2,max=50"`
	Ward        string  `json:"ward" validate:"required,min=2,max=50"`
	Latitude    float64 `json:"latitude" validate:"required"`
	Longitude   float64 `json:"longitude" validate:"required"`
}

// UpdateHomestayRequest - Request to update a homestay
type UpdateHomestayRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string  `json:"description,omitempty" validate:"omitempty,min=10,max=1000"`
	Address     *string  `json:"address,omitempty" validate:"omitempty,min=5,max=200"`
	City        *string  `json:"city,omitempty" validate:"omitempty,min=2,max=50"`
	District    *string  `json:"district,omitempty" validate:"omitempty,min=2,max=50"`
	Ward        *string  `json:"ward,omitempty" validate:"omitempty,min=2,max=50"`
	Latitude    *float64 `json:"latitude,omitempty" validate:"omitempty"`
	Longitude   *float64 `json:"longitude,omitempty" validate:"omitempty"`
	Status      *string  `json:"status,omitempty" validate:"omitempty,oneof=active inactive"`
}

// HomestayListRequest - Request to get list of homestays
type HomestayListRequest struct {
	Page     int    `json:"page" form:"page" validate:"min=1"`
	PageSize int    `json:"pageSize" form:"pageSize" validate:"min=1,max=100"`
	Search   string `json:"search" form:"search" validate:"omitempty"`
	Status   string `json:"status" form:"status" validate:"omitempty,oneof=active inactive"`
	City     string `json:"city" form:"city" validate:"omitempty"`
	District string `json:"district" form:"district" validate:"omitempty"`
	CheckIn  string `json:"checkIn" form:"checkIn" validate:"omitempty,datetime=2006-01-02"`
	CheckOut string `json:"checkOut" form:"checkOut" validate:"omitempty,datetime=2006-01-02"`
	Guests   int    `json:"guests" form:"guests" validate:"omitempty,min=1"`
}

// HomestayListResponse - Response for homestay list
type HomestayListResponse struct {
	Homestays []Homestay `json:"homestays"`
	Total     int        `json:"total"`
	Page      int        `json:"page"`
	PageSize  int        `json:"pageSize"`
	TotalPage int        `json:"totalPage"`
}

// HomestayDetailResponse - Response for homestay detail
type HomestayDetailResponse struct {
	Homestay Homestay `json:"homestay"`
	Rooms    []Room   `json:"rooms,omitempty"`
}

// HomestayStatsResponse - Response for homestay statistics
type HomestayStatsResponse struct {
	TotalHomestays  int     `json:"totalHomestays"`
	ActiveHomestays int     `json:"activeHomestays"`
	TotalRooms      int     `json:"totalRooms"`
	AvailableRooms  int     `json:"availableRooms"`
	TotalBookings   int     `json:"totalBookings"`
	TotalRevenue    float64 `json:"totalRevenue"`
	MonthlyRevenue  float64 `json:"monthlyRevenue"`
	OccupancyRate   float64 `json:"occupancyRate"`
}
