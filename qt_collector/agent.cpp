#include <QCoreApplication>
#include "agent.h"
#include "event_analyzer.h"
#include <QDebug>

using qt_collector::Agent;
using qt_collector::UserEventAnalyzer;

Agent *Agent::gAgent_ = nullptr; // agent 实例

Agent::Agent()
    : eventAnalyzer_(new UserEventAnalyzer(*this))
{
    assert(gAgent_ == nullptr);
    gAgent_ = this;

    // 创建数据目录
    QDir dir(DataFileDir);
    if (!dir.exists()) {
        QDir().mkpath(DataFileDir);
    }

    // 创建数据文件
    qint64 curTimeOfMS = QDateTime::currentDateTime().toMSecsSinceEpoch();
    QString dataFileName = DataFileDir + "/data_" + QString::number(curTimeOfMS, 10) + ".csv";
    qDebug() << "dataFileName = " << dataFileName;
    dataFile.setFileName(dataFileName);
    openSuccess = dataFile.open(QFile::WriteOnly | QFile::Append);
    qDebug() << "openSuccess = " << openSuccess;


    // qApp: 应用实例，应用退出
    connect(qApp, SIGNAL(aboutToQuit()), this, SLOT(onAppAboutToQuit()));
    // 用户事件
    connect(eventAnalyzer_, SIGNAL(userEvent(QStringList &)),
            this, SLOT(onUserEvent(QStringList &)));
    // 安装事件过滤器
    QCoreApplication::instance()->installEventFilter(eventAnalyzer_);

    onAppStart();
}

Agent::~Agent()
{
    if (dataFile.exists() && openSuccess) {
        dataFile.close();
    }
}

void Agent::onUserEvent(QStringList &list)
{
    writeData(list);
}

void Agent::onAppAboutToQuit()
{
    qint64 curTimeOfMS = QDateTime::currentDateTime().toMSecsSinceEpoch();
    QStringList quitData = {QString::number(AppQuit, 10), QString::number(curTimeOfMS, 10)};
    writeData(quitData);
}

void Agent::onAppStart()
{
    writeData(FileHeader);

    qint64 curTimeOfMS = QDateTime::currentDateTime().toMSecsSinceEpoch();
    QStringList startData = {QString::number(AppStart, 10), QString::number(curTimeOfMS, 10)};
    writeData(startData);
}

void Agent::writeData(QStringList &list)
{
    if (!dataFile.exists() || !openSuccess) {
        return ;
    }

    if (list.isEmpty()) {
        return ;
    }

    QTextStream output(&dataFile);
    for (int idx = 0; idx < list.length(); idx++) {
        output << ConvertString2CSV(list[idx]);
        output << ((idx < list.length() - 1) ? "," : "\n");
    }
}

