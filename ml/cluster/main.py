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


app = Flask(__name__)

# k-means模型
kmeans = joblib.load('./cluster/model/kmeans_model.pkl')
kmeans_scaler = joblib.load('./cluster/model/kmeans_scaler.pkl')
kmeans_pca = joblib.load('./cluster/model/kmeans_pca.pkl')

@app.route('/kmeans_cluster', methods=['POST'])
def kmeans_cluster():
    #提取请求数据
    form = request.form
    print("form=", form)
    # 转换为numpy数组
    reqStr = form['behavior_duration']
    print("reqStr=",reqStr)
    reqMap = json.loads(reqStr)
    param = [0, 0, 0, 0, 0]
    for key, value in reqMap.items():
        idx = int(key)-67
        if idx < 0 or idx >= 5:
            continue
        print("key=",key, "value=",value)
        param[idx] = value
    print("param=",param)
    data = np.array(param).reshape(1, -1)
    print(data)

    data_scaled = kmeans_scaler.transform(data)
    print('data_scaled')
    print(data_scaled)

    data_pca = kmeans_pca.transform(data_scaled)
    print('data_pca')
    print(data_pca)


    # 使用k-means模型进行预测
    features = data_pca[0].reshape(1,-1)
    print("features=", features)
    cluster = kmeans.predict(features)
    print("cluster=", cluster)
    # 返回预测结果
    return jsonify({'cluster': str(cluster[0])})

if __name__ == '__main__':
    app.run(port=5000)