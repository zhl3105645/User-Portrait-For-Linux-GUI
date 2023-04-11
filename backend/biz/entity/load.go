package entity

import "context"

type ILoad interface {
	Load(ctx context.Context) error
}
