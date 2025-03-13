package utils

import (
	"fmt"
)

type DataSource struct {
	User     string
	Password string
	Host     string
	Database string
}

func BuildDatasourceName(ds DataSource) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable&timezone=UTC",
		ds.User,
		ds.Password,
		ds.Host,
		ds.Database,
	)
}
