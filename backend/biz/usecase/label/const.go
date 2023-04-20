package label

type Operator string

const (
	Unknown Operator = ""
	GE      Operator = ">="
	GT      Operator = ">"
	EQ      Operator = "="
	LE      Operator = "<"
	LT      Operator = "<="
)

func getOperator(s string) Operator {
	switch Operator(s) {
	case GE:
		return GE
	case GT:
		return GT
	case EQ:
		return EQ
	case LE:
		return LE
	case LT:
		return LT
	}
	return Unknown
}

type ConvertRule struct {
	Op     Operator
	XValue float64
	YValue int64
}

func (c *ConvertRule) Match(value float64) bool {
	switch c.Op {
	case GE:
		return value >= c.XValue
	case GT:
		return value > c.XValue
	case EQ:
		return value == c.XValue
	case LE:
		return value < c.XValue
	case LT:
		return value <= c.XValue
	}

	return false
}
