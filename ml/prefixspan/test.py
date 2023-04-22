from numpy import number
from prefixspan import PrefixSpan
import os
import pandas as pd
import json


################## 全局变量
# 数据路径
dir_path = "D:\\hadoop_data\\event_data\\2"
# 区分事件数据的index
event_indexs = [0, 3, 4,6,7,8]
# prefix db 最大长度
max_seq_length = 200 

################## 函数
def prefixDemo(db, fre):
    ps = PrefixSpan(db)
    print(ps.frequent(fre))
    # print(ps.frequent(2, closed=True))
    # print(ps.topk(5, closed=True))
    # print(ps.frequent(2, generator=True))
    # print(ps.topk(5, generator=True))

# 读取行为原始数据 [][][]
def load_data():
    paths = os.listdir(dir_path)
    print(paths)
    data = []
    for idx in range(len(paths)):
        path = paths[idx]
        strs = path.split("_")
        if len(strs) != 2:
            return
        df = pd.read_csv(dir_path + "\\" + path)
        df = df.fillna("")
        data.append(df.values)
    return data

# 将单条记录转换成可区分event
def get_custom_event(event):
    custom_event = ""
    for idx in range(len(event_indexs)):
        event_idx = event_indexs[idx]
        custom_event = custom_event +"|"+ str(event[event_idx])
    return custom_event

# 区分事件，并对其编号
def number_event(data):
    # custom_event -> number 
    res = {}
    number = 1
    # 记录
    for i in range(len(data)):
        events = data[i]
        # 事件
        for idx in range(len(events)):
            event = events[idx]
            custom_event = get_custom_event(event)
            if custom_event in res.keys():
                continue
            res[custom_event] = number
            number = number + 1    
    return res 

# 将事件数据转换成编号序列 [][]
def get_number_data(event_data, event2number):
    res = []
    for i in range(len(event_data)):
        events = event_data[i]
        single_res = []
        # 事件
        for idx in range(len(events)):
            event = events[idx]
            custom_event = get_custom_event(event)
            if custom_event in event2number.keys():
                number = event2number[custom_event]
                single_res.append(number)
        if len(single_res) > 0:
            res.append(single_res)
    return res 

# 将数组按照size划分
def arr_size(arr,size):
    s=[]
    for i in range(0,int(len(arr))+1,size):
        if i+size<=len(arr):
            c=arr[i:i+size]
            s.append(c)
        else:
            c=arr[i:len(arr)]
            s.append(c)
    newlist = [x for x in s if x]
    return newlist

def split_number(numbers):
    res = []
    for idx in range(len(numbers)):
        arrs = arr_size(numbers[idx], max_seq_length)
        for jdx in range(len(arrs)):
            res.append(arrs[jdx])
    return res


event_data = load_data()
# print(data)
event2number = number_event(event_data)
# print(event2number)
json_str = json.dumps(event2number, indent=4)
with open('./prefixspan/event2number.json', 'w') as json_file:
    json_file.write(json_str)
numbers = get_number_data(event_data, event2number)
#print(numbers)

db = split_number(numbers)
print(db)
print(len(db))
#prefixDemo(db, int(len(db) * 0.1))