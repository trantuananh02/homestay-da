package logic

import (
	"context"
	"homestay-be/cmd/svc"
	"homestay-be/cmd/types"
	"mime/multipart"

	"github.com/zeromicro/go-zero/core/logx"
)

type UploadFileLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	file   *multipart.FileHeader
}

func NewUploadFileLogic(ctx context.Context, svcCtx *svc.ServiceContext, file *multipart.FileHeader) *UploadFileLogic {
	return &UploadFileLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		file:   file,
	}
}

func (l *UploadFileLogic) UploadFile(ctx context.Context, userID int64) (*types.UploadFileRes, error) {
	file, err := l.file.Open()
	if err != nil {
		logx.Error("UploadFileLogic.UploadFile: failed to open file", err)
		return nil, err
	}

	defer file.Close()

	url, err := l.svcCtx.CldClient.UploadImage(l.ctx, file, userID)
	if err != nil {
		logx.Error("UploadFileLogic.UploadFile: failed to upload image", err)
		return nil, err
	}

	return &types.UploadFileRes{
		Url: url,
	}, nil
}
