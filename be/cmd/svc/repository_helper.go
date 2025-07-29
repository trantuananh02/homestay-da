package svc

import (
	"context"
	"homestay-be/cmd/database/model"
)

// RepositoryHelper cung cấp các phương thức tiện ích cho việc sử dụng repositories
type RepositoryHelper struct {
	svc *ServiceContext
}

// NewRepositoryHelper tạo instance mới của RepositoryHelper
func NewRepositoryHelper(svc *ServiceContext) *RepositoryHelper {
	return &RepositoryHelper{svc: svc}
}

// User helpers
func (h *RepositoryHelper) CreateUser(ctx context.Context, req *model.UserCreateRequest) (*model.User, error) {
	return h.svc.UserRepo.Create(ctx, req)
}

func (h *RepositoryHelper) GetUserByID(ctx context.Context, id int) (*model.User, error) {
	return h.svc.UserRepo.GetByID(ctx, id)
}

func (h *RepositoryHelper) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return h.svc.UserRepo.GetByEmail(ctx, email)
}

func (h *RepositoryHelper) UpdateUser(ctx context.Context, id int, req *model.UserUpdateRequest) (*model.User, error) {
	return h.svc.UserRepo.Update(ctx, id, req)
}

func (h *RepositoryHelper) DeleteUser(ctx context.Context, id int) error {
	return h.svc.UserRepo.Delete(ctx, id)
}

func (h *RepositoryHelper) ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int, error) {
	return h.svc.UserRepo.List(ctx, page, pageSize)
}

func (h *RepositoryHelper) SearchUsers(ctx context.Context, name, email, role string, page, pageSize int) ([]*model.User, int, error) {
	return h.svc.UserRepo.Search(ctx, name, email, role, page, pageSize)
}

// Homestay helpers
func (h *RepositoryHelper) CreateHomestay(ctx context.Context, req *model.HomestayCreateRequest) (*model.Homestay, error) {
	return h.svc.HomestayRepo.Create(ctx, req)
}

func (h *RepositoryHelper) GetHomestayByID(ctx context.Context, id int) (*model.Homestay, error) {
	return h.svc.HomestayRepo.GetByID(ctx, id)
}

func (h *RepositoryHelper) UpdateHomestay(ctx context.Context, id int, req *model.HomestayUpdateRequest) (*model.Homestay, error) {
	return h.svc.HomestayRepo.Update(ctx, id, req)
}

func (h *RepositoryHelper) DeleteHomestay(ctx context.Context, id int) error {
	return h.svc.HomestayRepo.Delete(ctx, id)
}

func (h *RepositoryHelper) ListHomestays(ctx context.Context, page, pageSize int) ([]*model.Homestay, int, error) {
	return h.svc.HomestayRepo.List(ctx, page, pageSize)
}

func (h *RepositoryHelper) SearchHomestays(ctx context.Context, req *model.HomestaySearchRequest) ([]*model.Homestay, int, error) {
	return h.svc.HomestayRepo.Search(ctx, req)
}

func (h *RepositoryHelper) GetHomestaysByOwner(ctx context.Context, ownerID int, page, pageSize int) ([]*model.Homestay, int, error) {
	return h.svc.HomestayRepo.GetByOwnerID(ctx, ownerID, page, pageSize)
}

// Room helpers
func (h *RepositoryHelper) CreateRoom(ctx context.Context, req *model.RoomCreateRequest) (*model.Room, error) {
	return h.svc.RoomRepo.Create(ctx, req)
}

func (h *RepositoryHelper) GetRoomByID(ctx context.Context, id int) (*model.Room, error) {
	return h.svc.RoomRepo.GetByID(ctx, id)
}

func (h *RepositoryHelper) UpdateRoom(ctx context.Context, id int, req *model.RoomUpdateRequest) (*model.Room, error) {
	return h.svc.RoomRepo.Update(ctx, id, req)
}

func (h *RepositoryHelper) DeleteRoom(ctx context.Context, id int) error {
	return h.svc.RoomRepo.Delete(ctx, id)
}

func (h *RepositoryHelper) ListRooms(ctx context.Context, page, pageSize int) ([]*model.Room, int, error) {
	return h.svc.RoomRepo.List(ctx, page, pageSize)
}

func (h *RepositoryHelper) SearchRooms(ctx context.Context, req *model.RoomSearchRequest) ([]*model.Room, int, error) {
	return h.svc.RoomRepo.Search(ctx, req)
}

func (h *RepositoryHelper) GetRoomsByHomestay(ctx context.Context, homestayID int, page, pageSize int) ([]*model.Room, int, error) {
	return h.svc.RoomRepo.GetByHomestayID(ctx, homestayID, page, pageSize)
}

func (h *RepositoryHelper) GetAvailableRooms(ctx context.Context, homestayID int, checkIn, checkOut string, numGuests int) ([]*model.Room, error) {
	return h.svc.RoomRepo.GetAvailableRooms(ctx, homestayID, checkIn, checkOut, numGuests)
}

// Room Availability helpers
func (h *RepositoryHelper) CreateRoomAvailability(ctx context.Context, req *model.RoomAvailabilityCreateRequest) (*model.RoomAvailability, error) {
	return h.svc.RoomAvailabilityRepo.Create(ctx, req)
}

func (h *RepositoryHelper) GetRoomAvailabilityByID(ctx context.Context, id int) (*model.RoomAvailability, error) {
	return h.svc.RoomAvailabilityRepo.GetByID(ctx, id)
}

func (h *RepositoryHelper) UpdateRoomAvailability(ctx context.Context, id int, req *model.RoomAvailabilityUpdateRequest) (*model.RoomAvailability, error) {
	return h.svc.RoomAvailabilityRepo.Update(ctx, id, req)
}

func (h *RepositoryHelper) DeleteRoomAvailability(ctx context.Context, id int) error {
	return h.svc.RoomAvailabilityRepo.Delete(ctx, id)
}

func (h *RepositoryHelper) ListRoomAvailabilities(ctx context.Context, page, pageSize int) ([]*model.RoomAvailability, int, error) {
	return h.svc.RoomAvailabilityRepo.List(ctx, page, pageSize)
}

func (h *RepositoryHelper) SearchRoomAvailabilities(ctx context.Context, req *model.RoomAvailabilitySearchRequest) ([]*model.RoomAvailability, int, error) {
	return h.svc.RoomAvailabilityRepo.Search(ctx, req)
}

func (h *RepositoryHelper) GetRoomAvailabilitiesByRoom(ctx context.Context, roomID int, page, pageSize int) ([]*model.RoomAvailability, int, error) {
	return h.svc.RoomAvailabilityRepo.GetByRoomID(ctx, roomID, page, pageSize)
}

func (h *RepositoryHelper) GetRoomAvailabilitiesByDateRange(ctx context.Context, roomID int, startDate, endDate string) ([]*model.RoomAvailability, error) {
	return h.svc.RoomAvailabilityRepo.GetByDateRange(ctx, roomID, startDate, endDate)
}

func (h *RepositoryHelper) CreateRoomAvailabilityBatch(ctx context.Context, req *model.RoomAvailabilityBatchRequest) error {
	return h.svc.RoomAvailabilityRepo.CreateBatch(ctx, req)
}

func (h *RepositoryHelper) CheckRoomAvailability(ctx context.Context, roomID int, checkIn, checkOut string) (bool, error) {
	return h.svc.RoomAvailabilityRepo.CheckAvailability(ctx, roomID, checkIn, checkOut)
}
