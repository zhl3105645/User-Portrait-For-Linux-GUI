package event_data

// MaxNoOperateTimeS 允许最大不操作时间
const MaxNoOperateTimeS = 10

type EventType string

const (
	AppStart   EventType = "1"
	AppQuit    EventType = "2"
	MouseClick EventType = "3"
	MouseMove  EventType = "4"
	KeyClick   EventType = "5"
	MouseWheel EventType = "6"
	Shortcut   EventType = "7"
)

type MouseClickType string

const (
	One MouseClickType = "1"
	Two MouseClickType = "2"
)

type MouseClickButton string

const (
	LeftButton  MouseClickButton = "1"
	RightButton MouseClickButton = "2"
)

type MouseMoveType string

const (
	MoveBegin MouseMoveType = "1"
	MoveEnd   MouseMoveType = "2"
)

type KeyClickType string

const (
	Single    KeyClickType = "1"
	Component KeyClickType = "2"
)

type ComponentType string

const (
	None      ComponentType = "-1" // 未定义组件
	Button    ComponentType = "1"  // 按钮，QPushButton QToolButton QRadioButton QCheckBox
	Combo     ComponentType = "2"  // 下拉组合框   QComboBox QFontComboBox
	Text      ComponentType = "3"  // 文本编辑 QLineEdit QTextEdit QPlainTextEdit QLabel
	Spin      ComponentType = "4"  // 滚轮 QSpinBox QDoubleSpinBox
	Slider    ComponentType = "5"  // 滑块 QDial QSlider QScrollBar
	Calendar  ComponentType = "6"  // 日历
	Lcd       ComponentType = "7"  // lcd数字
	Progress  ComponentType = "8"  // 进度条
	List      ComponentType = "9"  // 列表视图
	Tree      ComponentType = "10" // 树状视图
	Table     ComponentType = "11" // 表视图
	Column    ComponentType = "12" // 列视图
	Action    ComponentType = "13" // 命令
	Container ComponentType = "14" // 容器
)

const (
	EventTypeIndex        = 0
	EventTimeIndex        = 1
	MousePos              = 2
	MouseClickTypeIndex   = 3
	MouseClickButtonIndex = 4
	MouseMoveTypeIndex    = 5
	KeyClickTypeIndex     = 6
	KeyCodeIndex          = 7
	ComponentNameIndex    = 8
	ComponentTypeIndex    = 9
	ComponentExtraIndex   = 10
)
