#include <QtCore/QObject>

namespace qt_collector
{

class UserEventAnalyzer;


class Agent
#ifndef Q_MOC_RUN
    final
#endif
    : public QObject
{
    Q_OBJECT
public:
    explicit Agent();
    ~Agent();

    // Agent is neither copyable nor movable.
    Agent(const Agent &) = delete;
    Agent &operator=(const Agent &) = delete;

    static Agent *instance() { return gAgent_; }

private slots:
    void onUserEvent(const QString &);
    void onAppAboutToQuit();

private:
    static Agent *gAgent_;
    qt_collector::UserEventAnalyzer *eventAnalyzer_ = nullptr;
};

} // namespace qt_event_collector


