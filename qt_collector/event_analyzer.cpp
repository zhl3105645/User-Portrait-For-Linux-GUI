#include <QtCore/QObject>
#include <QtCore/QEvent>
#include "event_analyzer.h"

#include "agent.h"

using qt_event_collector::UserEventAnalyzer;
using qt_event_collector::Agent;

UserEventAnalyzer::UserEventAnalyzer(qt_event_collector::Agent &agent, QObject *parent)
    : QObject(parent), agent_(agent)
{

}

bool UserEventAnalyzer::eventFilter(QObject *obj, QEvent *event)
{
    switch (event->type()) {
    case QEvent::KeyPress:
    case QEvent::KeyRelease:
    case QEvent::MouseMove:
    case QEvent::MouseButtonPress:
    case QEvent::MouseButtonRelease:
    case QEvent::MouseButtonDblClick:
        emit userEvent("测试");
    default:
        break;
    }

    return QObject::eventFilter(obj, event);
}
