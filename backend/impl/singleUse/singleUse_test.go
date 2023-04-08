package singleUse

import (
	"backend/impl/rule"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"testing"
)

func TestName(t *testing.T) {
	a := make([]*RuleData, 100)
	for i := 0; i < 100; i++ {
		a[i] = &RuleData{
			ID:   int64(rand.Intn(10)),
			Time: fmt.Sprintf("time%d", i),
		}
	}
	b := make([]*rule.BehaviorRule, 5)
	for i := 0; i < 5; i++ {
		s1 := make([]string, 3)
		for j := 0; j < 3; j++ {
			s1[j] = strconv.FormatInt(int64(rand.Intn(10)), 10)
		}

		s2 := make([]string, 3)
		for j := 0; j < 3; j++ {
			s2[j] = strconv.FormatInt(int64(rand.Intn(10)), 10)
		}

		b[i] = &rule.BehaviorRule{
			Id: i,
			Behaviors: []string{
				fmt.Sprintf("(%s)", strings.Join(s1, ",")),
				fmt.Sprintf("(%s)", strings.Join(s2, ",")),
			},
		}
	}

	fmt.Println("a....")

	for _, data := range a {
		fmt.Printf("(%d,%s) ", data.ID, data.Time)
	}
	fmt.Printf("\n")

	fmt.Println("b....")
	for _, data := range b {
		fmt.Printf("(%d,%s) ", data.Id, strings.Join(data.Behaviors, "&&"))
	}
	fmt.Printf("\n")

	res := getBehaviorRuleIDs(a, b)
	fmt.Printf("res...")
	for _, data := range res {
		fmt.Printf("(%d,%s) ", data.ID, data.Time)
	}
}
