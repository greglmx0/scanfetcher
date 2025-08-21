package db

import (
	"log"
	"os"
	"path/filepath"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	// Use the modernc.org/sqlite pure-Go driver (no cgo).
	_ "modernc.org/sqlite"
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

	_ = os.Chown(dbPath, os.Getuid(), os.Getgid())

	// create the file if it doesn't exist in dbPath
	dbName = filepath.Join(dbPath, dbName)
	log.Printf("Initialisation de la base de données à l'emplacement: %s\n", dbName)
	if _, err := os.Stat(dbName); os.IsNotExist(err) {
		file, err := os.Create(dbName)
		if err != nil {
			return nil, err
		}
		file.Close()
		log.Printf("Fichier de base de données créé: %s\n", dbName)
	} else {
		log.Printf("Fichier de base de données déjà existant: %s\n", dbName)
	}

	// Use the provided dbName as the sqlite file. modernc.org/sqlite registers
	// the driver under the name "sqlite" so GORM's sqlite driver will use it.
	dsn := "file:" + dbName + "?_busy_timeout=5000"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
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
