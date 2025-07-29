package logic

import (
	"context"
	"encoding/json"
	"errors"
	"homestay-be/cmd/database/model"
	"homestay-be/cmd/svc"
	"homestay-be/cmd/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HomestayLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHomestayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayLogic {
	return &HomestayLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateHomestay - Create a new homestay
func (h *HomestayLogic) CreateHomestay(req *types.CreateHomestayRequest, hostID int) (*types.HomestayDetailResponse, error) {
	modelReq := &model.HomestayCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		City:        req.City,
		District:    req.District,
		Ward:        req.Ward,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		OwnerID:     hostID,
	}
	created, err := h.svcCtx.HomestayRepo.Create(h.ctx, modelReq)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	resp := types.Homestay{
		ID:          created.ID,
		Name:        created.Name,
		Description: created.Description,
		Address:     created.Address,
		City:        created.City,
		District:    created.District,
		Ward:        created.Ward,
		Latitude:    created.Latitude,
		Longitude:   created.Longitude,
		HostID:      created.OwnerID,
		Status:      created.Status,
		CreatedAt:   created.CreatedAt,
		UpdatedAt:   created.UpdatedAt,
	}
	return &types.HomestayDetailResponse{Homestay: resp}, nil
}

// GetHomestayByID - Get homestay by ID
func (h *HomestayLogic) GetHomestayByID(homestayID, hostID int) (*types.HomestayDetailResponse, error) {
	found, err := h.svcCtx.HomestayRepo.GetByID(h.ctx, homestayID)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	if found.OwnerID != hostID {
		return nil, errors.New("Không có quyền truy cập homestay này")
	}

	// lấy review
	reviews, _, err := h.svcCtx.ReviewRepo.GetByHomestayID(h.ctx, homestayID, 1, 100)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	var reviewList []types.Review
	var totalRating int
	var totalReviews int
	var rating float64

	for _, review := range reviews {
		reviewList = append(reviewList, types.Review{
			ID:         review.ID,
			GuestID:    review.UserID,
			GuestName:  review.UserName,
			HomestayID: review.HomestayID,
			BookingID:  review.BookingID,
			Rating:     review.Rating,
			Comment:    review.Comment,
			CreatedAt:  review.CreatedAt,
		})
		totalRating += review.Rating
		totalReviews++
	}

	if totalReviews > 0 {
		rating = float64(totalRating) / float64(totalReviews)
	}

	resp := types.Homestay{
		ID:           found.ID,
		Name:         found.Name,
		Description:  found.Description,
		Address:      found.Address,
		City:         found.City,
		District:     found.District,
		Ward:         found.Ward,
		Latitude:     found.Latitude,
		Longitude:    found.Longitude,
		HostID:       found.OwnerID,
		Status:       found.Status,
		CreatedAt:    found.CreatedAt,
		UpdatedAt:    found.UpdatedAt,
		Reviews:      reviewList,
		Rating:       rating,
		TotalReviews: totalReviews,
	}
	return &types.HomestayDetailResponse{Homestay: resp}, nil
}

// GetHomestayList - Get list of homestays for a host
func (h *HomestayLogic) GetHomestayList(req *types.HomestayListRequest, hostID int) (*types.HomestayListResponse, error) {
	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	// Tạo search request
	searchReq := &model.HomestaySearchRequest{
		OwnerID:  &hostID,
		Page:     page,
		PageSize: pageSize,
	}

	// Thêm các filter nếu có
	if req.Search != "" {
		searchReq.Name = &req.Search
	}
	if req.City != "" {
		searchReq.City = &req.City
	}
	if req.District != "" {
		searchReq.District = &req.District
	}
	if req.Status != "" {
		searchReq.Status = &req.Status
	}

	homestays, total, err := h.svcCtx.HomestayRepo.Search(h.ctx, searchReq)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	respList := make([]types.Homestay, 0, len(homestays))
	for _, hst := range homestays {
		var rooms []types.Room
		roomModels, _, err := h.svcCtx.RoomRepo.GetByHomestayID(h.ctx, hst.ID, 0, 0)
		if err != nil {
			logx.Error(err)
			return nil, err
		}

		for _, rm := range roomModels {
			rooms = append(rooms, types.Room{
				ID:          rm.ID,
				Name:        rm.Name,
				Description: rm.Description,
				Type:        rm.Type,
				Capacity:    rm.Capacity,
				Price:       rm.Price,
				Status:      rm.Status,
				CreatedAt:   rm.CreatedAt,
				UpdatedAt:   rm.UpdatedAt,
				HomestayID:  rm.HomestayID,
			})
		}

		logx.Info("Rooms: ", rooms)

		// lấy review
		reviews, _, err := h.svcCtx.ReviewRepo.GetByHomestayID(h.ctx, hst.ID, 1, 100)
		if err != nil {
			logx.Error(err)
			return nil, err
		}

		var reviewList []types.Review
		var totalRating int
		var totalReviews int
		var rating float64

		for _, review := range reviews {
			reviewList = append(reviewList, types.Review{
				ID:         review.ID,
				GuestID:    review.UserID,
				GuestName:  review.UserName,
				HomestayID: review.HomestayID,
				BookingID:  review.BookingID,
				Rating:     review.Rating,
				Comment:    review.Comment,
				CreatedAt:  review.CreatedAt,
			})
			totalRating += review.Rating
			totalReviews++
		}

		if totalReviews > 0 {
			rating = float64(totalRating) / float64(totalReviews)
		}

		respList = append(respList, types.Homestay{
			ID:           hst.ID,
			Name:         hst.Name,
			Description:  hst.Description,
			Address:      hst.Address,
			City:         hst.City,
			District:     hst.District,
			Ward:         hst.Ward,
			Latitude:     hst.Latitude,
			Longitude:    hst.Longitude,
			HostID:       hst.OwnerID,
			Status:       hst.Status,
			Rooms:        rooms,
			CreatedAt:    hst.CreatedAt,
			UpdatedAt:    hst.UpdatedAt,
			Reviews:      reviewList,
			Rating:       rating,
			TotalReviews: totalReviews,
		})
	}
	totalPage := (total + pageSize - 1) / pageSize
	return &types.HomestayListResponse{
		Homestays: respList,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

// UpdateHomestay - Update homestay
func (h *HomestayLogic) UpdateHomestay(homestayID int, req *types.UpdateHomestayRequest, hostID int) (*types.HomestayDetailResponse, error) {
	// Kiểm tra quyền
	found, err := h.svcCtx.HomestayRepo.GetByID(h.ctx, homestayID)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	if found.OwnerID != hostID {
		return nil, errors.New("Không có quyền cập nhật homestay này")
	}

	modelReq := &model.HomestayUpdateRequest{
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		City:        req.City,
		District:    req.District,
		Ward:        req.Ward,
		Latitude:    req.Latitude,
		Longitude:   req.Longitude,
		Status:      req.Status,
	}

	updated, err := h.svcCtx.HomestayRepo.Update(h.ctx, homestayID, modelReq)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	resp := types.Homestay{
		ID:          updated.ID,
		Name:        updated.Name,
		Description: updated.Description,
		Address:     updated.Address,
		City:        updated.City,
		District:    updated.District,
		Ward:        updated.Ward,
		Latitude:    updated.Latitude,
		Longitude:   updated.Longitude,
		HostID:      updated.OwnerID,
		Status:      updated.Status,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}
	return &types.HomestayDetailResponse{Homestay: resp}, nil
}

// DeleteHomestay - Delete homestay
func (h *HomestayLogic) DeleteHomestay(homestayID, hostID int) error {
	found, err := h.svcCtx.HomestayRepo.GetByID(h.ctx, homestayID)
	if err != nil {
		logx.Error(err)
		return err
	}
	if found.OwnerID != hostID {
		return errors.New("Không có quyền xóa homestay này")
	}
	return h.svcCtx.HomestayRepo.Delete(h.ctx, homestayID)
}

// GetHomestayStats - Get homestay statistics for a host
func (h *HomestayLogic) GetHomestayStats(hostID int) (*types.HomestayStatsResponse, error) {
	// Lấy tất cả homestay của host
	homestays, _, err := h.svcCtx.HomestayRepo.GetByOwnerID(h.ctx, hostID, 1, 1000)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	totalHomestays := len(homestays)
	activeHomestays := 0
	totalRooms := 0
	availableRooms := 0
	totalBookings := 0
	totalRevenue := 0.0
	// bookingSeen := make(map[int]bool) // Đếm booking duy nhất

	for _, hst := range homestays {
		if hst.Status == "active" {
			activeHomestays++
		}
		// Lấy tất cả phòng của homestay này
		rooms, _, err := h.svcCtx.RoomRepo.GetByHomestayID(h.ctx, hst.ID, 1, 1000)
		if err != nil {
			logx.Error(err)
			return nil, err
		}
		totalRooms += len(rooms)
		for _, room := range rooms {
			if room.Status == "available" {
				availableRooms++
			}
			// Lấy tất cả booking của phòng này qua bảng booking_room
			bookings, _, err := h.svcCtx.BookingRepo.GetByRoom(h.ctx, room.ID, 1, 1000)
			if err != nil {
				logx.Error(err)
				continue
			}
			for _, booking := range bookings {
				switch booking.Status {
				case "confirmed", "completed":
					totalRevenue += booking.TotalAmount
				}
				totalBookings++
			}
		}
	}

	return &types.HomestayStatsResponse{
		TotalHomestays:  totalHomestays,
		ActiveHomestays: activeHomestays,
		TotalRooms:      totalRooms,
		AvailableRooms:  availableRooms,
		TotalBookings:   totalBookings,
		TotalRevenue:    totalRevenue,
	}, nil
}

// GetHomestayStatsByID - Get homestay statistics for a specific homestay
func (h *HomestayLogic) GetHomestayStatsByID(homestayID, hostID int) (*types.HomestayStatsResponse, error) {
	var availableRooms int
	var totalBookings int
	var totalRevenue float64

	_, err := h.svcCtx.HomestayRepo.GetByID(h.ctx, homestayID)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	rooms, totalRooms, err := h.svcCtx.RoomRepo.GetByHomestayID(h.ctx, homestayID, 1, 1000)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	for _, room := range rooms {
		if room.Status == "available" {
			availableRooms++
		}

		bookings, _, err := h.svcCtx.BookingRepo.GetByRoom(h.ctx, room.ID, 1, 1000)
		if err != nil {
			logx.Error(err)
			continue
		}
		for _, booking := range bookings {
			switch booking.Status {
			case "confirmed", "completed":
				totalRevenue += booking.TotalAmount
			}
			totalBookings++
		}
	}

	return &types.HomestayStatsResponse{
		TotalHomestays:  1, // Chỉ có 1 homestay được yêu cầu
		ActiveHomestays: 1,
		TotalRooms:      totalRooms,
		AvailableRooms:  availableRooms,
		TotalBookings:   totalBookings,
		TotalRevenue:    totalRevenue,
	}, nil
}

// ToggleHomestayStatus - Toggle homestay status between active and inactive
func (h *HomestayLogic) ToggleHomestayStatus(homestayID, hostID int) (*types.HomestayDetailResponse, error) {
	// Kiểm tra quyền
	found, err := h.svcCtx.HomestayRepo.GetByID(h.ctx, homestayID)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	if found.OwnerID != hostID {
		return nil, errors.New("Không có quyền thay đổi trạng thái homestay này")
	}

	// Xác định trạng thái mới
	newStatus := "inactive"
	if found.Status == "inactive" {
		newStatus = "active"
	}

	// Cập nhật trạng thái
	modelReq := &model.HomestayUpdateRequest{
		Status: &newStatus,
	}

	updated, err := h.svcCtx.HomestayRepo.Update(h.ctx, homestayID, modelReq)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	resp := types.Homestay{
		ID:          updated.ID,
		Name:        updated.Name,
		Description: updated.Description,
		Address:     updated.Address,
		City:        updated.City,
		District:    updated.District,
		Ward:        updated.Ward,
		Latitude:    updated.Latitude,
		Longitude:   updated.Longitude,
		HostID:      updated.OwnerID,
		Status:      updated.Status,
		CreatedAt:   updated.CreatedAt,
		UpdatedAt:   updated.UpdatedAt,
	}
	return &types.HomestayDetailResponse{Homestay: resp}, nil
}

// logic guest
func (h *HomestayLogic) GetPublicHomestayList(req *types.HomestayListRequest) (*types.HomestayListResponse, error) {
	logx.Info(req)

	page := req.Page
	if page < 1 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize < 1 {
		pageSize = 10
	}

	searchReq := &model.HomestaySearchRequest{
		Page:     page,
		PageSize: pageSize,
	}

	// Thêm các filter nếu có
	if req.Search != "" {
		searchReq.Name = &req.Search
	}
	if req.City != "" {
		searchReq.City = &req.City
	}
	if req.District != "" {
		searchReq.District = &req.District
	}
	if req.Status != "" {
		searchReq.Status = &req.Status
	}
	if req.CheckIn != "" {
		searchReq.CheckIn = &req.CheckIn
	}
	if req.CheckOut != "" {
		searchReq.CheckOut = &req.CheckOut
	}
	if req.Guests > 0 {
		searchReq.GuestCount = &req.Guests
	}

	var homestays []*model.Homestay
	var total int
	var err error

	if req.CheckIn != "" || req.CheckOut != "" || req.Guests > 0 {
		homestays, total, err = h.svcCtx.HomestayRepo.SearchAvailable(h.ctx, searchReq)
		} else {
		homestays, total, err = h.svcCtx.HomestayRepo.Search(h.ctx, searchReq)
	}


	if err != nil {
		logx.Error(err)
		return nil, err
	}

	respList := make([]types.Homestay, 0, len(homestays))
	for _, hst := range homestays {
		var rooms []types.Room
		roomModels, _, err := h.svcCtx.RoomRepo.GetByHomestayID(h.ctx, hst.ID, 1, 100)
		if err != nil {
			logx.Error(err)
			return nil, err
		}

		for _, rm := range roomModels {

			var images []string
			if len(rm.Images) > 0 {
				// unmarshal JSON string to slice
				err := json.Unmarshal([]byte(rm.Images), &images)
				if err != nil {
					logx.Error(err)
					return nil, err
				}
			}
			rooms = append(rooms, types.Room{
				ID:          rm.ID,
				Name:        rm.Name,
				Description: rm.Description,
				Capacity:    rm.Capacity,
				Type:        rm.Type,
				Price:       rm.Price,
				Status:      rm.Status,
				CreatedAt:   rm.CreatedAt,
				UpdatedAt:   rm.UpdatedAt,
				HomestayID:  rm.HomestayID,
				Images:      images,
			})
		}

		// lấy review
		reviews, _, err := h.svcCtx.ReviewRepo.GetByHomestayID(h.ctx, hst.ID, 1, 100)
		if err != nil {
			logx.Error(err)
			return nil, err
		}

		var reviewList []types.Review
		var totalRating int
		var totalReviews int
		var rating float64

		for _, review := range reviews {
			reviewList = append(reviewList, types.Review{
				ID:         review.ID,
				GuestID:    review.UserID,
				GuestName:  review.UserName,
				HomestayID: review.HomestayID,
				BookingID:  review.BookingID,
				Rating:     review.Rating,
				Comment:    review.Comment,
				CreatedAt:  review.CreatedAt,
			})
			totalRating += review.Rating
			totalReviews++
		}

		if totalReviews > 0 {
			rating = float64(totalRating) / float64(totalReviews)
		}

		respList = append(respList, types.Homestay{
			ID:           hst.ID,
			Name:         hst.Name,
			Description:  hst.Description,
			Address:      hst.Address,
			City:         hst.City,
			District:     hst.District,
			Ward:         hst.Ward,
			Latitude:     hst.Latitude,
			Longitude:    hst.Longitude,
			HostID:       hst.OwnerID,
			Status:       hst.Status,
			Rooms:        rooms,
			CreatedAt:    hst.CreatedAt,
			UpdatedAt:    hst.UpdatedAt,
			Reviews:      reviewList,
			Rating:       rating,
			TotalReviews: totalReviews,
		})
	}

	totalPage := (total + pageSize - 1) / pageSize
	return &types.HomestayListResponse{
		Homestays: respList,
		Total:     total,
		Page:      page,
		PageSize:  pageSize,
		TotalPage: totalPage,
	}, nil
}

func (h *HomestayLogic) GetPublicHomestayByID(homestayID int) (*types.HomestayDetailResponse, error) {
	found, err := h.svcCtx.HomestayRepo.GetByID(h.ctx, homestayID)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	var rooms []types.Room
	roomModels, _, err := h.svcCtx.RoomRepo.GetByHomestayID(h.ctx, homestayID, 1, 100)
	if err != nil {
		logx.Error(err)
		return nil, err
	}
	for _, rm := range roomModels {
		var images []string
		var amenities []string
		if len(rm.Images) > 0 {
			// unmarshal JSON string to slice
			err := json.Unmarshal([]byte(rm.Images), &images)
			if err != nil {
				logx.Error(err)
				return nil, err
			}
		}
		if len(rm.Amenities) > 0 {
			// unmarshal JSON string to slice
			err := json.Unmarshal([]byte(rm.Amenities), &amenities)
			if err != nil {
				logx.Error(err)
				return nil, err
			}
		}

		logx.Info("Room: ", rm.Images, rm.Amenities)

		rooms = append(rooms, types.Room{
			ID:          rm.ID,
			Name:        rm.Name,
			Description: rm.Description,
			Capacity:    rm.Capacity,
			Type:        rm.Type,
			Price:       rm.Price,
			Status:      rm.Status,
			CreatedAt:   rm.CreatedAt,
			UpdatedAt:   rm.UpdatedAt,
			HomestayID:  rm.HomestayID,
			Images:      images,
			Amenities:   amenities,
		})
	}

	// lấy review
	reviews, _, err := h.svcCtx.ReviewRepo.GetByHomestayID(h.ctx, homestayID, 1, 100)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	var reviewList []types.Review
	var totalRating int
	var totalReviews int
	var rating float64

	for _, review := range reviews {
		reviewList = append(reviewList, types.Review{
			ID:         review.ID,
			GuestID:    review.UserID,
			GuestName:  review.UserName,
			HomestayID: review.HomestayID,
			BookingID:  review.BookingID,
			Rating:     review.Rating,
			Comment:    review.Comment,
			CreatedAt:  review.CreatedAt,
		})
		totalRating += review.Rating
		totalReviews++
	}

	if totalReviews > 0 {
		rating = float64(totalRating) / float64(totalReviews)
	}

	resp := types.Homestay{
		ID:           found.ID,
		Name:         found.Name,
		Description:  found.Description,
		Address:      found.Address,
		City:         found.City,
		District:     found.District,
		Ward:         found.Ward,
		Latitude:     found.Latitude,
		Longitude:    found.Longitude,
		HostID:       found.OwnerID,
		Status:       found.Status,
		Rooms:        rooms,
		CreatedAt:    found.CreatedAt,
		UpdatedAt:    found.UpdatedAt,
		Reviews:      reviewList,
		Rating:       rating,
		TotalReviews: totalReviews,
	}
	return &types.HomestayDetailResponse{Homestay: resp}, nil
}

func (h *HomestayLogic) GetTopHomestays(limit int) ([]types.Homestay, error) {
	if limit <= 0 {
		limit = 8 // Default limit
	}

	homestays, err := h.svcCtx.HomestayRepo.GetTopHomestays(h.ctx, limit)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	respList := make([]types.Homestay, 0, len(homestays))
	for _, hst := range homestays {
		var rooms []types.Room
		roomModels, _, err := h.svcCtx.RoomRepo.GetByHomestayID(h.ctx, hst.ID, 1, 100)
		if err != nil {
			logx.Error(err)
			return nil, err
		}

		for _, rm := range roomModels {
			logx.Info("Room: ", rm.Images)
			var images []string

			if len(rm.Images) > 0 {
				// unmarshal JSON string to slice
				err := json.Unmarshal([]byte(rm.Images), &images)
				if err != nil {
					logx.Error(err)
					return nil, err
				}
			}

			rooms = append(rooms, types.Room{
				ID:          rm.ID,
				Name:        rm.Name,
				Description: rm.Description,
				Capacity:    rm.Capacity,
				Price:       rm.Price,
				Status:      rm.Status,
				CreatedAt:   rm.CreatedAt,
				UpdatedAt:   rm.UpdatedAt,
				HomestayID:  rm.HomestayID,
				Images:      images, // Thêm trường Images
			})
		}

		// lấy review
		reviews, _, err := h.svcCtx.ReviewRepo.GetByHomestayID(h.ctx, hst.ID, 1, 100)
		if err != nil {
			logx.Error(err)
			return nil, err
		}

		var reviewList []types.Review
		var totalRating int
		var totalReviews int
		var rating float64

		for _, review := range reviews {
			reviewList = append(reviewList, types.Review{
				ID:         review.ID,
				GuestID:    review.UserID,
				GuestName:  review.UserName,
				HomestayID: review.HomestayID,
				BookingID:  review.BookingID,
				Rating:     review.Rating,
				Comment:    review.Comment,
				CreatedAt:  review.CreatedAt,
			})
			totalRating += review.Rating
			totalReviews++
		}

		if totalReviews > 0 {
			rating = float64(totalRating) / float64(totalReviews)
		}

		respList = append(respList, types.Homestay{
			ID:           hst.ID,
			Name:         hst.Name,
			Description:  hst.Description,
			Address:      hst.Address,
			City:         hst.City,
			District:     hst.District,
			Ward:         hst.Ward,
			Latitude:     hst.Latitude,
			Longitude:    hst.Longitude,
			HostID:       hst.OwnerID,
			Status:       hst.Status,
			Rooms:        rooms,
			CreatedAt:    hst.CreatedAt,
			UpdatedAt:    hst.UpdatedAt,
			Reviews:      reviewList,
			Rating:       rating,
			TotalReviews: totalReviews,
		})
	}
	return respList, nil
}


// get reviews for a homestay
func (h *HomestayLogic) GetHomestayReviews(homestayID int) ([]types.Review, error) {
	reviews, _, err := h.svcCtx.ReviewRepo.GetByHomestayID(h.ctx, homestayID, 1, 100)
	if err != nil {
		logx.Error(err)
		return nil, err
	}

	var reviewList []types.Review
	for _, review := range reviews {
		reviewList = append(reviewList, types.Review{
			ID:         review.ID,
			GuestID:    review.UserID,
			GuestName:  review.UserName,
			HomestayID: review.HomestayID,
			BookingID:  review.BookingID,
			Rating:     review.Rating,
			Comment:    review.Comment,
			CreatedAt:  review.CreatedAt,
		})
	}

	return reviewList, nil
}