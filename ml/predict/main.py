from flask import Flask, request, jsonify
import numpy as np
from sklearn.cluster import KMeans
from sklearn.metrics import silhouette_score
from sklearn.decomposition import PCA
from sklearn.preprocessing import StandardScaler
import pandas as pd
import numpy as np
import matplotlib.pyplot as plt
from matplotlib import colors as mcolors
import joblib
import json
import torch
import torch.nn as nn
from sklearn.model_selection import train_test_split
from torch.utils.data import TensorDataset, DataLoader
import matplotlib.pyplot as plt
from pylab import mpl

from test.lstm1 import LSTMModel, input_size, hidden_size, num_layers, num_classes


app = Flask(__name__)

# lstm 模型
model = LSTMModel(input_size=input_size, hidden_size=hidden_size, num_layers=num_layers, num_classes=num_classes)
model.load_state_dict(torch.load("./predict/model/model"))

@app.route('/lstm', methods=['POST'])
def kmeans_cluster():
    #提取请求数据
    form = request.form
    print("form=", form)
    # 转换为numpy数组
    req = form['event_rule_ids']
    arr = req.split(",")
    print(arr)
    arr = np.array(arr).astype(float)
    # 获取序列长度
    sequence_length = arr.shape[0]
    # 将一维数组转换为形状为(1, sequence_length, 1)的三维数组
    arr_3d = arr.reshape(1, sequence_length, 1)
    
    print(arr_3d)
    current_sequence = torch.tensor(arr_3d, dtype=torch.float32)
    print(current_sequence)
    
    next_event_prob = model(current_sequence)
    next_event_prob = nn.functional.softmax(next_event_prob, dim=1)
    next_event_prob = next_event_prob.detach().numpy().flatten()
    list_prop = next_event_prob.tolist()
    #print(list_prop)

    res = {}
    for idx in range(len(list_prop)):
        res[idx-2] = list_prop[idx]

    print(res)    
    res = json.dumps(res)
    
    # 返回预测结果
    return jsonify({'next_event_prob': res})

if __name__ == '__main__':
    app.run(port=5001)