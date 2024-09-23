package config

import (
	"fmt"

	"github.com/ahmaddzidnii/go-fiber-rest-api/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB;

func ConnectDatabase() {
	// Declare dsn (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", Config("DATABASE_USER"), Config("DATABASE_PASSWORD"), Config("DATABASE_HOST"), Config("DATABASE_PORT"), Config("DATABASE_NAME"));
	
	db,err := gorm.Open(mysql.Open(dsn), &gorm.Config{});

	if(err != nil){
		panic(err);
	}

	db.AutoMigrate(&models.Book{}, &models.User{});

	DB = db;
}