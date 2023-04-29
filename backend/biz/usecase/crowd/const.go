package crowd

import (
	"strconv"
)

const (
	UnionAnd = 1
	UnionOr  = 2
)

const (
	DivideGT = 1
	DivideGE = 2
	DivideEQ = 3
	DivideLS = 4
	DivideLE = 5
)

func Union(ope int64, arr1 []int64, arr2 []int64) []int64 {
	res := make([]int64, 0)

	map1 := make(map[int64]bool)
	map2 := make(map[int64]bool)
	mapAll := make(map[int64]bool)
	for _, userId := range arr1 {
		map1[userId] = true
		mapAll[userId] = true
	}
	for _, userId := range arr2 {
		map2[userId] = true
		mapAll[userId] = true
	}
	if ope == UnionAnd {
		for userId, _ := range map1 {
			if _, ok := map2[userId]; ok {
				res = append(res, userId)
			}
		}
	} else if ope == UnionOr {
		for userId, _ := range mapAll {
			res = append(res, userId)
		}
	}

	return res
}

func MatchDivide(ope int64, userVal string, compareVal string) bool {
	if ope == DivideEQ {
		return userVal == compareVal
	}

	userIntVal, err1 := strconv.ParseInt(userVal, 10, 64)
	compareIntVal, err2 := strconv.ParseInt(compareVal, 10, 64)
	if err1 != nil || err2 != nil {
		return false
	}

	switch ope {
	case DivideGT:
		return userIntVal > compareIntVal
	case DivideGE:
		return userIntVal >= compareIntVal
	case DivideLS:
		return userIntVal < compareIntVal
	case DivideLE:
		return userIntVal <= compareIntVal
	}

	return false
}
