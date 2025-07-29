package server

import (
	"fmt"
	"log"
	"homestay-be/cmd/config"
	"homestay-be/cmd/handler"
	"homestay-be/cmd/svc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	config *config.Config
	router *gin.Engine
	svcCtx *svc.ServiceContext
}

// NewServer tạo instance mới của server
func NewServer(cfg *config.Config) *Server {
	// Set Gin mode
	if cfg.Http.Path == "localhost" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	
	router := gin.Default()
	
	// Thêm middleware cơ bản
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	
	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})
	
	// Tạo service context
	svcCtx := svc.NewServiceContext(*cfg)
	
	return &Server{
		config: cfg,
		router: router,
		svcCtx: svcCtx,
	}
}

// SetupRoutes thiết lập các routes cho server bằng cách sử dụng handler
func (s *Server) SetupRoutes() {
	// Đăng ký tất cả handlers từ handler/router.go
	handler.RegisterHandlers(s.router, s.svcCtx)
}

// Start khởi động server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%s", s.config.Http.Path, s.config.Http.Port)
	log.Printf("Server đang khởi động tại: %s", addr)
	
	return s.router.Run(addr)
}

// GetRouter trả về router instance (để test)
func (s *Server) GetRouter() *gin.Engine {
	return s.router
}

// GetServiceContext trả về service context
func (s *Server) GetServiceContext() *svc.ServiceContext {
	return s.svcCtx
} 