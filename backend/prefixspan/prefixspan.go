package prefixspan

// the PrefixSpan algorithm, as defined in http://hanj.cs.illinois.edu/pdf/span01.pdf.

// "_" in the paper.
const placeholder = -1

func itoa(item int) string {
	if item == placeholder {
		return "_"
	}
	r := rune(int('a') + item)
	return string(r)
}

// An ItemSet is a subset of items.
type ItemSet []int

// String implements Stringer.
func (set ItemSet) String() string {
	if len(set) == 1 {
		return itoa(set[0])
	}

	s := "("
	for _, it := range set {
		s += itoa(it)
	}
	s += ")"
	return s
}

// contains checks if an item is in an itemset.
func (set ItemSet) contains(item int) bool {
	for _, it := range set {
		if it == item {
			return true
		}
	}
	return false
}

// A sequence is an ordered list of itemsets.
type Sequence []ItemSet

// String implements Stringer.
func (seq Sequence) String() string {
	s := "<"
	for _, set := range seq {
		s += set.String()
	}
	s += ">"
	return s
}

// copy makes a shallow copy of the sequence.
func (seq Sequence) copy() Sequence {
	copied := make(Sequence, len(seq))
	copy(copied, seq)
	return copied
}

// suffix returns the suffix of seq w.r.t. the prefix of seq until indices i, j.
// 返回序列后缀，从第i+1个项集开始保留，第i+1个去除前j个项
func (seq Sequence) suffix(i, j int) Sequence {
	suffix := seq[i:]
	if len(suffix) == 0 {
		return nil
	}
	itemSetTrail := suffix[0][j+1:]
	if len(itemSetTrail) == 0 {
		suffix = suffix[1:]
	} else {
		suffix = suffix.copy()
		suffix[0] = append(ItemSet{placeholder}, itemSetTrail...)
	}
	return suffix
}

// sequencePostfix returns the suffix of seq w.r.t the itemset containing only
// item (noted <item>). It returns nil if and only if such suffix doesn't exist.
func (seq Sequence) sequencePostfix(item int) Sequence {
	// 寻找序列前缀item所在位置，并返回对应的后缀
	for i, set := range seq {
		for j, it := range set {
			if it == item {
				// 第一个项集是 _a 形式
				if i == 0 && j > 0 && set[0] == placeholder {
					continue
				}
				return seq.suffix(i, j)
			}
		}
	}
	return nil
}

// itemSetPostfix returns the suffix of seq w.r.t. itemSet, assuming _ represents
// itemSet without its last item.
func (seq Sequence) itemSetPostfix(itemSet ItemSet) Sequence {
	lastItem := itemSet[len(itemSet)-1] // 被添加的项
	for i, set := range seq {
		// 判断 _a 形式
		if set[0] == placeholder && set[1] == lastItem {
			return seq.suffix(i, 1)
		}

		// abcd bc 形式
		delta := 0
		for j, it := range set {
			if it == itemSet[delta] {
				delta++
				if delta == len(itemSet) {
					return seq.suffix(i, j)
				}
			} else {
				delta = 0
			}
		}
	}
	return nil
}

// frequentItems returns frequent sequential patterns in db that have a length equal to one.
func frequentItems(db []Sequence, minSupport int) ItemSet {
	var list ItemSet // 满足最小支持度的项集
	m := make(map[int]int)
	for _, seq := range db {
		exist := make(map[int]bool)
		for _, itemSet := range seq {
			for _, item := range itemSet {
				if item == placeholder {
					continue
				}

				exist[item] = true
			}
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

// appendToSequence takes an alpha-projected database db, and returns an
// alpha'-projected database, given alpha' is alpha appended with <item>. The
// length of the result is the frequency of the pattern alpha'.
func appendToSequence(db []Sequence, minSupport int, item int) []Sequence {
	var projected []Sequence
	for _, seq := range db {
		// 寻找对应的后缀，并返回
		suffix := seq.sequencePostfix(item)
		if suffix == nil {
			continue
		}

		projected = append(projected, suffix)
	}
	return projected
}

// appendToItemSet takes an alpha-projected database db, and returns an
// alpha'-projected database, given alpha' is alpha with item added to its last
// itemset. The length of the result is the frequency of the pattern alpha'.
func appendToItemSet(db []Sequence, minSupport int, itemSet ItemSet) []Sequence {
	var projected []Sequence
	for _, seq := range db {
		suffix := seq.itemSetPostfix(itemSet)
		if suffix == nil {
			continue
		}

		projected = append(projected, suffix)
	}
	return projected
}

// prefixSpan minSupport 最小支持度 pattern 当前前缀序列 返回满足条件的序列集
func prefixSpan(db []Sequence, minSupport int, pattern Sequence) []Sequence {
	// Scan db once, find the set of frequent items
	items := frequentItems(db, minSupport) // 寻找当前db中长度为1的频繁序列

	var patterns []Sequence // 结果
	if len(pattern) > 0 {
		patterns = append(patterns, pattern) // 当前序列
	}

	for _, item := range items {
		// Append item to the last element of pattern to form a sequential pattern
		// _b 形式添加到模式最后一个项集
		if len(pattern) > 0 {
			lastItemSet := pattern[len(pattern)-1] // 最后一个项集
			if !lastItemSet.contains(item) {       // 不包含 _b, 添加
				lastItemSet = append(lastItemSet, item)
				projected := appendToItemSet(db, minSupport, lastItemSet)
				if len(projected) >= minSupport {
					superPattern := pattern.copy()
					superPattern[len(superPattern)-1] = lastItemSet // 更新模式
					patterns = append(patterns, prefixSpan(projected, minSupport, superPattern)...)
				}
			}
		}

		// Append item to pattern to form a sequential pattern.
		// b 形式扩展模式一个项集
		projected := appendToSequence(db, minSupport, item) // 构建新的投影数据库
		if len(projected) >= minSupport {                   // 大小满足最小支持度
			superPattern := append(pattern, ItemSet{item}) // 更新模式
			patterns = append(patterns, prefixSpan(projected, minSupport, superPattern)...)
		}
	}

	return patterns
}

// PrefixSpan searches frequent sequences in the database db. Sequences must
// only contain positive values.
func PrefixSpan(db []Sequence, minSupport int) []Sequence {
	return prefixSpan(db, minSupport, nil)
}
