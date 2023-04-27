import os
from numpy import longlong
import numpy as np
import pandas as pd

source_dir = 'C:\\Users\\鸿泽量天\\Desktop\\data_RedPandaIDE'
target_dir = 'C:\\Users\\鸿泽量天\\Desktop\\utf-data'

for filename in os.listdir(source_dir):
    if filename.endswith('.csv'):
        file_path = os.path.join(source_dir, filename)
        df = pd.read_csv(file_path, encoding='gbk')
        # 类型转换为 整型
        df['事件类型'] = df['事件类型'].fillna(0)
        df['事件时间'] = df['事件时间'].fillna(0)
        df['鼠标点击类型'] = df['鼠标点击类型'].fillna(0)
        df['鼠标点击按键'] = df['鼠标点击按键'].fillna(0)
        df['鼠标移动类型'] = df['鼠标移动类型'].fillna(0)
        df['键盘点击类型'] = df['键盘点击类型'].fillna(0)
        df['组件类型'] = df['组件类型'].fillna(0)

        df["事件类型"] = df["事件类型"].astype(longlong)
        df["事件时间"] = df["事件时间"].astype(longlong)
        df["鼠标点击类型"] = df["鼠标点击类型"].astype(longlong)
        df["鼠标点击按键"] = df["鼠标点击按键"].astype(longlong)
        df["鼠标移动类型"] = df["鼠标移动类型"].astype(longlong)
        df["键盘点击类型"] = df["键盘点击类型"].astype(longlong)
        df["组件类型"] = df["组件类型"].astype(longlong)

        target_path = os.path.join(target_dir, filename)
        df.to_csv(target_path, index=False, encoding='utf-8')