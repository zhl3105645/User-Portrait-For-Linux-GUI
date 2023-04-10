import pandas as pd
from sklearn.cluster import DBSCAN
from sklearn.preprocessing import StandardScaler
from sklearn.decomposition import PCA
import matplotlib.pyplot as plt
from matplotlib import colors as mcolors
from sklearn.metrics import silhouette_score

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

dbscan = DBSCAN()
dbscan.fit(data_pca)

labels = dbscan.labels_

score = silhouette_score(data_pca, labels)
print(f"Silhouette Coefficient: {score}")

# 输出结果
n_clusters = len(set(labels)) - (1 if -1 in labels else 0)
print(f"Estimated number of clusters: {n_clusters}")

for i in range(n_clusters):
    print(f"Cluster {i+1}:")
    print(data[labels==i])

# 可视化结果
n_clusters = len(set(labels)) - (1 if -1 in labels else 0)
colors = list(mcolors.TABLEAU_COLORS.values())

for i in range(n_clusters):
    plt.scatter(data_pca[labels==i, 0], data_pca[labels==i, 1], c=colors[i], label=f'Cluster {i+1}')

plt.legend()
plt.show()