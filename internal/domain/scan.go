package domain

type Scan struct {
	ID           int
	Name         string
	Url          string
	LastScanRead string
	WebsiteID    uint
}
