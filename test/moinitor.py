import threading
import time
from pynput import mouse, keyboard
import subprocess

def mouse_listener():
    def on_move(x, y):
        print(f'Mouse moved to ({x}, {y})')
        
    def on_click(x, y, button, pressed):
        if pressed:
            print(f'Mouse clicked at ({x}, {y}) with {button}')
        
    def on_scroll(x, y, dx, dy):
        print(f'Mouse scrolled at ({x}, {y}) ({dx}, {dy})')
    
    with mouse.Listener(on_move=on_move, on_click=on_click, on_scroll=on_scroll) as listener:
        listener.join()

def keyboard_listener():
    def on_press(key):
        try:
            print(f'Key {key.char} pressed')
        except AttributeError:
            print(f'Key {key} pressed')
        
    def on_release(key):
        try:
            print(f'Key {key.char} released')
        except AttributeError:
            print(f'Key {key} released')
        
    with keyboard.Listener(on_press=on_press, on_release=on_release) as listener:
        listener.join()

def window_listener():
    while True:
        window = subprocess.run(['xprop', '-root', '_NET_ACTIVE_WINDOW'], capture_output=True)
        window_id = window.stdout.split()[-1]
        window_info = subprocess.run(['xwininfo', '-id', window_id], capture_output=True)
        print("win_info", window_info.stdout.decode())
        time.sleep(1)

t1 = threading.Thread(target=mouse_listener)
t2 = threading.Thread(target=keyboard_listener)
t3 = threading.Thread(target=window_listener)

t1.start()
t2.start()
t3.start()

t1.join()
t2.join()
t3.join()
