from sklearn.cluster import KMeans
from sklearn.metrics import silhouette_score
from sklearn.decomposition import PCA
from sklearn.preprocessing import StandardScaler
import pandas as pd
import matplotlib.pyplot as plt
from matplotlib import colors as mcolors

def k_means(k, data_pca):
    # 创建 KMeans 模型，并将数据聚类为k组
    kmeans = KMeans(n_clusters=k, random_state=0).fit(data_pca)

    # 获取聚类结果
    labels = kmeans.labels_

    # 获取聚类中心
    cluster_centers = kmeans.cluster_centers_

    # print("聚类结果：", labels)
    # print("聚类中心：", cluster_centers)

    score = silhouette_score(data_pca, labels)
    # print(f"Silhouette Coefficient: {score}")
    # 可视化结果
    n_clusters = len(set(labels)) - (1 if -1 in labels else 0)
    colors = list(mcolors.TABLEAU_COLORS.values())

    for i in range(n_clusters):
        plt.scatter(data_pca[labels==i, 0], data_pca[labels==i, 1], c=colors[i])

    plt.legend()
    plt.show()
    return score


def run():
    # 读取数据
    df = pd.read_csv('./cluster/test_data.csv')
    data = df.values
    # print('data')
    # print(data)

    scaler = StandardScaler()
    data_scaled = scaler.fit_transform(data)
    # print('data_scaled')
    # print(data_scaled)

    pca = PCA(n_components=2)
    data_pca = pca.fit_transform(data_scaled)

    for k in range(2, 10):
        score = k_means(k, data_pca)
        print('k=', k, 'score=', score)

run()