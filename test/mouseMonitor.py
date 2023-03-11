from evdev import InputDevice, ecodes

# 打开鼠标设备
dev = InputDevice('/dev/input/eventX')

# 循环监听事件
for event in dev.read_loop():
    # 如果事件是鼠标事件
    if event.type == ecodes.EV_KEY:
        # 如果按下左键
        if event.code == ecodes.BTN_LEFT and event.value == 1:
            print('Left button pressed')
        # 如果按下右键
        elif event.code == ecodes.BTN_RIGHT and event.value == 1:
            print('Right button pressed')
    # 如果事件是鼠标移动事件
    elif event.type == ecodes.EV_REL:
        # 如果是水平移动
        if event.code == ecodes.REL_X:
            print('Horizontal movement:', event.value)
        # 如果是垂直移动
        elif event.code == ecodes.REL_Y:
            print('Vertical movement:', event.value)
