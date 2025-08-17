package handler

import (
	"homestay-be/cmd/logic"
	"homestay-be/cmd/svc"
	"homestay-be/cmd/types"
	"homestay-be/core/response"

	"github.com/gin-gonic/gin"
)

type UploadFileHandler struct {
	svcCtx *svc.ServiceContext
}

func NewUploadFileHandler(svcCtx *svc.ServiceContext) *UploadFileHandler {
	return &UploadFileHandler{
		svcCtx: svcCtx,
	}
}

func (h *UploadFileHandler) UploadFile(c *gin.Context) {
	ctx := c.Request.Context()

	// Lấy file từ form-data
	file, err := c.FormFile("image")
	if err != nil {
		response.ResponseError(c, response.BadRequest, "File không hợp lệ")
		return
	}

	// Lấy userID từ context (nếu có authentication)
	userID := int64(1) // Default value
	if userInterface, exists := c.Get("user"); exists {
		if user, ok := userInterface.(*types.UserInfo); ok {
			userID = int64(user.ID)
		}
	}

	// Gọi logic xử lý upload
	logic := logic.NewUploadFileLogic(ctx, h.svcCtx, file)
	resp, err := logic.UploadFile(ctx, userID)
	if err != nil {
		response.ResponseError(c, response.InternalServerError, err.Error())
		return
	}

	response.ResponseSuccess(c, resp)
}
