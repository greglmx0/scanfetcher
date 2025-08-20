package repository

import (
	"scanfetcher/internal/db"
	"scanfetcher/internal/domain"

	"gorm.io/gorm"
)

type ScanRepository struct {
	db *gorm.DB
}

func NewScanRepository(dbConn *gorm.DB) *ScanRepository {
	return &ScanRepository{
		db: dbConn,
	}
}

func (r *ScanRepository) Create(scan domain.Scan) error {
	dbScan := db.Scan{
		Name:         scan.Name,
		Url:          scan.Url,
		LastScanRead: scan.LastScanRead,
		WebsiteRef:   scan.WebsiteID,
	}
	return r.db.Create(&dbScan).Error
}

func (r *ScanRepository) GetByWebsiteID(websiteID int) ([]domain.Scan, error) {
	var dbScans []db.Scan
	if err := r.db.Where("website_ref = ?", websiteID).Find(&dbScans).Error; err != nil {
		return nil, err
	}

	var scans []domain.Scan
	for _, s := range dbScans {
		scans = append(scans, domain.Scan{
			ID:           int(s.ID),
			Name:         s.Name,
			Url:          s.Url,
			LastScanRead: s.LastScanRead,
			WebsiteID:    s.WebsiteRef,
		})
	}

	return scans, nil
}

func (r *ScanRepository) UpdateLastScanRead(scanID int, lastScanRead string) error {
	return r.db.Model(&db.Scan{}).Where("id = ?", scanID).Update("last_scan_read", lastScanRead).Error
}
