// internal/scheduler/cron.go
package scheduler

import (
	"log"
	"os"
	"time"

	"context"
	"scanfetcher/internal/service"

	"github.com/robfig/cron/v3"
)

type CronScheduler struct {
	c   *cron.Cron
	svc *service.ScanService
	// optionnel: logger, metrics…
}

func NewCronScheduler(svc *service.ScanService) *CronScheduler {
	return &CronScheduler{
		c:   cron.New(),
		svc: svc,
	}
}

func (s *CronScheduler) Start() error {
	spec := os.Getenv("CRON_JOB_SCHEDULE")
	if spec == "" {
		spec = "@every 10m"
	}
	_, err := s.c.AddFunc(spec, func() {
		log.Println("[cron] tick at", time.Now().Format(time.RFC3339))
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
		defer cancel()

		if _, err := s.svc.RunAll(ctx); err != nil {
			log.Printf("[cron] runAll error: %v", err)
		}
	})
	if err != nil {
		return err
	}
	s.c.Start()
	log.Printf("[cron] started with schedule: %s", spec)
	return nil
}

func (s *CronScheduler) Stop() {
	s.c.Stop()
}
