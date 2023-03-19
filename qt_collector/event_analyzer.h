#include <QKeySequence>
#include <QPoint>
#include <QtCore/QDateTime>
#include <QtCore/QEvent>
#include <QtCore/QObject>
#include <QtCore/QString>
#include <functional>

#include <QtEvents>
#include <QDateTime>
#include <QCursor>
#include <QApplication>
#include <cstdlib>
#include <cstring>
#include <QDebug>
#include <QWidget>
#include <QPushButton>

#include <component_annlyzer.h>

namespace qt_collector
{
class Agent;

enum EventType {
    AppStart = 1, // 应用启动
    AppQuit = 2, // 应用退出
    MouseClick = 3, // 鼠标点击
    MouseMove = 4, // 鼠标移动
    KeyClick = 5, // 键盘点击
    Shortcut = 6 // 快捷键
};

enum MouseButtonType {
    Left = 1, // 左键
    Right = 2 // 右键
};

enum MouseClickType {
    One = 1, // 单击
    Two = 2  // 双击
};

enum MouseMoveType {
    Begin = 1, // 起点
    End = 2  // 终点
};

enum KeyClickType {
    Single = 1, // 单键
    Component = 2 // 组合键
};

class UserEventAnalyzer
#ifndef Q_MOC_RUN
    final
#endif
    : public QObject
{
    Q_OBJECT
signals:
    void userEvent(QStringList &);

public:
    UserEventAnalyzer(Agent &agent);

private:
    // 重写事件过滤
    bool eventFilter(QObject *obj, QEvent *event) override;

    // 根据格式生成数据
    QStringList geneDataInForm();
    // 根据位置生成组件信息
    QStringList geneComponent();

private:
    Agent &agent_;
    struct {
        QObject *obj;       //!< the same as in QObject::eventFilter
        QEvent *event;      //!< the same as in QObject::eventFilter

        QPoint globalPos; // 全局坐标
        EventType type; // 事件类型
        MouseClickType mouseClickType; // 鼠标点击类型
        MouseButtonType mouseButtonType; // 鼠标按键类型
        MouseMoveType mouseMoveType; // 鼠标移动类型
        KeyClickType keyClickType; // 键盘点击类型
        QString keyValue; // 键盘输入
    } eventInfo; //事件信息
    struct {
        QEvent::Type type = QEvent::None;
        QDateTime timestamp;
        int key = -1;
    } lastKeyEvent_; // 上次键盘事件
    struct {
        QEvent::Type type = QEvent::None;
        QDateTime timestamp;
        QStringList lastRes; // 上次鼠标移动的记录
    } lastMouseMoveEvent_; // 上次鼠标移动事件
    struct {
        QEvent::Type type = QEvent::None;
        QDateTime timestamp;
        Qt::MouseButton button;
    } lastMouseClickEvent_; // 上次鼠标点击事件
    std::list<ComponentAnalyzer> componentAnalyzer_;
};

}

