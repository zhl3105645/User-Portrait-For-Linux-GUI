package label_gene

import (
	"backend/biz/hadoop"
	"backend/cmd/dal"
	"context"
	"testing"
)

func TestGene(t *testing.T) {
	dal.Init()
	hadoop.Init(context.Background())

	Gene(2, Gender)
	Gene(2, Age)
	Gene(2, Career)
	Gene(2, UseTime)
	Gene(2, UsePeriod)
	Gene(2, UseActivity)
	Gene(2, ProgramLanguage)
	Gene(2, CodeSpeed)
	Gene(2, CodeAbility)
	Gene(2, ShortcutFre)
	Gene(2, GitFre)
	Gene(2, GitNorm)
	Gene(2, BehaviorPrefer)
}
