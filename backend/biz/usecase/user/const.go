package user

type Gender int

const (
	Man   = 1
	Woman = 2
)

var Gender2Desc = map[Gender]string{
	Man:   "男",
	Woman: "女",
}
