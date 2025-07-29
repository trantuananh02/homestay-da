package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"homestay-be/cmd/database/model"
	"homestay-be/cmd/svc"
	"homestay-be/cmd/types"
	"math"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type RoomLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRoomLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RoomLogic {
	return &RoomLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateRoom - Create a new room
func (r *RoomLogic) CreateRoom(req *types.CreateRoomRequest, hostID int) (*types.RoomDetailResponse, error) {
	logx.Info("Creating room:", req)

	// Kiểm tra quyền sở hữu homestay
	homestay, err := r.svcCtx.HomestayRepo.GetByID(r.ctx, req.HomestayID)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("homestay không tồn tại")
	}
	if homestay.OwnerID != hostID {
		logx.Error(err)
		return nil, fmt.Errorf("không có quyền tạo room cho homestay này")
	}

	// Tạo room request
	roomReq := &model.RoomCreateRequest{
		HomestayID:  req.HomestayID,
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Capacity:    req.Capacity,
		Price:       req.Price,
		PriceType:   req.PriceType,
		Images:      req.Images,
		Amenities:   req.Amenities,
	}

	// Tạo room
	room, err := r.svcCtx.RoomRepo.Create(r.ctx, roomReq)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("lỗi tạo room: %w", err)
	}

	// Tạo response
	response := &types.RoomDetailResponse{
		Room: types.Room{
			ID:          room.ID,
			HomestayID:  room.HomestayID,
			Name:        room.Name,
			Description: room.Description,
			Type:        room.Type,
			Capacity:    room.Capacity,
			Price:       room.Price,
			PriceType:   room.PriceType,
			Status:      room.Status,
			CreatedAt:   room.CreatedAt,
			UpdatedAt:   room.UpdatedAt,
		},
		Homestay: types.Homestay{
			ID:   homestay.ID,
			Name: homestay.Name,
		},
	}

	return response, nil
}

// GetRoomByID - Get room by ID
func (r *RoomLogic) GetRoomByID(roomID, hostID int) (*types.RoomDetailResponse, error) {
	// Lấy thông tin room
	room, err := r.svcCtx.RoomRepo.GetByID(r.ctx, roomID)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("room không tồn tại")
	}

	// Kiểm tra quyền sở hữu
	homestay, err := r.svcCtx.HomestayRepo.GetByID(r.ctx, room.HomestayID)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("homestay không tồn tại")
	}
	if homestay.OwnerID != hostID {
		logx.Error(err)
		return nil, fmt.Errorf("không có quyền truy cập room này")
	}

	// Lấy danh sách availability (nếu có)
	availabilities, _, err := r.svcCtx.RoomAvailabilityRepo.GetByRoomID(r.ctx, roomID, 1, 1000)
	if err != nil {
		// Không báo lỗi nếu không có availability
		availabilities = []*model.RoomAvailability{}
	}

	// Chuyển đổi Amenities string sang types []string
	var amenities []string
	if len(room.Amenities) > 0 {
		err = json.Unmarshal([]byte(room.Amenities), &amenities)
		if err != nil {
			logx.Error(err)
			return nil, err
		}
	}

	var images []string
	if len(room.Images) > 0 {
		err = json.Unmarshal([]byte(room.Images), &images)
		if err != nil {
			logx.Error(err)
			return nil, err
		}
	}

	// Tạo response
	response := &types.RoomDetailResponse{
		Room: types.Room{
			ID:          room.ID,
			HomestayID:  room.HomestayID,
			Name:        room.Name,
			Description: room.Description,
			Type:        room.Type,
			Capacity:    room.Capacity,
			Price:       room.Price,
			PriceType:   room.PriceType,
			Status:      room.Status,
			CreatedAt:   room.CreatedAt,
			UpdatedAt:   room.UpdatedAt,
			Amenities:   amenities,
			Images:      images,
		},
		Homestay: types.Homestay{
			ID:   homestay.ID,
			Name: homestay.Name,
		},
	}

	// Thêm availabilities vào response
	for _, avail := range availabilities {
		response.Availabilities = append(response.Availabilities, types.RoomAvailability{
			ID:        avail.ID,
			RoomID:    avail.RoomID,
			Date:      avail.Date,
			Status:    avail.Status,
			Price:     avail.Price,
			CreatedAt: avail.CreatedAt,
			UpdatedAt: avail.UpdatedAt,
		})
	}

	return response, nil
}

// GetRoomList - Get list of rooms for a homestay
func (r *RoomLogic) GetRoomList(req *types.RoomListRequest, hostID int) (*types.RoomListResponse, error) {
	// Kiểm tra quyền sở hữu homestay
	logx.Info(req.HomestayID)

	// Thiết lập giá trị mặc định cho pagination
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// Tạo search request
	searchReq := &model.RoomSearchRequest{
		HomestayID: &req.HomestayID,
		Status:     &req.Status,
		Type:       &req.Type,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	// Tìm kiếm rooms
	rooms, total, err := r.svcCtx.RoomRepo.Search(r.ctx, searchReq)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("lỗi tìm kiếm rooms: %w", err)
	}

	// Chuyển đổi sang types.Room
	var roomList []types.Room
	for _, room := range rooms {

		// Chuyển đổi Images string sang types []string
		var images []string
		if len(room.Images) > 0 {
			err = json.Unmarshal([]byte(room.Images), &images)
			if err != nil {
				logx.Error(err)
				return nil, err
			}
		}

		roomList = append(roomList, types.Room{
			ID:          room.ID,
			HomestayID:  room.HomestayID,
			Name:        room.Name,
			Description: room.Description,
			Type:        room.Type,
			Capacity:    room.Capacity,
			Price:       room.Price,
			PriceType:   room.PriceType,
			Status:      room.Status,
			CreatedAt:   room.CreatedAt,
			UpdatedAt:   room.UpdatedAt,
			Images:      images,
		})
	}

	// Tính tổng số trang
	totalPage := int(math.Ceil(float64(total) / float64(req.PageSize)))

	response := &types.RoomListResponse{
		Rooms:     roomList,
		Total:     total,
		Page:      req.Page,
		PageSize:  req.PageSize,
		TotalPage: totalPage,
	}

	return response, nil
}

// UpdateRoom - Update room
func (r *RoomLogic) UpdateRoom(roomID int, req *types.UpdateRoomRequest, hostID int) (*types.RoomDetailResponse, error) {
	logx.Info("Updating room:", roomID, req)

	// Lấy thông tin room hiện tại
	room, err := r.svcCtx.RoomRepo.GetByID(r.ctx, roomID)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("room không tồn tại")
	}

	// Kiểm tra quyền sở hữu
	homestay, err := r.svcCtx.HomestayRepo.GetByID(r.ctx, room.HomestayID)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("homestay không tồn tại")
	}
	if homestay.OwnerID != hostID {
		logx.Error(err)
		return nil, fmt.Errorf("không có quyền cập nhật room này")
	}

	// Tạo update request
	updateReq := &model.RoomUpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		Capacity:    req.Capacity,
		Price:       req.Price,
		PriceType:   req.PriceType,
		Status:      req.Status,
		Images:      req.Images,
		Amenities:   req.Amenities,
	}

	// Cập nhật room
	updatedRoom, err := r.svcCtx.RoomRepo.Update(r.ctx, roomID, updateReq)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("lỗi cập nhật room: %w", err)
	}

	// Tạo response
	response := &types.RoomDetailResponse{
		Room: types.Room{
			ID:          updatedRoom.ID,
			HomestayID:  updatedRoom.HomestayID,
			Name:        updatedRoom.Name,
			Description: updatedRoom.Description,
			Type:        updatedRoom.Type,
			Capacity:    updatedRoom.Capacity,
			Price:       updatedRoom.Price,
			PriceType:   updatedRoom.PriceType,
			Status:      updatedRoom.Status,
			CreatedAt:   updatedRoom.CreatedAt,
			UpdatedAt:   updatedRoom.UpdatedAt,
		},
		Homestay: types.Homestay{
			ID:   homestay.ID,
			Name: homestay.Name,
		},
	}

	return response, nil
}

// DeleteRoom - Delete room
func (r *RoomLogic) DeleteRoom(roomID, hostID int) error {
	// Lấy thông tin room
	room, err := r.svcCtx.RoomRepo.GetByID(r.ctx, roomID)
	if err != nil {
		logx.Error(err)
		return fmt.Errorf("room không tồn tại")
	}

	// Kiểm tra quyền sở hữu
	homestay, err := r.svcCtx.HomestayRepo.GetByID(r.ctx, room.HomestayID)
	if err != nil {
		logx.Error(err)
		return fmt.Errorf("homestay không tồn tại")
	}
	if homestay.OwnerID != hostID {
		logx.Error(err)
		return fmt.Errorf("không có quyền xóa room này")
	}
	// TODO: Kiểm tra booking đang hoạt động trước khi xóa room (chưa implement GetActiveBookingsByRoomID)
	// activeBookings, err := r.svcCtx.BookingRepo.GetActiveBookingsByRoomID(r.ctx, roomID)
	// if err != nil {
	// 	return fmt.Errorf("lỗi kiểm tra booking: %w", err)
	// }
	// if len(activeBookings) > 0 {
	// 	return fmt.Errorf("không thể xóa room vì có booking đang hoạt động")
	// }

	// Xóa room
	err = r.svcCtx.RoomRepo.Delete(r.ctx, roomID)
	if err != nil {
		logx.Error(err)
		return fmt.Errorf("lỗi xóa room: %w", err)
	}

	return nil
}

// CreateAvailability - Create room availability
func (r *RoomLogic) CreateAvailability(req *types.CreateAvailabilityRequest, hostID int) (interface{}, error) {
	// Lấy thông tin room
	room, err := r.svcCtx.RoomRepo.GetByID(r.ctx, req.RoomID)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("room không tồn tại")
	}

	// Kiểm tra quyền sở hữu
	homestay, err := r.svcCtx.HomestayRepo.GetByID(r.ctx, room.HomestayID)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("homestay không tồn tại")
	}
	if homestay.OwnerID != hostID {
		logx.Error(err)
		return nil, fmt.Errorf("không có quyền tạo availability cho room này")
	}

	// Kiểm tra ngày không được trong quá khứ
	if req.Date.Before(time.Now().Truncate(24 * time.Hour)) {
		logx.Error(err)
		return nil, fmt.Errorf("không thể tạo availability cho ngày trong quá khứ")
	}

	// Tạo availability request
	availabilityReq := &model.RoomAvailabilityCreateRequest{
		RoomID: req.RoomID,
		Date:   req.Date,
		Status: req.Status,
		Price:  req.Price,
	}

	// Tạo availability
	availability, err := r.svcCtx.RoomAvailabilityRepo.Create(r.ctx, availabilityReq)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("lỗi tạo availability: %w", err)
	}

	// Tạo response
	response := types.RoomAvailability{
		ID:        availability.ID,
		RoomID:    availability.RoomID,
		Date:      availability.Date,
		Status:    availability.Status,
		Price:     availability.Price,
		CreatedAt: availability.CreatedAt,
		UpdatedAt: availability.UpdatedAt,
	}

	return response, nil
}

// UpdateAvailability - Update room availability
func (r *RoomLogic) UpdateAvailability(availabilityID int, req *types.UpdateAvailabilityRequest, hostID int) (interface{}, error) {
	// Lấy thông tin availability
	availability, err := r.svcCtx.RoomAvailabilityRepo.GetByID(r.ctx, availabilityID)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("availability không tồn tại")
	}

	// Lấy thông tin room
	room, err := r.svcCtx.RoomRepo.GetByID(r.ctx, availability.RoomID)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("room không tồn tại")
	}

	// Kiểm tra quyền sở hữu
	homestay, err := r.svcCtx.HomestayRepo.GetByID(r.ctx, room.HomestayID)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("homestay không tồn tại")
	}
	if homestay.OwnerID != hostID {
		logx.Error(err)
		return nil, fmt.Errorf("không có quyền cập nhật availability này")
	}

	// Tạo update request
	updateReq := &model.RoomAvailabilityUpdateRequest{
		Status: req.Status,
		Price:  req.Price,
	}

	// Cập nhật availability
	updatedAvailability, err := r.svcCtx.RoomAvailabilityRepo.Update(r.ctx, availabilityID, updateReq)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("lỗi cập nhật availability: %w", err)
	}

	// Tạo response
	response := types.RoomAvailability{
		ID:        updatedAvailability.ID,
		RoomID:    updatedAvailability.RoomID,
		Date:      updatedAvailability.Date,
		Status:    updatedAvailability.Status,
		Price:     updatedAvailability.Price,
		CreatedAt: updatedAvailability.CreatedAt,
		UpdatedAt: updatedAvailability.UpdatedAt,
	}

	return response, nil
}

// BulkUpdateAvailability - Update multiple availabilities
func (r *RoomLogic) BulkUpdateAvailability(req *types.BulkAvailabilityRequest, hostID int) error {
	// Lấy thông tin room
	room, err := r.svcCtx.RoomRepo.GetByID(r.ctx, req.RoomID)
	if err != nil {
		logx.Error(err)
		return fmt.Errorf("room không tồn tại")
	}

	// Kiểm tra quyền sở hữu
	homestay, err := r.svcCtx.HomestayRepo.GetByID(r.ctx, room.HomestayID)
	if err != nil {
		logx.Error(err)
		return fmt.Errorf("homestay không tồn tại")
	}
	if homestay.OwnerID != hostID {
		logx.Error(err)
		return fmt.Errorf("không có quyền cập nhật availability cho room này")
	}

	// Kiểm tra ngày bắt đầu không được trong quá khứ
	if req.StartDate.Before(time.Now().Truncate(24 * time.Hour)) {
		logx.Error(err)
		return fmt.Errorf("không thể cập nhật availability cho ngày trong quá khứ")
	}

	// Tạo map exclude dates để kiểm tra nhanh
	excludeMap := make(map[time.Time]bool)
	for _, date := range req.ExcludeDates {
		excludeMap[date.Truncate(24*time.Hour)] = true
	}

	// Tạo availabilities cho từng ngày trong khoảng
	currentDate := req.StartDate
	for currentDate.Before(req.EndDate) || currentDate.Equal(req.EndDate) {
		// Bỏ qua ngày trong exclude list
		if !excludeMap[currentDate.Truncate(24*time.Hour)] {
			// Kiểm tra xem availability đã tồn tại chưa
			// TODO: Bổ sung hàm kiểm tra tồn tại availability theo ngày nếu cần (chưa implement GetByRoomIDAndDate)
			// existingAvailability, err := r.svcCtx.RoomAvailabilityRepo.GetByRoomIDAndDate(r.ctx, req.RoomID, currentDate)
			// if err != nil {
			// 	// Nếu không tồn tại, tạo mới
			// 	availabilityReq := &model.RoomAvailabilityCreateRequest{
			// 		RoomID: req.RoomID,
			// 		Date:   currentDate,
			// 		Status: req.Status,
			// 		Price:  req.Price,
			// 	}
			// 	_, err = r.svcCtx.RoomAvailabilityRepo.Create(r.ctx, availabilityReq)
			// 	if err != nil {
			// 		return fmt.Errorf("lỗi tạo availability cho ngày %s: %w", currentDate.Format("2006-01-02"), err)
			// 	}
			// } else {
			// 	// Nếu đã tồn tại, cập nhật
			// 	updateReq := &model.RoomAvailabilityUpdateRequest{
			// 		Status: &req.Status,
			// 		Price:  req.Price,
			// 	}
			// 	_, err = r.svcCtx.RoomAvailabilityRepo.Update(r.ctx, existingAvailability.ID, updateReq)
			// 	if err != nil {
			// 		return fmt.Errorf("lỗi cập nhật availability cho ngày %s: %w", currentDate.Format("2006-01-02"), err)
			// 	}
			// }
		}

		// Tăng ngày lên 1
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return nil
}

// GetRoomStats - Get room statistics for a homestay
func (r *RoomLogic) GetRoomStats(homestayID, hostID int) (*types.RoomStatsResponse, error) {
	// Kiểm tra quyền sở hữu homestay
	homestay, err := r.svcCtx.HomestayRepo.GetByID(r.ctx, homestayID)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("homestay không tồn tại")
	}
	if homestay.OwnerID != hostID {
		logx.Error(err)
		return nil, fmt.Errorf("không có quyền truy cập thống kê homestay này")
	}

	// Lấy danh sách tất cả rooms của homestay
	rooms, _, err := r.svcCtx.RoomRepo.GetByHomestayID(r.ctx, homestayID, 1, 1000) // Lấy tất cả rooms
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("lỗi lấy danh sách rooms: %w", err)
	}

	// Tính toán thống kê
	stats := &types.RoomStatsResponse{
		TotalRooms:       len(rooms),
		AvailableRooms:   0,
		OccupiedRooms:    0,
		MaintenanceRooms: 0,
		AveragePrice:     0,
		TotalRevenue:     0,
		OccupancyRate:    0,
	}

	var totalPrice float64
	var availableCount int

	for _, room := range rooms {
		totalPrice += room.Price

		switch room.Status {
		case "available":
			stats.AvailableRooms++
			availableCount++
		case "occupied":
			stats.OccupiedRooms++
		case "maintenance":
			stats.MaintenanceRooms++
		}
	}

	// Tính average price
	if len(rooms) > 0 {
		stats.AveragePrice = totalPrice / float64(len(rooms))
	}

	// Tính occupancy rate
	if len(rooms) > 0 {
		stats.OccupancyRate = float64(stats.OccupiedRooms) / float64(len(rooms)) * 100
	}

	// TODO: Tính total revenue từ booking history
	// Đây là placeholder, cần implement logic tính revenue thực tế
	stats.TotalRevenue = 0

	return stats, nil
}

func (r *RoomLogic) GetPublicRoomList(req *types.RoomListRequest) (*types.RoomListResponse, error) {
	logx.Infof("Lấy danh sách phòng cho homestay ID: %d", req.HomestayID)

	// Lấy thông tin homestay (không kiểm tra quyền sở hữu)
	homestay, err := r.svcCtx.HomestayRepo.GetByID(r.ctx, req.HomestayID)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("homestay không tồn tại")
	}

	// (Tuỳ chọn) Chỉ hiển thị nếu homestay ở trạng thái hoạt động
	if homestay.Status != "active" {
		return nil, fmt.Errorf("homestay không khả dụng")
	}

	// Pagination mặc định
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}
	if req.PageSize > 100 {
		req.PageSize = 100
	}

	// Tạo search request
	searchReq := &model.RoomSearchRequest{
		HomestayID: &req.HomestayID,
		Status:     &req.Status,
		Type:       &req.Type,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	// Tìm kiếm phòng
	rooms, total, err := r.svcCtx.RoomRepo.Search(r.ctx, searchReq)
	if err != nil {
		logx.Error(err)
		return nil, fmt.Errorf("lỗi tìm kiếm rooms: %w", err)
	}

	// Chuyển sang types.Room
	var roomList []types.Room
	for _, room := range rooms {
		roomList = append(roomList, types.Room{
			ID:          room.ID,
			HomestayID:  room.HomestayID,
			Name:        room.Name,
			Description: room.Description,
			Type:        room.Type,
			Capacity:    room.Capacity,
			Price:       room.Price,
			PriceType:   room.PriceType,
			Status:      room.Status,
			CreatedAt:   room.CreatedAt,
			UpdatedAt:   room.UpdatedAt,
		})
	}

	// Tổng số trang
	totalPage := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &types.RoomListResponse{
		Rooms:     roomList,
		Total:     total,
		Page:      req.Page,
		PageSize:  req.PageSize,
		TotalPage: totalPage,
	}, nil
}
