from evdev import InputDevice, categorize, ecodes

dev = InputDevice('/dev/input/event1')

# 循环监听事件
for event in dev.read_loop():
    if event.type == ecodes.EV_KEY:
        key_event = categorize(event)
        if key_event.keystate == key_event.key_down:
            print(key_event.keycode + "be pressed down")
        elif key_event.keystate == key_event.key_up:
            print(key_event.keycode + "be pressed up")