package optimize_prefixspan

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
)

var testMinSupport = 2

var testDB = []Sequence{
	{0, 1, 2, 3, 4},
	{1, 1, 1, 3, 4},
	{2, 1, 2, 2, 0},
	{1, 1, 1, 2, 2},
}

func sortDB(db []Sequence) {
	sort.Slice(db, func(i, j int) bool {
		a, b := db[i], db[j]

		for k := 0; k < len(a) && k < len(b); k++ {
			if a[k] == b[k] {
				continue
			}
			return a[k] < b[k]
		}
		return len(a) < len(b)
	})
}

func TestPrefixSpan(t *testing.T) {
	result, cnts := PrefixSpan(testDB, testMinSupport)

	for i, res := range result {
		fmt.Printf("res = %v, cnt = %d\n", res, cnts[i])
	}

	want := []Sequence{
		{0},
		{1},
		{1, 2},
		{1, 2, 2},
		{1, 3},
		{1, 3, 4},
		{1, 4},
		{1, 1},
		{1, 1, 1},
		{2},
		{2, 2},
		{3},
		{3, 4},
		{4},
	}

	sortDB(want)
	sortDB(result)

	if !reflect.DeepEqual(result, want) {
		t.Errorf("PrefixSpan() = \n%v\n, want \n%v", result, want)
	}
}
