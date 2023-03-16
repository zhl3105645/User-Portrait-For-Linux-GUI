#include <QCoreApplication>
#include "agent.h"
#include "event_analyzer.h"

using qt_collector::Agent;
using qt_collector::UserEventAnalyzer;

Agent *Agent::gAgent_ = nullptr; // agent 实例

Agent::Agent()
    : eventAnalyzer_(new UserEventAnalyzer(*this, this))
{
    assert(gAgent_ == nullptr);
    gAgent_ = this;
    // qApp: 应用实例，应用退出
    connect(qApp, SIGNAL(aboutToQuit()), this, SLOT(onAppAboutToQuit()));
    // 用户事件
    connect(eventAnalyzer_, SIGNAL(userEvent(const QString &)),
            this, SLOT(onUserEvent(const QString &)));
    // 安装事件过滤器
    QCoreApplication::instance()->installEventFilter(eventAnalyzer_);
}

Agent::~Agent()
{

}

void Agent::onAppAboutToQuit()
{
    qDebug("程序退出");
}

void Agent::onUserEvent(const QString &)
{
    qDebug("接收到事件");
}
