package types

import "time"

// Room types
type Room struct {
	ID          int       `json:"id" db:"id"`
	HomestayID  int       `json:"homestayId" db:"homestay_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Type        string    `json:"type" db:"type"` // single, double, family, dormitory
	Capacity    int       `json:"capacity" db:"capacity"`
	Price       float64   `json:"price" db:"price"`
	PriceType   string    `json:"priceType" db:"price_type"` // per_night, per_person
	Status      string    `json:"status" db:"status"`        // available, occupied, maintenance
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
	Images      []string  `json:"images" db:"image_urls"`
	Amenities   []string  `json:"amenities" db:"amenities"`
}

// CreateRoomRequest - Request to create a new room
type CreateRoomRequest struct {
	HomestayID  int      `json:"homestayId" validate:"required"`
	Name        string   `json:"name" validate:"required,min=2,max=100"`
	Description string   `json:"description" validate:"required,min=10,max=500"`
	Type        string   `json:"type" validate:"required,oneof=single double family dormitory"`
	Capacity    int      `json:"capacity" validate:"required,min=1,max=20"`
	Price       float64  `json:"price" validate:"required,min=0"`
	PriceType   string   `json:"priceType" validate:"required,oneof=per_night per_person"`
	Images      []string `json:"images"`
	Amenities   []string `json:"amenities"`
}

// UpdateRoomRequest - Request to update a room
type UpdateRoomRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string  `json:"description,omitempty" validate:"omitempty,min=10,max=500"`
	Type        *string  `json:"type,omitempty" validate:"omitempty,oneof=single double family dormitory"`
	Capacity    *int     `json:"capacity,omitempty" validate:"omitempty,min=1,max=20"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
	PriceType   *string  `json:"priceType,omitempty" validate:"omitempty,oneof=per_night per_person"`
	Status      *string  `json:"status,omitempty" validate:"omitempty,oneof=available occupied maintenance"`
	Images      []string `json:"images,omitempty"`
	Amenities   []string `json:"amenities,omitempty"`
}

// RoomListRequest - Request to get list of rooms
type RoomListRequest struct {
	HomestayID int      `json:"homestayId" form:"homestayId" validate:"required"`
	Page       int      `json:"page" form:"page" validate:"min=1"`
	PageSize   int      `json:"pageSize" form:"pageSize" validate:"min=1,max=100"`
	Status     string   `json:"status" form:"status" validate:"omitempty,oneof=available occupied maintenance"`
	Type       string   `json:"type" form:"type" validate:"omitempty,oneof=single double family dormitory"`
	MinPrice   *float64 `json:"minPrice" form:"minPrice" validate:"omitempty,min=0"`
	MaxPrice   *float64 `json:"maxPrice" form:"maxPrice" validate:"omitempty,min=0"`
}

// RoomListResponse - Response for room list
type RoomListResponse struct {
	Rooms     []Room `json:"rooms"`
	Total     int    `json:"total"`
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
	TotalPage int    `json:"totalPage"`
}

// RoomDetailResponse - Response for room detail
type RoomDetailResponse struct {
	Room           Room               `json:"room"`
	Homestay       Homestay           `json:"homestay"`
	Availabilities []RoomAvailability `json:"availabilities,omitempty"`
}

// RoomAvailability types
type RoomAvailability struct {
	ID        int       `json:"id" db:"id"`
	RoomID    int       `json:"roomId" db:"room_id"`
	Date      time.Time `json:"date" db:"date"`
	Status    string    `json:"status" db:"status"` // available, booked, blocked
	Price     *float64  `json:"price" db:"price"`   // Override price for specific date
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// CreateAvailabilityRequest - Request to create room availability
type CreateAvailabilityRequest struct {
	RoomID int       `json:"roomId" validate:"required"`
	Date   time.Time `json:"date" validate:"required"`
	Status string    `json:"status" validate:"required,oneof=available booked blocked"`
	Price  *float64  `json:"price,omitempty" validate:"omitempty,min=0"`
}

// UpdateAvailabilityRequest - Request to update room availability
type UpdateAvailabilityRequest struct {
	Status *string  `json:"status,omitempty" validate:"omitempty,oneof=available booked blocked"`
	Price  *float64 `json:"price,omitempty" validate:"omitempty,min=0"`
}

// BulkAvailabilityRequest - Request to update multiple availabilities
type BulkAvailabilityRequest struct {
	RoomID       int         `json:"roomId" validate:"required"`
	StartDate    time.Time   `json:"startDate" validate:"required"`
	EndDate      time.Time   `json:"endDate" validate:"required"`
	Status       string      `json:"status" validate:"required,oneof=available booked blocked"`
	Price        *float64    `json:"price,omitempty" validate:"omitempty,min=0"`
	ExcludeDates []time.Time `json:"excludeDates,omitempty"`
}

// RoomStatsResponse - Response for room statistics
type RoomStatsResponse struct {
	TotalRooms       int     `json:"totalRooms"`
	AvailableRooms   int     `json:"availableRooms"`
	OccupiedRooms    int     `json:"occupiedRooms"`
	MaintenanceRooms int     `json:"maintenanceRooms"`
	AveragePrice     float64 `json:"averagePrice"`
	TotalRevenue     float64 `json:"totalRevenue"`
	OccupancyRate    float64 `json:"occupancyRate"`
}
