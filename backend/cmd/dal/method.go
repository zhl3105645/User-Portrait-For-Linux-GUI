package dal

import "gorm.io/gen"

type QueryAll interface {
	// SELECT * FROM @@table
	FindAll() ([]gen.T, error)
}

type MethodForApp interface {
	// Insert @@table (app_name) value (@name)
	InsertOne(name string) ([]gen.T, error)
}
