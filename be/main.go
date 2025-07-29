package main

import (
	"homestay-be/cmd/config"
	"homestay-be/cmd/server"
	"log"
	"os"
)

func main() {
	// Đường dẫn đến file config
	configPath := "configs/config.yaml"

	// Kiểm tra xem file config có tồn tại không
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("File config không tồn tại: %s", configPath)
	}

	// Đọc config
	conf, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Không thể đọc config: %v", err)
	}

	log.Printf("Config đã được đọc thành công")
	log.Printf("Server sẽ chạy tại: %s:%s", conf.Http.Path, conf.Http.Port)
	log.Printf("Database: %s", conf.Database.Driver)

	// Tạo server mới
	srv := server.NewServer(conf)

	// Thiết lập routes
	srv.SetupRoutes()

	// Khởi động server
	if err := srv.Start(); err != nil {
		log.Fatalf("Không thể khởi động server: %v", err)
	}
}
