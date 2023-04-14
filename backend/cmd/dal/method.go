package dal

import "gorm.io/gen"

type QueryAll interface {
	// SELECT * FROM @@table
	FindAll() ([]gen.T, error)
}
