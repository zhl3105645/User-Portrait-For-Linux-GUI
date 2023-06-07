from pandas.core import api
from sklearn.cluster import KMeans
from sklearn.metrics import silhouette_score
from sklearn.decomposition import PCA
from sklearn.preprocessing import StandardScaler
import pandas as pd
import matplotlib.pyplot as plt
from matplotlib import colors as mcolors
import joblib
import numpy as np
from pylab import mpl
mpl.rcParams['font.sans-serif'] = ['FangSong']
mpl.rcParams['axes.unicode_minus'] = False

def k_means(k, data_pca, scaler, pca):
    # 创建 KMeans 模型，并将数据聚类为k组
    kmeans = KMeans(n_clusters=k, random_state=0).fit(data_pca)

    # 获取聚类结果
    labels = kmeans.labels_

    # 获取聚类中心
    cluster_centers = kmeans.cluster_centers_

    # print("聚类结果：", labels)
    # print("聚类中心：", cluster_centers)

    # 将聚类结果从PCA空间转换回原始数据空间（但仍然是标准化的）
    cluster_centers_scaled = pca.inverse_transform(cluster_centers)
    # 将标准化的聚类结果转换回原始数据空间
    cluster_centers_original = scaler.inverse_transform(cluster_centers_scaled)
    #print("cluster_centers_scaled=",cluster_centers_scaled)
    print("cluster_centers_original=",cluster_centers_original)
    draw_center(cluster_centers_original)

    score = silhouette_score(data_pca, labels)
    # print(f"Silhouette Coefficient: {score}")
    # 可视化结果
    n_clusters = len(set(labels)) - (1 if -1 in labels else 0)
    colors = list(mcolors.TABLEAU_COLORS.values())

    for i in range(n_clusters):
        plt.scatter(data_pca[labels==i, 0], data_pca[labels==i, 1], c=colors[i])
    plt.scatter(cluster_centers[:,0], cluster_centers[:,1], marker='x', s=200, linewidths=3,
            color='r', zorder=10)

    plt.legend()
    plt.show()
    return score, labels

def save_model(k, data_pca, scaler, pca):
    kmeans = KMeans(n_clusters=k, random_state=0).fit(data_pca)
    
    labels = kmeans.labels_
    # print("labels=",labels)
    # 保存模型
 
    joblib.dump(kmeans, './cluster/model/kmeans_model.pkl')
    joblib.dump(scaler, './cluster/model/kmeans_scaler.pkl')
    joblib.dump(pca, './cluster/model/kmeans_pca.pkl') 

def draw_center(data):
    plt.ylabel('行为时长(ms)')
    # 设置柱状图的宽度
    bar_width = 0.1

    # 设置柱状图的位置
    x1 = np.arange(len(data))
    x2 = [x + bar_width for x in x1]
    x3 = [x + bar_width for x in x2]
    x4 = [x + bar_width for x in x3]
    x5 = [x + bar_width for x in x4]

    # 绘制柱状图
    plt.bar(x1, [row[0] for row in data], width=bar_width, label='编码')
    plt.bar(x2, [row[1] for row in data], width=bar_width, label='测试')
    plt.bar(x3, [row[2] for row in data], width=bar_width, label='调试')
    plt.bar(x4, [row[3] for row in data], width=bar_width, label='浏览')
    plt.bar(x5, [row[4] for row in data], width=bar_width, label='使用工具')

    # 设置x轴刻度
    plt.xticks([x + 2 * bar_width for x in x1], ['聚类1', '聚类2', '聚类3'])

    # 添加图例
    plt.legend()

    # 显示图形
    plt.show()
    return

def draw_score(k_nums, scores):
    plt.xlabel('参数k')
    plt.ylabel('轮廓系数')
    plt.plot(k_nums,scores)
    plt.show()
    return 


def run():
    # 读取数据
    df = pd.read_csv('./cluster/data/test_data2.csv')
    data = df.values
    #print('data = ')
    #print(data)
    userIds = []
    for idx in range(len(data)):
        userIds.append(data[idx][0])
    #print("user_id = " , userIds)

    df = df.drop('user_id', axis=1)
    data = df.values
    #print(data)

    scaler = StandardScaler()
    data_scaled = scaler.fit_transform(data)
    #print('data_scaled = ')
    #print(data_scaled)

    pca = PCA(n_components=2)
    data_pca = pca.fit_transform(data_scaled)
    #print('data_pca = ')
    #print(data_pca)

    k_nums = []
    scores = []

    for k in range(3, 4):
        score,labels = k_means(k, data_pca, scaler,pca)
        print('k=', k, 'score=', score)
        k_nums.append(k)
        scores.append(score)
        # 使用column_stack函数将两个数组组合成一个二维数组
        res = np.column_stack((userIds, labels))
        # 使用savetxt函数将二维数组保存到csv文件中
        np.savetxt('./cluster/data/result.csv', res, delimiter=',',fmt='%d')

    #draw_score(k_nums, scores)
run()


