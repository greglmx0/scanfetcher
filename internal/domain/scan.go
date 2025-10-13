package domain

type Scan struct {
	ID           uint
	Name         string
	Url          string
	LastScanRead string
	WebsiteID    uint
}
