package fixtures

import (
	"log"

	"scanfetcher/internal/db"

	"gorm.io/gorm"
)

func SeedData(db *gorm.DB) {
	seedWebsites(db)
	seedScans(db)
}

func seedWebsites(dbConn *gorm.DB) {
	websites := []db.Website{
		{Name: "mangas-origines.fr"},
	}

	for _, w := range websites {
		var existing db.Website
		if err := dbConn.Where("name = ?", w.Name).First(&existing).Error; err == gorm.ErrRecordNotFound {
			if err := dbConn.Create(&w).Error; err != nil {
				log.Printf("❌ Failed to create website %s: %v", w.Name, err)
			} else {
				log.Printf("✅ Website seeded: %s", w.Name)
			}
		}
	}
}

func seedScans(dbConn *gorm.DB) {
	// get website mangas-origines.fr add a scan with the name "the-beginning-after-the-end" and the url "https://mangas-origines.fr/oeuvre/640-the-beginning-after-the-end/" 

	website := db.Website{}
	if err := dbConn.Where("name = ?", "mangas-origines.fr").First(&website).Error; err != nil {
		log.Printf("❌ Failed to find website: %v", err)
		return
	}
	// Check if the scan already exists
	var existingScan db.Scan
	if err := dbConn.Where("name = ?", "the-beginning-after-the-end").First(&existingScan).Error; err == gorm.ErrRecordNotFound {
		// Create the scan
		scan := db.Scan{
			Name:        "the-beginning-after-the-end",
			Url:         "https://mangas-origines.fr/oeuvre/the-beginning-after-the-end/",
			WebsiteRef:   website.ID,
			LastScanRead: "211",
		}
		if err := dbConn.Create(&scan).Error; err != nil {
			log.Printf("❌ Failed to create scan %s: %v", scan.Name, err)
		} else {
			log.Printf("✅ Scan seeded: %s", scan.Name)
		}
	} else {
		log.Printf("✅ Scan already exists: %s", existingScan.Name)
	}

}
