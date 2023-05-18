package hadoop

// 事件类型
const (
	AppStart   = 1
	AppQuit    = 2
	MouseClick = 3
	MouseMove  = 4
	KeyClick   = 5
	MouseWheel = 6
	Shortcut   = 7
)

// 鼠标点击类型
const (
	One = 1
	Two = 2
)

// 鼠标点击按键
const (
	LeftButton  = 1
	RightButton = 2
)

// 鼠标移动类型
const (
	MoveBegin = 1
	MoveEnd   = 2
)

// 键盘点击类型
const (
	Single    = 1
	Component = 2
)
