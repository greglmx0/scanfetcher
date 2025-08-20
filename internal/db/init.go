package db

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Username string `gorm:"not null"`
}

type Website struct {
	gorm.Model
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
}

type Scan struct {
	gorm.Model
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"not null"`
	Url          string `gorm:"not null"`
	LastScanRead string `gorm:"not null"`
	WebsiteRef   uint
	Website      Website `gorm:"foreignKey:WebsiteRef"`
}

// InitDB initialise et retourne une connexion à SQLite avec GORM
func InitDB(dbPath string, dbName string) (*gorm.DB, error) {

	dir := filepath.Dir(dbPath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// change owner of the file to the current user
	if err := os.Chown(dbPath, os.Getuid(), os.Getgid()); err != nil {
		return nil, err
	}

	// create the file if it doesn't exist
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		file, err := os.Create(dbName)
		if err != nil {
			return nil, err
		}
		defer file.Close()
	}

	db, err := gorm.Open(sqlite.Open(dbPath+"/"+dbName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&User{}, &Website{}, &Scan{})
	if err != nil {
		return nil, err
	}

	log.Println("Base de données initialisée avec succès 🚀")
	return db, nil
}
