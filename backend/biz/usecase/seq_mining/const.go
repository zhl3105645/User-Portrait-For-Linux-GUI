package seq_mining

const (
	StatusBegin = 1 // 开始执行
	StatusRun   = 2 // 执行中
	StatusEnd   = 3 // 执行完成
)

type Result struct {
	Cnt     int
	Numbers []int
}

type EventNumber struct {
	Number int
	Event  string
}
