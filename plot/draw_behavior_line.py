import matplotlib.pyplot as plt

data_str = "(0,1680873522817)(1,1680873651246)(4,1680873666334)(4,1680873670439)(0,1680873678367)(0,1680873690679)(2,1680873722385)(2,1680873722416)(2,1680873753086)(0,1680873759442)(0,1680873773799)(0,1680873833542)"
data_list = [tuple(map(int, s.strip("()").split(","))) for s in data_str.split(")") if s]

x = [d[1] for d in data_list]
y = [d[0] for d in data_list]

plt.rcParams['font.family'] = 'SimHei'
plt.rcParams['axes.unicode_minus'] = False

plt.plot(x, y, drawstyle='steps-post')
plt.xlabel("时间")
plt.ylabel("行为")
plt.yticks([0,1,2,3,4,5],['未操作','编码','测试','调试','浏览','使用工具'])
plt.show()