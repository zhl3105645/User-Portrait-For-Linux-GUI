package seq_mining

import (
	"backend/biz/hadoop"
	"backend/cmd/dal"
	"context"
	"testing"
)

func TestGene(t *testing.T) {
	dal.Init()
	hadoop.Init(context.Background())
	Gene(2, 2)
}
