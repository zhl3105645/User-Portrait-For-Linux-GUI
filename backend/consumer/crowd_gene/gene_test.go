package crowd_gene

import (
	"backend/cmd/dal"
	"testing"
)

func TestGene(t *testing.T) {
	dal.Init()

	Gene(2, 1)
}
