import subprocess

def get_window_info(app_name):
    # 使用 xdotool 命令获取当前激活窗口的 ID
    window_id = subprocess.check_output(["xdotool", "getactivewindow"], universal_newlines=True).strip()
    print("window_id=", window_id)

    # 使用 xprop 命令获取窗口的 WM_CLASS 属性
    output = subprocess.check_output(["xprop", "-id", window_id, "WM_CLASS"], universal_newlines=True)
    print("output=", output)
    class_name = output.split("=")[1].strip().split(",")[0].strip().strip('"')
    print("class_name=", class_name)

    if class_name.lower() != app_name.lower():
        return None

    # 使用 xwininfo 命令获取窗口信息
    output = subprocess.check_output(["xwininfo", "-id", window_id], universal_newlines=True)
    print(output)
    # 从输出中获取窗口大小和位置信息
    width = int([x.split(":")[1].strip() for x in output.split("\n") if "Width" in x][0])
    height = int([x.split(":")[1].strip() for x in output.split("\n") if "Height" in x][0])
    x = int([x.split(":")[1].strip() for x in output.split("\n") if "Absolute upper-left X" in x][0])
    y = int([x.split(":")[1].strip() for x in output.split("\n") if "Absolute upper-left Y" in x][0])

    return {"id": window_id, "x": x, "y": y, "width": width, "height": height}

if __name__ == "__main__":
    s = get_window_info("Code")
    print(s)