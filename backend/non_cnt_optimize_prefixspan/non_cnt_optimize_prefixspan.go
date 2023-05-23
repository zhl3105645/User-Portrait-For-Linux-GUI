package non_cnt_optimize_prefixspan

type Sequence []int

// copy makes a shallow copy of the sequence.
func (seq Sequence) copy() Sequence {
	copied := make(Sequence, len(seq))
	copy(copied, seq)
	return copied
}

// sequencePostfix 寻找前缀item所在位置，并返回后缀以及是否存在前缀
func (seq Sequence) sequencePostfix(item int) (Sequence, bool) {
	// 寻找序列前缀item所在位置，并返回对应的后缀
	for i, it := range seq {
		if it == item {
			if i == len(seq)-1 {
				return nil, true
			}

			return seq[i+1:].copy(), true
		}
	}
	return nil, false
}

// frequentItems returns frequent sequential patterns in db that have a length equal to one.
func frequentItems(db []Sequence, minSupport int) []int {
	var list []int // 满足最小支持度的item
	m := make(map[int]int)
	for _, seq := range db {
		exist := make(map[int]bool)
		for _, item := range seq {
			exist[item] = true
		}
		for item, _ := range exist {
			m[item] = m[item] + 1
		}
	}

	for item, cnt := range m {
		if cnt >= minSupport {
			list = append(list, item)
		}
	}

	return list
}

// appendToSequence 根据前缀item构建新的数据集, 返回长度大于0的数据集和前缀的支持度
func appendToSequence(db []Sequence, minSupport int, item int) ([]Sequence, int) {
	var projected []Sequence
	support := 0
	for _, seq := range db {
		// 寻找对应的后缀，并返回
		suffix, exist := seq.sequencePostfix(item)
		if exist {
			support++
		}
		if suffix == nil {
			continue
		}

		projected = append(projected, suffix)
	}
	return projected, support
}

func prefixSpan(db []Sequence, minSupport int, pattern Sequence, patternSupport int) []Sequence {
	items := frequentItems(db, minSupport) // 寻找当前db中长度为1的频繁序列

	var patterns []Sequence // 序列结果
	if len(pattern) > 0 {
		patterns = append(patterns, pattern) // 当前序列
	}

	for _, item := range items {
		// Append item to pattern to form a sequential pattern.
		// b 形式扩展模式一个项集
		projected, support := appendToSequence(db, minSupport, item) // 构建新的投影数据库
		if support >= minSupport {                                   // 大小满足最小支持度
			superPattern := append(pattern, item) // 更新模式
			newPatterns := prefixSpan(projected, minSupport, superPattern, support)
			patterns = append(patterns, newPatterns...)
		}
	}

	return patterns
}

func PrefixSpan(db []Sequence, minSupport int) []Sequence {
	return prefixSpan(db, minSupport, nil, 0)
}
