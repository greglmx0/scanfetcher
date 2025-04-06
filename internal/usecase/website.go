package usecase

import (
	"scanfetcher/internal/domain"
)

type WebsiteRepository interface {
	Create(domain.Website) error
	GetByID(id int) (*domain.Website, error)
	GetAll() ([]domain.Website, error)
}

type WebsiteUseCase struct {
	repo WebsiteRepository
}

func NewWebsiteUseCase(repo WebsiteRepository) *WebsiteUseCase {
	return &WebsiteUseCase{repo: repo}
}

func (uc *WebsiteUseCase) CreateWebsite(name string) error {
	website := domain.Website{Name: name}
	return uc.repo.Create(website)
}

func (uc *WebsiteUseCase) GetWebsite(id int) (*domain.Website, error) {
	return uc.repo.GetByID(id)
}

func (uc *WebsiteUseCase) ListWebsites() ([]domain.Website, error) {
	return uc.repo.GetAll()
}
