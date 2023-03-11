from evdev import InputDevice, ecodes

# 打开鼠标设备
dev = InputDevice('/dev/input/event2')
# print(ecodes.EV_REL, ecodes.EV_KEY, ecodes.EV_ABS, ecodes.EV_SYN)
# print(ecodes.ABS_X, ecodes.ABS_Y)

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
    elif event.type == ecodes.EV_ABS:
        print(event.code, event.value)
        if event.code == ecodes.ABS_X:
            print("X: ", event.value)
        elif event.code == ecodes.ABS_Y:
            print("Y: ", event.value)
