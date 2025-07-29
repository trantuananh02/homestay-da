package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// HTTP Status Codes
const (
	Success             = 0
	BadRequest          = 400
	Unauthorized        = 401
	Forbidden           = 403
	NotFound            = 404
	MethodNotAllowed    = 405
	Conflict            = 409
	UnprocessableEntity = 422
	TooManyRequests     = 429
	InternalServerError = 500
	BadGateway          = 502
	ServiceUnavailable  = 503
)

// Error Messages
const (
	// General errors
	MsgSuccess             = "Thành công"
	MsgBadRequest          = "Yêu cầu không hợp lệ"
	MsgUnauthorized        = "Không được phép truy cập"
	MsgForbidden           = "Không có quyền truy cập"
	MsgNotFound            = "Không tìm thấy"
	MsgInternalServerError = "Lỗi hệ thống"
	MsgServiceUnavailable  = "Dịch vụ không khả dụng"

	// Authentication errors
	MsgTokenRequired      = "Token không được cung cấp"
	MsgTokenInvalid       = "Token không hợp lệ"
	MsgTokenExpired       = "Token đã hết hạn"
	MsgInvalidCredentials = "Email hoặc mật khẩu không đúng"
	MsgUserNotFound       = "Người dùng không tồn tại"
	MsgUserAlreadyExists  = "Người dùng đã tồn tại"
	MsgPasswordMismatch   = "Mật khẩu không khớp"
	MsgInvalidRole        = "Vai trò không hợp lệ"

	// Validation errors
	MsgInvalidData     = "Dữ liệu không hợp lệ"
	MsgRequiredField   = "Trường này là bắt buộc"
	MsgInvalidEmail    = "Email không hợp lệ"
	MsgInvalidPhone    = "Số điện thoại không hợp lệ"
	MsgInvalidID       = "ID không hợp lệ"
	MsgInvalidDate     = "Ngày không hợp lệ"
	MsgInvalidPrice    = "Giá không hợp lệ"
	MsgInvalidCapacity = "Sức chứa không hợp lệ"

	// Homestay errors
	MsgHomestayNotFound          = "Homestay không tồn tại"
	MsgHomestayAlreadyExists     = "Homestay đã tồn tại"
	MsgHomestayNameRequired      = "Tên homestay là bắt buộc"
	MsgHomestayAddressRequired   = "Địa chỉ homestay là bắt buộc"
	MsgHomestayHasActiveBookings = "Không thể xóa homestay có booking đang hoạt động"
	MsgHomestayAccessDenied      = "Không có quyền truy cập homestay này"
	MsgHomestayUpdateFailed      = "Cập nhật homestay thất bại"
	MsgHomestayDeleteFailed      = "Xóa homestay thất bại"

	// Room errors
	MsgRoomNotFound          = "Phòng không tồn tại"
	MsgRoomAlreadyExists     = "Phòng đã tồn tại"
	MsgRoomNameRequired      = "Tên phòng là bắt buộc"
	MsgRoomTypeRequired      = "Loại phòng là bắt buộc"
	MsgRoomCapacityRequired  = "Sức chứa phòng là bắt buộc"
	MsgRoomPriceRequired     = "Giá phòng là bắt buộc"
	MsgRoomHasActiveBookings = "Không thể xóa phòng có booking đang hoạt động"
	MsgRoomAccessDenied      = "Không có quyền truy cập phòng này"
	MsgRoomUpdateFailed      = "Cập nhật phòng thất bại"
	MsgRoomDeleteFailed      = "Xóa phòng thất bại"

	// Availability errors
	MsgAvailabilityNotFound     = "Availability không tồn tại"
	MsgAvailabilityExists       = "Availability đã tồn tại cho ngày này"
	MsgAvailabilityDateRequired = "Ngày availability là bắt buộc"
	MsgAvailabilityInvalidDate  = "Ngày bắt đầu phải trước ngày kết thúc"
	MsgAvailabilityUpdateFailed = "Cập nhật availability thất bại"

	// Booking errors
	MsgBookingNotFound        = "Booking không tồn tại"
	MsgBookingAlreadyExists   = "Booking đã tồn tại"
	MsgBookingDateConflict    = "Ngày booking bị trùng"
	MsgBookingRoomUnavailable = "Phòng không khả dụng cho ngày này"
	MsgBookingInvalidDate     = "Ngày booking không hợp lệ"
	MsgBookingUpdateFailed    = "Cập nhật booking thất bại"
	MsgBookingDeleteFailed    = "Xóa booking thất bại"

	// File upload errors
	MsgFileUploadFailed = "Tải file thất bại"
	MsgFileTooLarge     = "File quá lớn"
	MsgInvalidFileType  = "Loại file không hợp lệ"
	MsgFileNotFound     = "File không tồn tại"

	// Database errors
	MsgDatabaseError            = "Lỗi cơ sở dữ liệu"
	MsgDatabaseConnectionFailed = "Kết nối cơ sở dữ liệu thất bại"
	MsgDatabaseQueryFailed      = "Truy vấn cơ sở dữ liệu thất bại"

	// API errors
	MsgAPIEndpointNotFound  = "API endpoint không tồn tại"
	MsgInvalidRequestMethod = "Phương thức request không hợp lệ"
	MsgRateLimitExceeded    = "Vượt quá giới hạn request"
	MsgRequestTimeout       = "Request timeout"
)

// ResponseResult chứa thông tin kết quả
type ResponseResult struct {
	Code    int    `json:"code" example:"0"`
	Message string `json:"message" example:"Thành công"`
}

// ResponseData chứa dữ liệu thực tế
type ResponseData struct {
	Result ResponseResult `json:"result"`
	Data   interface{}    `json:"data"`
}

// ResponseJSON trả về response với format chuẩn
func ResponseJSON(c *gin.Context, code int, message string, data interface{}) {
	response := ResponseData{
		Result: ResponseResult{
			Code:    code,
			Message: message,
		},
		Data: data,
	}
	c.JSON(http.StatusOK, response)
}

// ResponseSuccess trả về response thành công
func ResponseSuccess(c *gin.Context, data interface{}) {
	ResponseJSON(c, Success, MsgSuccess, data)
}

// ResponseError trả về response lỗi
func ResponseError(c *gin.Context, code int, message string) {
	ResponseJSON(c, code, message, nil)
}

// ResponseRedirect trả về redirect
func ResponseRedirect(c *gin.Context, statusCode int, destination string) {
	c.Redirect(statusCode, destination)
}

// ResponseBadRequest trả về lỗi 400
func ResponseBadRequest(c *gin.Context, message string) {
	ResponseError(c, BadRequest, message)
}

// ResponseUnauthorized trả về lỗi 401
func ResponseUnauthorized(c *gin.Context, message string) {
	ResponseError(c, Unauthorized, message)
}

// ResponseForbidden trả về lỗi 403
func ResponseForbidden(c *gin.Context, message string) {
	ResponseError(c, Forbidden, message)
}

// ResponseNotFound trả về lỗi 404
func ResponseNotFound(c *gin.Context, message string) {
	ResponseError(c, NotFound, message)
}

// ResponseInternalServerError trả về lỗi 500
func ResponseInternalServerError(c *gin.Context, message string) {
	ResponseError(c, InternalServerError, message)
}

// Helper functions for common error responses
func ResponseInvalidData(c *gin.Context, field string) {
	ResponseError(c, BadRequest, field+" "+MsgInvalidData)
}

func ResponseRequiredField(c *gin.Context, field string) {
	ResponseError(c, BadRequest, field+" "+MsgRequiredField)
}

func ResponseDatabaseError(c *gin.Context) {
	ResponseError(c, InternalServerError, MsgDatabaseError)
}

func ResponseFileUploadError(c *gin.Context) {
	ResponseError(c, BadRequest, MsgFileUploadFailed)
}

func ResponseAccessDenied(c *gin.Context) {
	ResponseError(c, Forbidden, MsgForbidden)
}

func ResponseResourceNotFound(c *gin.Context, resource string) {
	ResponseError(c, NotFound, resource+" "+MsgNotFound)
}