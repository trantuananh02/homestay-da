package svc

import (
	"homestay-be/cmd/config"
	"homestay-be/cmd/database/mysql"
	"homestay-be/cmd/database/repo"
	"homestay-be/cmd/mail"
	"homestay-be/cmd/storage"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config     config.Config
	DB         *sqlx.DB
	CldClient  *storage.CloudinaryClient
	MailClient *mail.MailClient // Thêm MailClient vào ServiceContext

	// Repositories
	UserRepo             repo.UserRepository
	HomestayRepo         repo.HomestayRepository
	RoomRepo             repo.RoomRepository
	RoomAvailabilityRepo repo.RoomAvailabilityRepository
	BookingRequestRepo   repo.BookingRequestRepository
	BookingRepo          repo.BookingRepository
	PaymentRepo          repo.PaymentRepository
	ReviewRepo           repo.ReviewRepository
}

func NewServiceContext(c config.Config) *ServiceContext {
	// Khởi tạo database connection
	dbConfig := &mysql.Config{
		Host:     c.Database.Host,
		Port:     c.Database.Port,
		User:     c.Database.User,
		Password: c.Database.Password,
		DBName:   c.Database.DBName,
		SSLMode:  c.Database.SSLMode,
	}

	db, err := mysql.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Khởi tạo repository factory
	repoFactory := mysql.NewRepositoryFactory(db)

	logx.Info(c.Http.Path, c.Storage.APIKey, c.Storage.APISecret)

	return &ServiceContext{
		Config:    c,
		DB:        db,
		CldClient: storage.NewCloudinaryClient(c.Storage.CloudName, c.Storage.APIKey, c.Storage.APISecret, "homestay"),
		MailClient: mail.NewMailClient(c.Mail.From, c.Mail.Username, c.Mail.Password), // Khởi tạo MailClient
		// Repositories
		UserRepo:             repoFactory.UserRepo,
		HomestayRepo:         repoFactory.HomestayRepo,
		RoomRepo:             repoFactory.RoomRepo,
		RoomAvailabilityRepo: repoFactory.RoomAvailabilityRepo,
		BookingRequestRepo:   repoFactory.BookingRequestRepo,
		BookingRepo:          repoFactory.BookingRepo,
		PaymentRepo:          repoFactory.PaymentRepo,
		ReviewRepo:           repoFactory.ReviewRepo,
	}
}

// Close đóng database connection
func (svc *ServiceContext) Close() error {
	if svc.DB != nil {
		return mysql.CloseConnection(svc.DB)
	}
	return nil
}
