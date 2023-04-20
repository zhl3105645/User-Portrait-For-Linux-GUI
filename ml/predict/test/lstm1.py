import torch
import torch.nn as nn
import pandas as pd
from sklearn.model_selection import train_test_split
from torch.utils.data import TensorDataset, DataLoader
import matplotlib.pyplot as plt
import numpy as np
from pylab import mpl
mpl.rcParams['font.sans-serif'] = ['FangSong']
mpl.rcParams['axes.unicode_minus'] = False

## 全局变量   
sequence_length=16 # 序列长度
num_classes=69 # 事件类型种类 66 + 1 + 2
input_size = 1          
hidden_size = 16         
num_layers = 2
lr=0.01

num_epochs=300
## 

def loadData():
    data = pd.read_csv('./predict/data/test_data.csv')
    sequences = data.iloc[:, :15].values
    labels = data.iloc[:, 15].values

    # 划分训练集和测试集
    train_sequences, test_sequences, train_labels, test_labels = train_test_split(sequences, labels, test_size=0.2, random_state=42)

    # 划分训练集和验证集
    train_sequences, val_sequences, train_labels, val_labels = train_test_split(train_sequences, train_labels, test_size=0.25, random_state=42)

    # 创建 TensorDataset 对象
    train_dataset = TensorDataset(torch.from_numpy(train_sequences).to(torch.float32), torch.from_numpy(train_labels).to(torch.long))
    val_dataset = TensorDataset(torch.from_numpy(val_sequences).to(torch.float32), torch.from_numpy(val_labels).to(torch.long))
    test_dataset = TensorDataset(torch.from_numpy(test_sequences).to(torch.float32), torch.from_numpy(test_labels).to(torch.long))

    # 创建 DataLoader 对象
    batch_size = 32
    train_loader = DataLoader(train_dataset, batch_size=batch_size, shuffle=True)
    val_loader = DataLoader(val_dataset, batch_size=batch_size, shuffle=False)
    test_loader = DataLoader(test_dataset, batch_size=batch_size, shuffle=False)
    return train_loader, val_loader, test_loader


class LSTMModel(nn.Module):
    def __init__(self, input_size, hidden_size, num_layers, num_classes):
        super(LSTMModel, self).__init__()
        self.hidden_size = hidden_size
        self.num_layers = num_layers
        self.lstm = nn.LSTM(input_size, hidden_size, num_layers, batch_first=True)
        self.fc = nn.Linear(hidden_size, num_classes)
    
    def forward(self, x):
        # 初始化隐藏状态和细胞状态
        h0 = torch.zeros(self.num_layers, x.size(0), self.hidden_size)
        c0 = torch.zeros(self.num_layers, x.size(0), self.hidden_size)
        
        # 前向传播LSTM
        out, _ = self.lstm(x, (h0, c0))
        
        # 解码最后一个时间步的隐藏状态
        out = self.fc(out[:, -1, :])
        return out

# 训练模型
def train(model, train_loader, criterion, optimizer, num_epochs):
    l = []
    for epoch in range(num_epochs):
        for i, (sequences, labels) in enumerate(train_loader):
            sequences = sequences.reshape(-1, sequence_length-1, input_size)
            # 前向传播
            outputs = model(sequences)
            # print(outputs.shape)
            # print(labels.shape)
            loss = criterion(outputs, labels)
            
            # 反向传播和优化
            optimizer.zero_grad()
            loss.backward()
            optimizer.step()

        if epoch % 100 == 0:
            print("Iteration: {} loss {}".format(epoch, loss.item()))
        l.append(loss.item())
    # 绘制损失函数
    plt.plot(l,'r')
    plt.xlabel('训练次数')
    plt.ylabel('loss')
    plt.title('损失函数下降曲线')
    plt.show()

# 验证模型
def validate(model, val_loader, criterion):
    with torch.no_grad():
        correct = 0
        total = 0
        for sequences, labels in val_loader:
            sequences = sequences.reshape(-1, sequence_length-1, input_size)
            outputs = model(sequences)
            _, predicted = torch.max(outputs.data, 1)
            total += labels.size(0)
            correct += (predicted == labels).sum().item()
        accuracy = 100 * correct / total
        print('Accuracy: {} %'.format(accuracy))

# 测试模型
def test(model, test_loader, criterion):
    with torch.no_grad():
        correct = 0
        total = 0
        for sequences, labels in test_loader:
            sequences = sequences.reshape(-1, sequence_length-1, input_size)
            outputs = model(sequences)
            _, predicted = torch.max(outputs.data, 1)
            total += labels.size(0)
            correct += (predicted == labels).sum().item()
        accuracy = 100 * correct / total
        print('Test Accuracy: {} %'.format(accuracy))

# 测试
def predict():
    model = LSTMModel(input_size=input_size, hidden_size=hidden_size, num_layers=num_layers, num_classes=num_classes)
    model.load_state_dict(torch.load("./predict/model/model"))  
    # 预测下一次事件各个编号的出现概率
    str = "59.0,59.0,59.0,59.0,59.0,59.0,1.0,0.0,60.0,61.0,61.0,62.0,63.0,63.0,64.0"
    arr = str.split(",")
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
    print(next_event_prob)


    # 绘制条形图
    plt.bar(range(len(next_event_prob)), next_event_prob)
    plt.xticks(range(len(next_event_prob)),  rotation='270')
    plt.xlabel('Event')
    plt.ylabel('Probability')
    plt.show()

def run():
    # 加载数据
    print('begin load data.....')
    train_loader, val_loader, test_loader = loadData()

    # 实例化模型
    model = LSTMModel(input_size=input_size, hidden_size=hidden_size, num_layers=num_layers, num_classes=num_classes)

    # 定义损失函数和优化器
    criterion = nn.CrossEntropyLoss()
    optimizer = torch.optim.Adam(model.parameters(), lr=lr)

    # 训练模型
    print('begin train model.....')
    train(model, train_loader, criterion, optimizer, num_epochs)

    # 保存模型参数
    torch.save(model.state_dict(), "./predict/model/model")

    # 加载模型参数
    # the_model = TheModelClass(*args, **kwargs)
    # the_model.load_state_dict(torch.load(PATH))

    # 验证模型
    print('begin validate model.....')
    validate(model, val_loader, criterion)

    # 开始测试
    print('begin test model.....')
    test(model, test_loader, criterion)

if __name__ == '__main__':
    #run()
    predict()