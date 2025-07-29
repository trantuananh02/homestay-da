package handler

import (
	"homestay-be/cmd/logic"
	"homestay-be/cmd/svc"
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

	// Gọi logic xử lý upload
	logic := logic.NewUploadFileLogic(ctx, h.svcCtx, file)
	resp, err := logic.UploadFile(ctx, 1)
	if err != nil {
		response.ResponseError(c, response.InternalServerError, err.Error())
		return
	}

	response.ResponseSuccess(c, resp)
}
