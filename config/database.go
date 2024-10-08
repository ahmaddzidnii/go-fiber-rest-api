package config

import (
	"fmt"

	"github.com/ahmaddzidnii/go-fiber-rest-api/models"
	// "gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB;

func ConnectDatabase() {
	// mendefinisikan dsn (Data Source Name)
	// dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", Config("DATABASE_USER"), Config("DATABASE_PASSWORD"), Config("DATABASE_HOST"), Config("DATABASE_PORT"), Config("DATABASE_NAME"));

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", Config("DATABASE_HOST"), Config("DATABASE_USER"), Config("DATABASE_PASSWORD"), Config("DATABASE_NAME"), Config("DATABASE_PORT"));
	
	db,err := gorm.Open(postgres.Open(dsn), &gorm.Config{});


	// Cek jika terjadi error saat koneksi ke database
	if(err != nil){
		panic(err);
	}

	// Migrasi database dan cek jika terjadi error
	if err:= db.AutoMigrate(&models.Book{}, &models.User{}); err != nil {
		panic(err);
	} else {
		fmt.Println("Database migrated");
	}

	DB = db;
}