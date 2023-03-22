import matplotlib.pyplot as plt
 
plt.rcParams['font.family'] = 'SimHei'
plt.rcParams['axes.unicode_minus'] = False
 
plt.figure(figsize=(20, 10), dpi=100)
x = ["0min", "1min", "2min", "3min", "4min", "5min", "6min", "7min", "8min", "9min", "10min", "11min", "12min"]
y = [1, 2, 2, 0, 0, 2, 2, 3, 3, 1, 2, 3, 2]
map = {
    0: "未操作",
    1: "打开项目",
    2: "查看地形图",
    3: "设置参数"
}
ylabel = []
for idx in range(len(y)):
    ylabel.append(map[y[idx]])

plt.xlabel("时间")
plt.ylabel("行为")

plt.yticks(y, ylabel)

plt.plot(x, y)

plt.legend()
plt.show()