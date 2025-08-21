package sqlite3

import (
	"database/sql"

	_ "modernc.org/sqlite" // ensure driver package is initialized
	drv "modernc.org/sqlite"
)

func init() {
	// Register the modernc driver under the name "sqlite3" so code that
	// expects to open driver "sqlite3" (eg gorm.io/driver/sqlite importing
	// mattn/go-sqlite3) will work.
	// modernc.org/sqlite exposes type Driver, use its zero value.
	sql.Register("sqlite3", &drv.Driver{})
}
