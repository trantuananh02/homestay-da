package types

type Booking struct {
	ID            int           `json:"id"`
	BookingCode   string        `json:"bookingCode"`
	HomestayID    int           `json:"homestayId"`
	CustomerName  string        `json:"customerName"`
	CustomerPhone string        `json:"customerPhone"`
	CustomerEmail string        `json:"customerEmail"`
	CheckIn       string        `json:"checkIn"`
	CheckOut      string        `json:"checkOut"`
	Nights        int           `json:"nights"`
	TotalAmount   float64       `json:"totalAmount"`
	PaidAmount    float64       `json:"paidAmount"`
	Status        string        `json:"status"`
	BookingDate   string        `json:"bookingDate"`
	PaymentMethod string        `json:"paymentMethod"`
	Rooms         []BookingRoom `json:"rooms"`
	Review        Review        `json:"review"`
}

type BookingRoom struct {
	RoomID   int     `json:"id"`
	RoomName string  `json:"name"`
	RoomType string  `json:"type"`
	Capacity int     `json:"capacity"`
	Price    float64 `json:"pricePerNight"`
	Nights   int     `json:"nights"`
	SubTotal float64 `json:"subtotal"`
}

type CreateBookingReq struct {
	HomestayID    int              `json:"homestayId"`
	CustomerName  string           `json:"customerName"`
	CustomerPhone string           `json:"customerPhone"`
	CustomerEmail string           `json:"customerEmail"`
	CheckIn       string           `json:"checkIn"`
	CheckOut      string           `json:"checkOut"`
	Guests        int              `json:"guests"`
	PaymentMethod string           `json:"paymentMethod"`
	Notes         string           `json:"notes,omitempty"`
	TotalAmount   float64          `json:"totalAmount"`
	PaidAmount    float64          `json:"paidAmount,omitempty"`
	Rooms         []BookingRoomReq `json:"rooms"`
}

type BookingRoomReq struct {
	RoomID   int     `json:"id"`
	RoomName string  `json:"name"`
	RoomType string  `json:"type"`
	Capacity int     `json:"capacity"`
	Price    float64 `json:"pricePerNight"`
	Nights   int     `json:"nights"`
	SubTotal float64 `json:"subtotal"`
}

type CreateBookingResp struct {
	Booking Booking `json:"booking"`
}

type FilterBookingReq struct {
	Status        *string `form:"status"`
	CustomerName  *string `form:"customerName"`
	CustomerEmail *string `form:"customerEmail"`
	CustomerPhone *string `form:"customerPhone"`
	DateFrom      *string `form:"dateFrom"`
	DateTo        *string `form:"dateTo"`
	Page          int     `form:"page"`
	PageSize      int     `form:"pageSize"`
}

type FilterBookingResp struct {
	Bookings []Booking `json:"bookings"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}

type Payment struct {
	ID            int     `json:"id"`
	BookingCode   string  `json:"bookingCode"`
	BookingID     int  `json:"bookingId"`
	Amount        float64 `json:"amount"`
	PaymentMethod string  `json:"paymentMethod"`
	PaymentStatus string  `json:"paymentStatus"`
	TransactionID string  `json:"transactionId"`
	PaymentDate   string  `json:"paymentDate"`
}

type BookingDetailResp struct {
	Booking  Booking   `json:"booking"`
	Payments []Payment `json:"payments"`
}

type UpdateBookingStatusReq struct {
	Status string `json:"status"`
}

type UpdateBookingStatusResp struct {
	Success bool `json:"success"`
}

type UploadFileReq struct {
}

type UploadFileRes struct {
	Url string `json:"url"`
}

type GetBookingsByHomestayIDResp struct {
	Bookings []Booking `json:"bookings"`
	Total    int       `json:"total"`
}

type BookingEmailData struct {
	GuestName    string
	HomestayName string
	CheckInDate  string
	CheckOutDate string
	Nights       int
	Rooms        int
	TotalPrice   string
	Year         int
	BookingLink  string
}

type CreateReviewReq struct {
	BookingID int    `json:"bookingId"`
	Comment   string `json:"comment"`
	Rating    int    `json:"rating"`
}

type FilterPaymentReq struct {
	BookingCode *string `form:"bookingCode"`
	Method      *string `form:"method"`
	DateFrom    *string `form:"dateFrom"`
	DateTo      *string `form:"dateTo"`
	Page        int     `form:"page"`
	PageSize    int     `form:"pageSize"`
}

type FilterPaymentResp struct {
	Payments []Payment `json:"payments"`
	Total    int       `json:"total"`
	Page     int       `json:"page"`
	PageSize int       `json:"pageSize"`
}