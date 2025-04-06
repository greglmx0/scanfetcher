package repository

import (
	"scanfetcher/internal/db"
	"scanfetcher/internal/domain"

	"gorm.io/gorm"
)

type WebsiteRepository struct {
	db *gorm.DB
}

func NewWebsiteRepository(dbConn *gorm.DB) *WebsiteRepository {
	return &WebsiteRepository{
		db: dbConn,
	}
}

func (r *WebsiteRepository) Create(website domain.Website) error {
	dbWebsite := db.Website{
		Name: website.Name,
	}
	return r.db.Create(&dbWebsite).Error
}

func (r *WebsiteRepository) GetByID(id int) (*domain.Website, error) {
	var dbWebsite db.Website
	if err := r.db.First(&dbWebsite, id).Error; err != nil {
		return nil, err
	}

	return &domain.Website{
		ID:   int(dbWebsite.ID),
		Name: dbWebsite.Name,
	}, nil
}

func (r *WebsiteRepository) GetAll() ([]domain.Website, error) {
	var dbWebsites []db.Website
	if err := r.db.Find(&dbWebsites).Error; err != nil {
		return nil, err
	}

	var websites []domain.Website
	for _, w := range dbWebsites {
		websites = append(websites, domain.Website{
			ID:   int(w.ID),
			Name: w.Name,
		})
	}
	return websites, nil
}
