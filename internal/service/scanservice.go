package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"scanfetcher/internal/domain"
	"scanfetcher/internal/notifier"
	"scanfetcher/internal/repository"
	"scanfetcher/internal/scraper"

	"github.com/go-rod/rod"
)

type ScanService struct {
	websites *repository.WebsiteRepository
	scans    *repository.ScanRepository
	scrapers scraper.Factory
	notify   notifier.Notifier
	chatID   int64
	// config
	perSiteTimeout time.Duration
}

func NewScanService(
	websites *repository.WebsiteRepository,
	scans *repository.ScanRepository,
	scrapers scraper.Factory,
	notify notifier.Notifier,
	chatID int64,
) *ScanService {
	return &ScanService{
		websites:       websites,
		scans:          scans,
		scrapers:       scrapers,
		notify:         notify,
		chatID:         chatID,
		perSiteTimeout: 2 * time.Minute,
	}
}

// RunAll: scanne tous les sites et retourne les nouveaux scans trouvés
func (s *ScanService) RunAll(ctx context.Context) ([]domain.Scan, error) {
	websites, err := s.websites.GetAll()
	if err != nil {
		return nil, fmt.Errorf("get websites: %w", err)
	}
	var allNew []domain.Scan

	for _, site := range websites {
		ctxSite, cancel := context.WithTimeout(ctx, s.perSiteTimeout)
		newScans, err := s.runForSite(ctxSite, site)
		cancel()

		if err != nil {
			log.Printf("[scan] site=%s err=%v", site.Name, err)
			continue
		}
		allNew = append(allNew, newScans...)
	}
	return allNew, nil
}

// RunForWebsiteID: exposé au handler HTTP
func (s *ScanService) RunForWebsiteID(ctx context.Context, websiteID uint) ([]domain.Scan, error) {
	site, err := s.websites.GetByID(websiteID)
	if err != nil {
		return nil, fmt.Errorf("get website %d: %w", websiteID, err)
	}
	return s.runForSite(ctx, *site)
}

func (s *ScanService) runForSite(ctx context.Context, site domain.Website) ([]domain.Scan, error) {
	// récupérer les configurations de scan pour ce site
	scanEntries, err := s.scans.GetByWebsiteID(site.ID)
	if err != nil {
		return nil, fmt.Errorf("get scans for site %d: %w", site.ID, err)
	}
	if len(scanEntries) == 0 {
		return nil, nil
	}

	// un seul browser par site (plus efficace que par scan)
	browser := rod.New().MustConnect()
	defer func() { _ = browser.Close() }()

	scr, err := s.scrapers.Get(site.Name)
	if err != nil {
		return nil, fmt.Errorf("scraper not found for %s: %w", site.Name, err)
	}

	var newScans []domain.Scan
	for _, entry := range scanEntries {
		select {
		case <-ctx.Done():
			return newScans, ctx.Err()
		default:
		}

		scans, err := scr.Scrape(browser, entry.Url)
		if err != nil {
			log.Printf("[scan] site=%s url=%s scrape err=%v", site.Name, entry.Url, err)
			continue
		}
		if len(scans) == 0 || scans[0].LastScanRead == entry.LastScanRead {
			log.Printf("[scan] site=%s no new results", site.Name)
			continue
		}

		if err := s.scans.UpdateLastScanRead(entry.ID, scans[0].LastScanRead); err != nil {
			log.Printf("[scan] site=%s update last read err=%v", site.Name, err)
			continue
		}

		// notification + agrégation
		msg := fmt.Sprintf(
			"Nouveau scan: %s\nURL: %s\nDernière lecture: %s\nSiteID: %d",
			scans[0].Name, scans[0].Url, scans[0].LastScanRead, entry.WebsiteID,
		)
		if s.notify != nil {
			if err := s.notify.SendMessage(s.chatID, msg); err != nil {
				log.Printf("[notify] site=%s err=%v", site.Name, err)
			}
		}
		newScans = append(newScans, scans...)
	}
	return newScans, nil
}
