#include <QtCore/QObject>
#include <QtCore/QEvent>

namespace qt_collector
{
class Agent;

class UserEventAnalyzer
#ifndef Q_MOC_RUN
    final
#endif
    : public QObject
{
    Q_OBJECT
signals:
    void userEvent(const QString &);

public:
    UserEventAnalyzer(Agent &agent, QObject *parent = nullptr);

private:
    bool eventFilter(QObject *obj, QEvent *event) override;

private:
    Agent &agent_;
};


}

