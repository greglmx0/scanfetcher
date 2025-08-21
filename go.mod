module scanfetcher

go 1.24.0

require (
	github.com/go-rod/rod v0.116.2
	github.com/go-telegram-bot-api/telegram-bot-api/v5 v5.5.1
	github.com/gorilla/mux v1.8.1
	github.com/robfig/cron v1.2.0
	gorm.io/gorm v1.30.1
	modernc.org/sqlite v1.38.2
)

require github.com/mattn/go-sqlite3 v1.14.22 // indirect

// Use local shim to satisfy imports of github.com/mattn/go-sqlite3 without CGO
replace github.com/mattn/go-sqlite3 => ./third_party/mattn_go_sqlite3

// removed indirect CGO-based driver (github.com/mattn/go-sqlite3)
require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v0.1.9 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	github.com/ysmood/fetchup v0.2.3 // indirect
	github.com/ysmood/goob v0.4.0 // indirect
	github.com/ysmood/got v0.40.0 // indirect
	github.com/ysmood/gson v0.7.3 // indirect
	github.com/ysmood/leakless v0.9.0 // indirect
	golang.org/x/exp v0.0.0-20250819193227-8b4c13bb791b // indirect
	golang.org/x/sys v0.35.0 // indirect
	golang.org/x/text v0.28.0 // indirect
	gorm.io/driver/sqlite v1.6.0
	modernc.org/libc v1.66.7 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
)
