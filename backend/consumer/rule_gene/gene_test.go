package rule_gene

import (
	"backend/biz/hadoop"
	"backend/cmd/dal"
	"context"
	"testing"
)

func TestGene(t *testing.T) {
	dal.Init()
	hadoop.Init(context.Background())

	Gene(2)
}
