package handler

import (
	"context"
	"homestay-be/cmd/logic"
	"homestay-be/cmd/middleware"
	"homestay-be/cmd/svc"
	"homestay-be/core/response"

	"github.com/gin-gonic/gin"
)

// RegisterHandlers đăng ký tất cả các routes và handlers
func RegisterHandlers(router *gin.Engine, serverCtx *svc.ServiceContext) {
	// Middleware
	router.Use(middleware.CORSMiddleware())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		response.ResponseSuccess(c, gin.H{"status": "ok"})
	})
	// Upload file handler
	uploadHandler := NewUploadFileHandler(serverCtx)
	router.POST("/upload", uploadHandler.UploadFile)

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public)
		authHandler := NewAuthHandler(serverCtx)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Public routes
		homestayHandler := NewHomestayHandler(logic.NewHomestayLogic(context.Background(), serverCtx))
		public := api.Group("/public")
		{
			public.GET("/homestays/top", homestayHandler.GetTopHomestays)
			public.GET("/homestays", homestayHandler.GetPublicHomestayList)
			public.GET("/homestays/:id", homestayHandler.GetPublicHomestayByID)
		}

		// Protected routes (cần authentication)
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(serverCtx))
		{
			// User profile
			protected.GET("/auth/profile", authHandler.GetProfile)
			protected.PUT("/auth/profile", authHandler.UpdateProfile)
			protected.POST("/auth/logout", authHandler.Logout)
		}

		// Host routes (cần role host hoặc admin)
		host := api.Group("/host")
		host.Use(middleware.AuthMiddleware(serverCtx))
		host.Use(middleware.RoleMiddleware("host", "admin"))
		{
			// Init context for logic layers
			ctx := context.Background()

			// Initialize logic layers
			homestayLogic := logic.NewHomestayLogic(ctx, serverCtx)
			roomLogic := logic.NewRoomLogic(ctx, serverCtx)
			bookingLogic := logic.NewBookingLogic(ctx, serverCtx)

			// Initialize handlers
			homestayHandler := NewHomestayHandler(homestayLogic)
			roomHandler := NewRoomHandler(roomLogic)
			bookingHandler := NewBookingHandler(bookingLogic)

			// Homestay management
			host.GET("/homestays", homestayHandler.GetHomestayList)
			host.POST("/homestays", homestayHandler.CreateHomestay)
			host.GET("/homestays/stats", homestayHandler.GetHomestayStats)
			host.GET("/homestays/:id", homestayHandler.GetHomestayByID)
			host.PUT("/homestays/:id", homestayHandler.UpdateHomestay)
			host.PUT("/homestays/:id/toggle-status", homestayHandler.ToggleHomestayStatus)
			host.DELETE("/homestays/:id", homestayHandler.DeleteHomestay)
			host.GET("/homestays/:id/stats", homestayHandler.GetHomestayStatsByID)
			host.GET("/homestays/:id/review", homestayHandler.GetHomestayReviews)

			// Room management
			host.POST("/rooms", roomHandler.CreateRoom)
			host.GET("/rooms", roomHandler.GetRoomList)
			host.GET("/rooms/:id", roomHandler.GetRoomByID)
			host.PUT("/rooms/:id", roomHandler.UpdateRoom)
			host.DELETE("/rooms/:id", roomHandler.DeleteRoom)

			// Room availability
			host.POST("/rooms/availability", roomHandler.CreateAvailability)
			host.PUT("/rooms/availability/:id", roomHandler.UpdateAvailability)

			// Room statistics
			host.GET("/homestays/:id/rooms/stats", roomHandler.GetRoomStats)

			// Booking requests
			host.GET("/booking", bookingHandler.FilterBookings)
			host.POST("/booking", bookingHandler.CreateBooking)
			host.GET("/booking/:id", bookingHandler.GetBookingDetail)
			host.GET("/homestays/:id/bookings", bookingHandler.GetBookingsByHomestayID)
			host.PUT("/booking/:id/status", bookingHandler.UpdateStatusBooking)

			// Payment management
			host.GET("/payments", bookingHandler.GetPayments)
		}

		// Guest routes (cần role guest)
		guest := api.Group("/guest")
		guest.Use(middleware.AuthMiddleware(serverCtx))
		guest.Use(middleware.RoleMiddleware("guest"))
		{
			ctx := context.Background()

			// Initialize logic layers
			homestayLogic := logic.NewHomestayLogic(ctx, serverCtx)
			roomLogic := logic.NewRoomLogic(ctx, serverCtx)
			bookingLogic := logic.NewBookingLogic(ctx, serverCtx)

			// Initialize handlers
			homestayHandler := NewHomestayHandler(homestayLogic)
			roomHandler := NewRoomHandler(roomLogic)
			bookingHandler := NewBookingHandler(bookingLogic)

			// Homestay management
			guest.GET("/homestays", homestayHandler.GetPublicHomestayList)
			guest.GET("/homestays/:id", homestayHandler.GetPublicHomestayByID)
			guest.GET("/rooms", roomHandler.GetRoomList)
			guest.GET("/homestays/:id/bookings", bookingHandler.GetBookingsByHomestayID)
			guest.GET("/booking/:id", bookingHandler.GetBookingDetail)
			guest.POST("/booking", bookingHandler.CreateGuestBooking)
			guest.GET("/booking", bookingHandler.GetGuestBookings)
			guest.PUT("/booking/:id/status", bookingHandler.UpdateStatusBooking)
			
			guest.POST("/review", bookingHandler.CreateReview)
		}
	}

	// 404 handler
	router.NoRoute(func(c *gin.Context) {
		response.ResponseError(c, response.NotFound, response.MsgAPIEndpointNotFound)
	})
}
