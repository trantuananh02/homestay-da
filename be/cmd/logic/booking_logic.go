package logic

import (
	"context"
	"errors"
	"fmt"
	"homestay-be/cmd/database/model"
	"homestay-be/cmd/svc"
	"homestay-be/cmd/types"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type BookingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewBookingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BookingLogic {
	return &BookingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// FilterBookings - Logic to filter bookings based on criteria
func (l *BookingLogic) FilterBookings(ctx context.Context, req *types.FilterBookingReq) (*types.FilterBookingResp, error) {
	logx.Info(req)

	// 1. Mapping filter sang model.BookingSearchRequest
	searchReq := &model.BookingSearchRequest{
		Status:        req.Status,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		CustomerPhone: req.CustomerPhone,
		Page:          req.Page,
		PageSize:      req.PageSize,
	}
	// Nếu có filter ngày
	if req.DateFrom != nil && *req.DateTo != "" {
		t, _ := time.Parse("2006-01-02", *req.DateFrom)
		searchReq.StartDate = &t
	}
	if req.DateTo != nil && *req.DateTo != "" {
		t, _ := time.Parse("2006-01-02", *req.DateTo)
		searchReq.EndDate = &t
	}

	// 2. Lấy danh sách booking
	bookings, total, err := l.svcCtx.BookingRepo.Search(ctx, searchReq)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	logx.Info(bookings)

	var respBookings []types.Booking
	for _, booking := range bookings {
		var nights int
		if booking.CheckIn.IsZero() || booking.CheckOut.IsZero() {
			nights = 0
		} else {
			nights = int(booking.CheckOut.Sub(booking.CheckIn).Hours() / 24) // Tính số đêm
		}

		// 3. Lấy danh sách phòng cho mỗi booking
		rooms, err := l.svcCtx.BookingRepo.GetRoomsByBookingID(ctx, booking.ID)
		if err != nil {
			logx.Error(err)
			return nil, err
		}
		var respRooms []types.BookingRoom
		for _, r := range rooms {
			respRooms = append(respRooms, types.BookingRoom{
				RoomID:   r.RoomID,
				RoomName: r.RoomName,
				RoomType: r.RoomType,
				Capacity: r.Capacity,
				Price:    r.Price,
				Nights:   nights,                    // Tính số đêm
				SubTotal: r.Price * float64(nights), // Tính số đêm
			})
		}

		// get review
		var review types.Review
		reviewModel, err := l.svcCtx.BookingRepo.GetReviewByBookingID(ctx, booking.ID)
		if err != nil {
			logx.Error(err)
		}
		if reviewModel != nil {
			review = types.Review{
				ID:         reviewModel.ID,
				BookingID:  reviewModel.BookingID,
				Rating:     reviewModel.Rating,
				Comment:    reviewModel.Comment,
				CreatedAt:  reviewModel.CreatedAt,
				GuestID:    reviewModel.UserID,
				GuestName:  reviewModel.UserName,
				HomestayID: reviewModel.HomestayID,
			}
		}

		respBookings = append(respBookings, types.Booking{
			ID:            booking.ID,
			BookingCode:   booking.BookingCode,
			CustomerName:  booking.Name,
			CustomerPhone: booking.Phone,
			CustomerEmail: booking.Email,
			CheckIn:       booking.CheckIn.Format("2006-01-02"),
			CheckOut:      booking.CheckOut.Format("2006-01-02"),
			TotalAmount:   booking.TotalAmount,
			Status:        booking.Status,
			Nights:        nights,
			BookingDate:   booking.CreatedAt.Format("2006-01-02 15:04:05"),
			PaymentMethod: booking.PaymentMethod,
			PaidAmount:    booking.PaidAmount,
			Rooms:         respRooms,
			Review:        review,
		})
	}

	return &types.FilterBookingResp{
		Bookings: respBookings,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}

// CreateBooking - Logic to create a new booking
func (l *BookingLogic) CreateBooking(ctx context.Context, req *types.CreateBookingReq) (*types.CreateBookingResp, error) {
	logx.Info(req)

	// validate request
	if req.CustomerName == "" || req.CustomerEmail == "" || req.CheckIn == "" || req.CheckOut == "" || len(req.Rooms) == 0 {
		return nil, errors.New("các trường bắt buộc không được để trống")
	}

	// validate check-in/check-out date
	checkIn, err := time.Parse("2006-01-02", req.CheckIn)
	if err != nil {
		return nil, errors.New("ngày check-in không hợp lệ")
	}
	checkOut, err := time.Parse("2006-01-02", req.CheckOut)
	if err != nil {
		return nil, errors.New("ngày check-out không hợp lệ")
	}
	if checkIn.After(checkOut) {
		return nil, errors.New("ngày check-in phải trước ngày check-out")
	}

	// // validate guests
	// if req.Guests <= 0 {
	// 	return nil, errors.New("số lượng khách phải lớn hơn 0")
	// }

	// kiểm tra phòng đã được đặt hay chưa
	for _, room := range req.Rooms {
		if room.RoomID == 0 {
			return nil, errors.New("phòng không hợp lệ")
		}
		// Kiểm tra xem phòng đã được đặt hay chưa
		exists, err := l.svcCtx.BookingRepo.CheckRoomExists(ctx, room.RoomID, checkIn, checkOut)
		if err != nil {
			logx.Error(err)
			return nil, err
		}
		if exists {
			return nil, errors.New("phòng đã được đặt trong khoảng thời gian này")
		}
	}
	// thêm booking
	var bookingCode string = "BK" + time.Now().Format("20060102150405")
	var status string = "confirmed" // Mặc định là confirmed
	if req.PaidAmount >= req.TotalAmount {
		status = "completed" // Nếu đã thanh toán đủ thì chuyển sang completed
	}
	bookingModel := &model.BookingCreateRequest{
		BookingCode:   bookingCode,
		HomestayID:    req.HomestayID,
		Name:          req.CustomerName,
		Email:         req.CustomerEmail,
		Phone:         req.CustomerPhone,
		CheckIn:       parseDate(req.CheckIn),
		CheckOut:      parseDate(req.CheckOut),
		Status:        status,
		TotalAmount:   req.TotalAmount,
		PaidAmount:    req.PaidAmount,
		PaymentMethod: req.PaymentMethod,
	}
	booking, err := l.svcCtx.BookingRepo.Create(ctx, bookingModel)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	if booking.PaidAmount > 0 {
		// Xử lý thanh toán
		payment := &model.PaymentCreateRequest{
			BookingID:     booking.ID,
			Amount:        booking.PaidAmount,
			PaymentMethod: booking.PaymentMethod,
			PaymentStatus: "completed",
			TransactionID: "",
			PaymentDate:   time.Now(),
		}
		_, err := l.svcCtx.PaymentRepo.Create(ctx, payment)
		if err != nil {
			logx.Error(err)
			return nil, err
		}
	}

	// 2. Tạo các bản ghi booking_room
	for _, room := range req.Rooms {
		roomModel := &model.BookingRoom{
			BookingID: booking.ID,
			RoomID:    room.RoomID,
			RoomName:  room.RoomName,
			RoomType:  room.RoomType,
			Capacity:  room.Capacity,
			Price:     room.Price,
			CreatedAt: time.Now(),
		}

		// Giả sử có hàm InsertBookingRoom trong repo
		_, err := l.svcCtx.BookingRepo.InsertBookingRoom(ctx, roomModel)
		if err != nil {
			logx.Error(err)
			return nil, err
		}
	}

	// 3. Lấy lại danh sách phòng vừa insert
	rooms, err := l.svcCtx.BookingRepo.GetRoomsByBookingID(ctx, booking.ID)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	// 4. Mapping sang types.BookingRoom
	var respRooms []types.BookingRoom
	for _, r := range rooms {
		respRooms = append(respRooms, types.BookingRoom{
			RoomID:   r.RoomID,
			RoomName: r.RoomName,
			RoomType: r.RoomType,
			Capacity: r.Capacity,
			Price:    r.Price,
			Nights:   int(booking.CheckOut.Sub(booking.CheckIn).Hours() / 24),                  // Tính số đêm
			SubTotal: r.Price * float64(int(booking.CheckOut.Sub(booking.CheckIn).Hours()/24)), // Tính số đêm
		})
	}

	// 5. Mapping sang types.Booking
	respBooking := types.Booking{
		ID:            booking.ID,
		BookingCode:   booking.BookingCode,
		CustomerName:  booking.Name,
		CustomerPhone: booking.Phone,
		CustomerEmail: booking.Email,
		CheckIn:       booking.CheckIn.Format("2006-01-02"),
		CheckOut:      booking.CheckOut.Format("2006-01-02"),
		TotalAmount:   booking.TotalAmount,
		PaidAmount:    booking.PaidAmount,
		BookingDate:   booking.CreatedAt.Format("2006-01-02 15:04:05"),
		PaymentMethod: booking.PaymentMethod,
		Status:        booking.Status,
		Nights:        int(booking.CheckOut.Sub(booking.CheckIn).Hours() / 24), // Tính số đêm
		Rooms:         respRooms,
	}

	return &types.CreateBookingResp{Booking: respBooking}, nil
}

// create guest booking
func (l *BookingLogic) CreateGuestBooking(ctx context.Context, req *types.CreateBookingReq) (*types.CreateBookingResp, error) {
	logx.Info(req)

	// validate request
	if req.CustomerName == "" || req.CustomerEmail == "" || req.CheckIn == "" || req.CheckOut == "" || len(req.Rooms) == 0 {
		return nil, errors.New("các trường bắt buộc không được để trống")
	}

	// validate check-in/check-out date
	checkIn, err := time.Parse("2006-01-02", req.CheckIn)
	if err != nil {
		return nil, errors.New("ngày check-in không hợp lệ")
	}
	checkOut, err := time.Parse("2006-01-02", req.CheckOut)
	if err != nil {
		return nil, errors.New("ngày check-out không hợp lệ")
	}
	if checkIn.After(checkOut) {
		return nil, errors.New("ngày check-in phải trước ngày check-out")
	}

	// // validate guests
	// if req.Guests <= 0 {
	// 	return nil, errors.New("số lượng khách phải lớn hơn 0")
	// }

	// kiểm tra phòng đã được đặt hay chưa
	for _, room := range req.Rooms {
		if room.RoomID == 0 {
			return nil, errors.New("phòng không hợp lệ")
		}
		// Kiểm tra xem phòng đã được đặt hay chưa
		exists, err := l.svcCtx.BookingRepo.CheckRoomExists(ctx, room.RoomID, checkIn, checkOut)
		if err != nil {
			logx.Error(err)
			return nil, err
		}
		if exists {
			return nil, errors.New("phòng đã được đặt trong khoảng thời gian này")
		}
	}
	// thêm booking
	var bookingCode string = "BK" + time.Now().Format("20060102150405")

	bookingModel := &model.BookingCreateRequest{
		BookingCode:   bookingCode,
		HomestayID:    req.HomestayID,
		Name:          req.CustomerName,
		Email:         req.CustomerEmail,
		Phone:         req.CustomerPhone,
		CheckIn:       parseDate(req.CheckIn),
		CheckOut:      parseDate(req.CheckOut),
		TotalAmount:   req.TotalAmount,
		Status:        "pending",
		PaidAmount:    req.PaidAmount,
		PaymentMethod: req.PaymentMethod,
	}
	booking, err := l.svcCtx.BookingRepo.Create(ctx, bookingModel)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	if booking.PaidAmount > 0 {
		// Xử lý thanh toán
		payment := &model.PaymentCreateRequest{
			BookingID:     booking.ID,
			Amount:        booking.PaidAmount,
			PaymentMethod: booking.PaymentMethod,
			PaymentStatus: "completed",
			TransactionID: "",
			PaymentDate:   time.Now(),
		}
		_, err := l.svcCtx.PaymentRepo.Create(ctx, payment)
		if err != nil {
			logx.Error(err)
			return nil, err
		}
	}

	// 2. Tạo các bản ghi booking_room
	for _, room := range req.Rooms {
		roomModel := &model.BookingRoom{
			BookingID: booking.ID,
			RoomID:    room.RoomID,
			RoomName:  room.RoomName,
			RoomType:  room.RoomType,
			Capacity:  room.Capacity,
			Price:     room.Price,
			CreatedAt: time.Now(),
		}

		// Giả sử có hàm InsertBookingRoom trong repo
		_, err := l.svcCtx.BookingRepo.InsertBookingRoom(ctx, roomModel)
		if err != nil {
			logx.Error(err)
			return nil, err
		}
	}

	// 3. Lấy lại danh sách phòng vừa insert
	rooms, err := l.svcCtx.BookingRepo.GetRoomsByBookingID(ctx, booking.ID)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	// 4. Mapping sang types.BookingRoom
	var respRooms []types.BookingRoom
	for _, r := range rooms {
		respRooms = append(respRooms, types.BookingRoom{
			RoomID:   r.RoomID,
			RoomName: r.RoomName,
			RoomType: r.RoomType,
			Capacity: r.Capacity,
			Price:    r.Price,
			Nights:   int(booking.CheckOut.Sub(booking.CheckIn).Hours() / 24),                  // Tính số đêm
			SubTotal: r.Price * float64(int(booking.CheckOut.Sub(booking.CheckIn).Hours()/24)), // Tính số đêm
		})
	}

	// 5. Mapping sang types.Booking
	respBooking := types.Booking{
		ID:            booking.ID,
		BookingCode:   booking.BookingCode,
		CustomerName:  booking.Name,
		CustomerPhone: booking.Phone,
		CustomerEmail: booking.Email,
		CheckIn:       booking.CheckIn.Format("2006-01-02"),
		CheckOut:      booking.CheckOut.Format("2006-01-02"),
		TotalAmount:   booking.TotalAmount,
		PaidAmount:    booking.PaidAmount,
		BookingDate:   booking.CreatedAt.Format("2006-01-02 15:04:05"),
		PaymentMethod: booking.PaymentMethod,
		Status:        booking.Status,
		Nights:        int(booking.CheckOut.Sub(booking.CheckIn).Hours() / 24), // Tính số đêm
		Rooms:         respRooms,
	}

	return &types.CreateBookingResp{Booking: respBooking}, nil
}

func parseDate(dateStr string) time.Time {
	t, _ := time.Parse("2006-01-02", dateStr)
	return t
}

// Get Detail Booking - Logic to get details of a booking
func (l *BookingLogic) GetBookingDetail(ctx context.Context, bookingID int) (*types.BookingDetailResp, error) {
	// 1. Lấy booking
	booking, err := l.svcCtx.BookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	// 2. Lấy danh sách phòng
	rooms, err := l.svcCtx.BookingRepo.GetRoomsByBookingID(ctx, bookingID)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	var nights int = int(booking.CheckOut.Sub(booking.CheckIn).Hours() / 24) // Tính số đêm
	var respRooms []types.BookingRoom
	for _, r := range rooms {
		respRooms = append(respRooms, types.BookingRoom{
			RoomID:   r.RoomID,
			RoomName: r.RoomName,
			RoomType: r.RoomType,
			Capacity: r.Capacity,
			Price:    r.Price,
			Nights:   nights,                    // Tính số đêm
			SubTotal: r.Price * float64(nights), // Tính số đêm
		})
	}

	// // 3. Lấy danh sách payment
	// payments, _, err := l.svcCtx.PaymentRepo.GetByBookingID(ctx, bookingID, 1, 100)
	// if err != nil {
	// 	logx.Error(err)
	// 	return nil, err
	// }
	// var respPayments []types.Payment
	// for _, p := range payments {
	// 	respPayments = append(respPayments, types.Payment{
	// 		ID:            p.ID,
	// 		Amount:        p.Amount,
	// 		PaymentMethod: p.PaymentMethod,
	// 		PaymentStatus: p.PaymentStatus,
	// 		TransactionID: p.TransactionID,
	// 		PaymentDate:   p.PaymentDate.Format("2006-01-02 15:04:05"),
	// 	})
	// }
	// 4. Mapping sang types.Booking
	respBooking := types.Booking{
		ID:            booking.ID,
		CustomerName:  booking.Name,
		CustomerPhone: booking.Phone,
		CustomerEmail: booking.Email,
		CheckIn:       booking.CheckIn.Format("2006-01-02"),
		CheckOut:      booking.CheckOut.Format("2006-01-02"),
		TotalAmount:   booking.TotalAmount,
		PaidAmount:    booking.PaidAmount,
		Status:        booking.Status,
		BookingCode:   booking.BookingCode,
		BookingDate:   booking.CreatedAt.Format("2006-01-02 15:04:05"),
		PaymentMethod: booking.PaymentMethod,
		Nights:        nights,
		Rooms:         respRooms,
	}
	return &types.BookingDetailResp{
		Booking: respBooking,
		// Payments: respPayments,
	}, nil
}

// UpdateBookingStatus - Logic to update the status of a booking
func (l *BookingLogic) UpdateBookingStatus(ctx context.Context, bookingID int, req *types.UpdateBookingStatusReq) (*types.UpdateBookingStatusResp, error) {
	updateReq := &model.BookingUpdateRequest{
		Status: &req.Status,
	}

	// get booking by ID
	booking, err := l.svcCtx.BookingRepo.GetByID(ctx, bookingID)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	// Validate status
	if booking == nil {
		return nil, errors.New("booking not found")
	}

	if booking.Status == req.Status || (req.Status != "confirmed" && req.Status != "cancelled" && req.Status != "completed") {
		return nil, errors.New("invalid booking status")
	}

	// Nếu status là "cancelled" thì cần kiểm tra xem booking có đang trong trạng thái "confirmed" hay không
	if req.Status == "cancelled" && booking.Status != "pending" && booking.Status != "confirmed" {
		return nil, errors.New("only pending or confirmed bookings can be cancelled")
	}

	// Nếu status là "completed" thì cần kiểm tra xem booking có đang trong trạng thái "confirmed" hay không
	if req.Status == "completed" && booking.Status != "confirmed" {
		return nil, errors.New("only confirmed bookings can be completed")
	}

	// Nếu status là "confirmed" thì cần kiểm tra xem booking có đang trong trạng thái "pending" hay không
	if req.Status == "confirmed" && booking.Status != "pending" {
		return nil, errors.New("only pending bookings can be confirmed")
	}

	// Nếu status là "completed" thì thêm payment tương ứng
	if req.Status == "completed" {
		// Giả sử có hàm AddPayment trong repo để thêm payment
		payment := &model.PaymentCreateRequest{
			BookingID:     booking.ID,
			Amount:        booking.TotalAmount - booking.PaidAmount, // Chỉ thêm phần còn thiếu
			PaymentMethod: booking.PaymentMethod,
			PaymentStatus: "completed",
			TransactionID: "",
			PaymentDate:   time.Now(),
		}
		_, err := l.svcCtx.PaymentRepo.Create(ctx, payment)
		if err != nil {
			logx.Error(err)
			return nil, err
		}

		// Cập nhật lại số tiền đã thanh toán
		booking.PaidAmount += payment.Amount
		updateReq.PaidAmount = &booking.PaidAmount
		if booking.PaidAmount >= booking.TotalAmount {
			statusCompleted := "completed"
			updateReq.Status = &statusCompleted
		}
	}

	// update booking status
	_, err = l.svcCtx.BookingRepo.UpdateStatus(ctx, bookingID, updateReq)
	if err != nil {
		logx.Error(err)
		return &types.UpdateBookingStatusResp{Success: false}, err
	}

	// Gửi email thông báo nếu booking đã được xác nhận
	if req.Status == "confirmed" {
		var homestayName string

		// Lấy thông tin homestay từ booking
		bookingRoom, err := l.svcCtx.BookingRepo.GetRoomsByBookingID(ctx, booking.ID)
		if err != nil {
			logx.Error("Lấy thông tin phòng thất bại:", err)
			return &types.UpdateBookingStatusResp{Success: false}, err
		}
		if len(bookingRoom) != 0 {
			room, err := l.svcCtx.RoomRepo.GetByID(ctx, bookingRoom[0].RoomID)
			if err != nil {
				logx.Error("Lấy thông tin phòng thất bại:", err)
				return &types.UpdateBookingStatusResp{Success: false}, err
			}
			if room != nil {
				homestay, err := l.svcCtx.HomestayRepo.GetByID(ctx, room.HomestayID)
				if err != nil {
					logx.Error("Lấy thông tin homestay thất bại:", err)
					return &types.UpdateBookingStatusResp{Success: false}, err
				}

				if homestay != nil {
					homestayName = homestay.Name
				}
			}
		}

		err = l.svcCtx.MailClient.SendBookingConfirmation(booking.Email, types.BookingEmailData{
			GuestName:    booking.Name,
			HomestayName: homestayName,
			CheckInDate:  booking.CheckIn.Format("02-01-2006"),
			CheckOutDate: booking.CheckOut.Format("02-01-2006"),
			Nights:       int(booking.CheckOut.Sub(booking.CheckIn).Hours() / 24),
			Rooms:        len(bookingRoom),
			TotalPrice:   fmt.Sprintf("%.2f", booking.TotalAmount),
			Year:         booking.CreatedAt.Year(),
			BookingLink:  "http://localhost:5173/bookings",
		})
		if err != nil {
			logx.Error("Gửi email xác nhận thất bại:", err)
			return &types.UpdateBookingStatusResp{Success: false}, err
		}
	}

	return &types.UpdateBookingStatusResp{Success: true}, nil
}

// GetBookingsByHomestayID - Logic to get bookings by homestay ID
func (l *BookingLogic) GetBookingsByHomestayID(ctx context.Context, homestayID int) (*types.GetBookingsByHomestayIDResp, error) {
	// 1. Lấy danh sách booking theo homestay ID
	bookings, total, err := l.svcCtx.BookingRepo.GetByHomestayID(ctx, homestayID, 1, 100)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	var respBookings []types.Booking
	for _, booking := range bookings {
		var rooms []types.BookingRoom
		var nights int = int(booking.CheckOut.Sub(booking.CheckIn).Hours() / 24) // Tính số đêm
		// 2. Lấy danh sách phòng cho mỗi booking
		bookingRooms, err := l.svcCtx.BookingRepo.GetRoomsByBookingID(ctx, booking.ID)
		if err != nil {
			logx.Error(err)
			return nil, err
		}

		for _, r := range bookingRooms {
			// Giả sử có trường RoomID trong BookingRoom
			room, err := l.svcCtx.RoomRepo.GetByID(ctx, r.RoomID)
			if err != nil {
				logx.Error(err)
				return nil, err
			}

			rooms = append(rooms, types.BookingRoom{
				RoomID:   r.RoomID,
				RoomName: room.Name,
				RoomType: room.Type,
				Capacity: r.Capacity,
				Price:    r.Price,
				SubTotal: r.Price * float64(nights), // Tính số đêm
				Nights:   nights,
			})
		}

		respBookings = append(respBookings, types.Booking{
			ID:            booking.ID,
			BookingCode:   booking.BookingCode,
			CustomerName:  booking.Name,
			CustomerPhone: booking.Phone,
			CustomerEmail: booking.Email,
			CheckIn:       booking.CheckIn.Format("2006-01-02"),
			CheckOut:      booking.CheckOut.Format("2006-01-02"),
			TotalAmount:   booking.TotalAmount,
			PaidAmount:    booking.PaidAmount,
			Status:        booking.Status,
			Rooms:         rooms,
			Nights:        nights,
			BookingDate:   booking.CreatedAt.Format("2006-01-02 15:04:05"),
			PaymentMethod: booking.PaymentMethod,
		})
	}

	return &types.GetBookingsByHomestayIDResp{
		Bookings: respBookings,
		Total:    total,
	}, nil
}

// Create Review - Logic to create a review for a booking
func (l *BookingLogic) CreateReview(ctx context.Context, userID int, req *types.CreateReviewReq) error {
	logx.Info(req)

	// Validate request
	if req.BookingID <= 0 || req.Comment == "" || req.Rating < 1 || req.Rating > 5 {
		return errors.New("các trường bắt buộc không được để trống và rating phải từ 1 đến 5")
	}

	// Lấy thông tin booking
	booking, err := l.svcCtx.BookingRepo.GetByID(ctx, req.BookingID)
	if err != nil {
		logx.Error(err)
		return err
	}

	if booking == nil {
		return errors.New("booking không tồn tại")
	}

	// get booking room
	bookingRooms, err := l.svcCtx.BookingRepo.GetRoomsByBookingID(ctx, req.BookingID)
	if err != nil {
		logx.Error(err)
		return err
	}

	if len(bookingRooms) == 0 {
		return errors.New("booking không có phòng nào")
	}

	// Lấy homestay ID từ phòng đầu tiên trong booking
	room, err := l.svcCtx.RoomRepo.GetByID(ctx, bookingRooms[0].RoomID)
	if err != nil {
		logx.Error(err)
		return err
	}
	if room == nil {
		return errors.New("phòng không tồn tại")
	}

	var homestayID int = room.HomestayID

	// Tạo review
	review := &model.ReviewCreateRequest{
		BookingID:  req.BookingID,
		HomestayID: homestayID,
		UserID:     userID,
		Comment:    req.Comment,
		Rating:     req.Rating,
	}

	// Lưu review vào DB
	if _, err := l.svcCtx.BookingRepo.CreateReview(ctx, review); err != nil {
		logx.Error(err)
		return err
	}

	return nil
}

// Get payments by user ID
func (l *BookingLogic) FilterPayment(ctx context.Context, userId int, req *types.FilterPaymentReq) (types.FilterPaymentResp, error) {
	logx.Info(req)

	var bookingCode string
	if req.BookingCode != nil && *req.BookingCode != "" {
		bookingCode = *req.BookingCode
	}

	bookings, _, err := l.svcCtx.BookingRepo.FilterByBookingCode(ctx, userId, bookingCode, 1, 1000)
	if err != nil {
		logx.Error(err)
		return types.FilterPaymentResp{}, err
	}

	var bookingIds []int
	for _, booking := range bookings {
		bookingIds = append(bookingIds, booking.ID)
	}

	// 1. Mapping filter sang model.PaymentSearchRequest
	searchReq := &model.PaymentSearchRequest{
		BookingIds:    bookingIds,
		PaymentMethod: req.Method,
		Page:          req.Page,
		PageSize:      req.PageSize,
	}

	if req.DateFrom != nil && *req.DateTo != "" {
		t, _ := time.Parse("2006-01-02", *req.DateFrom)
		searchReq.StartDate = &t
	}

	if req.DateTo != nil && *req.DateTo != "" {
		t, _ := time.Parse("2006-01-02", *req.DateTo)
		searchReq.EndDate = &t
	}

	// 2. Lấy danh sách payment theo user ID
	payments, total, err := l.svcCtx.PaymentRepo.Search(ctx, searchReq)
	if err != nil {
		logx.Error(err)
		return types.FilterPaymentResp{}, err
	}

	var respPayments []types.Payment
	for _, payment := range payments {
		respPayments = append(respPayments, types.Payment{
			ID:            payment.ID,
			Amount:        payment.Amount,
			BookingID:     payment.BookingID,
			BookingCode:   payment.BookingCode,
			PaymentMethod: payment.PaymentMethod,
			PaymentStatus: payment.PaymentStatus,
			TransactionID: payment.TransactionID,
			PaymentDate:   payment.PaymentDate.Format("2006-01-02 15:04:05"),
		})
	}

	return types.FilterPaymentResp{
		Payments: respPayments,
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
	}, nil
}
