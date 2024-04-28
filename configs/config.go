package configs

import (
	"fmt"
	"log"
	"os"
	"sistem_peminjaman_be/models"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (*gorm.DB, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	location, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		// Handle the error
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	dbConn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	dbConn = dbConn.Session(&gorm.Session{
		NowFunc: func() time.Time {
			return time.Now().In(location)
		},
	})

	return dbConn, nil

}

func AccountSeeder(db *gorm.DB) error {
	password := "$2a$10$QXBNiEWub5z3TX5LFewSy.atj0iARk1vCZDgzRQTDp5xOQopj4WRW"
	users := []models.User{
		{
			FullName:       "Admin",
			Email:          "admin@gmail.com",
			Password:       password,
			NIMNIP:         "25245802",
			ProfilePicture: "https://icon-library.com/images/default-user-icon/default-user-icon-13.jpg",
			Role:           "admin",
		},
		{
			FullName:       "User",
			Email:          "user@gmail.com",
			Password:       password,
			NIMNIP:         "25245804",
			ProfilePicture: "https://icon-library.com/images/default-user-icon/default-user-icon-13.jpg",
			Role:           "user",
		},
	}

	for _, user := range users {

		// Check if data already exists
		var count int64
		if err := db.Model(&models.User{}).Where(&user).Count(&count).Error; err != nil {
			return err
		}

		// If data exists, skip seeding
		if count > 0 {
			continue
		}
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}

	return nil
}

func MigrateDB(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.User{},
		&models.Notification{},
		&models.TemplateMessage{},
		&models.Lab{},
		&models.LabImage{},
		&models.HistorySearch{},
		&models.HistorySeenLab{},
	)
}
