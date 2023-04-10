from sklearn.decomposition import PCA
import numpy as np
import pandas as pd
from sklearn.preprocessing import StandardScaler

num_rows = 1000
num_cols = 5
col_names = ['a', 'b', 'c', 'd', 'e']
low = 0
high = 10000

data = np.random.randint(low, high, size=(num_rows, num_cols))
df = pd.DataFrame(data, columns=col_names)

df.to_csv('./cluster/test_data.csv', index=False)