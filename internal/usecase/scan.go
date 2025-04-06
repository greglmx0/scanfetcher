package usecase

import "scanfetcher/internal/domain"

type ScanRepository interface {
	Create(domain.Scan) error
	GetByWebsiteID(int) ([]domain.Scan, error)
}

type ScanUseCase struct {
	repo ScanRepository
}

func NewScanUseCase(repo ScanRepository) *ScanUseCase {
	return &ScanUseCase{repo: repo}
}

func (uc *ScanUseCase) CreateScan(s domain.Scan) error {
	return uc.repo.Create(s)
}

func (uc *ScanUseCase) GetScansForWebsite(websiteID int) ([]domain.Scan, error) {
	return uc.repo.GetByWebsiteID(websiteID)
}
