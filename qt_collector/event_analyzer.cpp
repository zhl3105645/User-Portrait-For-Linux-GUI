#include "event_analyzer.h"

using qt_collector::UserEventAnalyzer;
using qt_collector::ComponentAnalyzer;
using qt_collector::Agent;
using qt_collector::CustomComponent;

static constexpr int repeatKeyEventTimeoutMs = 20;
static constexpr int repeatMouseClickEventTimeoutMs = 20;
static constexpr int repeatMouseMoveEventTimeoutMs = 500;

// 键盘事件解析
static QKeySequence parseKeyReleaseEvent(QKeyEvent *keyEvent)
{
    if (keyEvent == nullptr) {
        return QStringLiteral();
    }

    if (keyEvent->key() == 0)
        return QStringLiteral("//some special key");

    int modifiers[4] = {0, 0, 0, 0};
    int curMod = 0;

    if (keyEvent->modifiers() & Qt::ShiftModifier)
        modifiers[curMod++] = Qt::ShiftModifier;

    if (keyEvent->modifiers() & Qt::AltModifier)
        modifiers[curMod++] = Qt::AltModifier;

    if (keyEvent->modifiers() & Qt::ControlModifier)
        modifiers[curMod++] = Qt::ControlModifier;

    if (keyEvent->modifiers() & Qt::MetaModifier)
        modifiers[curMod++] = Qt::MetaModifier;

    QKeySequence keySeq;
    switch (curMod) {
    case 1:
        if (keyEvent->key() != Qt::Key_unknown)
            keySeq = QKeySequence(modifiers[0], keyEvent->key());
        else
            keySeq = QKeySequence(modifiers[0]);
        break;
    case 2:
        if (keyEvent->key() != Qt::Key_unknown)
            keySeq = QKeySequence(modifiers[0], modifiers[1], keyEvent->key());
        else
            keySeq = QKeySequence(modifiers[0], modifiers[1]);
        break;
    case 3:
        if (keyEvent->key() != Qt::Key_unknown)
            keySeq = QKeySequence(modifiers[0], modifiers[1], modifiers[2], keyEvent->key());
        else
            keySeq = QKeySequence(modifiers[0], modifiers[1], modifiers[2]);
        break;
    case 4:
        keySeq = QKeySequence(keyEvent->text());
        break;
    case 0:
    default:
        keySeq = QKeySequence(keyEvent->key());
        break;
    }

    return keySeq;
}

QStringList geneQPushButton(QWidget *w);

UserEventAnalyzer::UserEventAnalyzer(Agent &agent)
    : agent_(agent)
{
    for (auto &&fun: {geneQPushButton, geneQLabel}) {
        componentAnalyzer_.emplace_back(fun);
    }
}

bool UserEventAnalyzer::eventFilter(QObject *obj, QEvent *event)
{
    if (obj == nullptr || event == nullptr) {
        return QObject::eventFilter(obj, event);
    }

    // 鼠标移动停止: 鼠标点击 键盘输入等事件
    if (lastMouseMoveEvent_.type == QEvent::MouseMove && event->type() != QEvent::MouseMove
        && (event->type() == QEvent::KeyPress || event->type() == QEvent::MouseButtonPress || event->type() == QEvent::MouseButtonDblClick)) {
        lastMouseMoveEvent_.type = QEvent::None;

        //qDebug() << "移动结束：" << lastMouseMoveEvent_.lastRes;
        // 保存结束数据
        emit userEvent(lastMouseMoveEvent_.lastRes);
    }

    switch (event->type()) {
    case QEvent::KeyPress:
        // 单击点击时不记入，松开时记录数据
        break;
    case QEvent::KeyRelease: {
        QDateTime now = QDateTime::currentDateTime();
        auto keyEvent = static_cast<QKeyEvent *>(event);

        if (keyEvent == nullptr) {
            break;
        }

        // 重复事件不记录
        if (lastKeyEvent_.type == event->type() && lastKeyEvent_.key == keyEvent->key()
                && std::llabs(now.msecsTo(lastKeyEvent_.timestamp)) < repeatKeyEventTimeoutMs) {
            break;
        }

        // ignore special keys alone
        if (keyEvent->key() == Qt::Key_Shift || keyEvent->key() == Qt::Key_Alt
            || keyEvent->key() == Qt::Key_Control || keyEvent->key() == Qt::Key_Meta) {
            break;
        }

        //qDebug() << "keyEvent.type=" << keyEvent->type() << "; keyEvent.key=" << keyEvent->key();

        QKeySequence keySeq = parseKeyReleaseEvent(keyEvent);

        eventInfo.obj = obj;
        eventInfo.event = event;
        eventInfo.globalPos = QCursor::pos();

        eventInfo.type = KeyClick;
        eventInfo.keyClickType = (keySeq.count() == 1 ? Single : Component);
        eventInfo.keyValue = keySeq.toString();

        QStringList res = geneDataInForm();

        lastKeyEvent_.type = event->type();
        lastKeyEvent_.timestamp = now;
        lastKeyEvent_.key = keyEvent->key();

        emit userEvent(res);
        break;
    }
    case QEvent::MouseMove: {
        eventInfo.obj = obj;
        eventInfo.event = event;
        eventInfo.globalPos = QCursor::pos();

        eventInfo.type = MouseMove;

        QDateTime now = QDateTime::currentDateTime();
        // 鼠标移动继续: 上次类型相同 & 时间间隔小于300ms
        if (lastMouseMoveEvent_.type == QEvent::MouseMove && std::llabs(now.msecsTo(lastMouseMoveEvent_.timestamp)) < repeatMouseMoveEventTimeoutMs) {
            // 每次假设结束
            eventInfo.mouseMoveType = End;

            QStringList res = geneDataInForm();
            lastMouseMoveEvent_.lastRes = res;
            lastMouseMoveEvent_.type = QEvent::MouseMove;
            lastMouseMoveEvent_.timestamp = now;
        } else if (lastMouseMoveEvent_.type == QEvent::MouseMove) {
            // 时间间隔过久
            // 保存结束数据
            //qDebug() << "移动结束：" << lastMouseMoveEvent_.lastRes;
            emit userEvent(lastMouseMoveEvent_.lastRes);
            // 重新开始移动
            eventInfo.mouseMoveType = Begin;

            QStringList res = geneDataInForm();
            lastMouseMoveEvent_.lastRes = res;
            lastMouseMoveEvent_.type = QEvent::MouseMove;
            lastMouseMoveEvent_.timestamp = now;


            //qDebug() << "移动开始：" << res;

            // 保存移动开始数据
            emit userEvent(res);
        } else {
            // 鼠标移动开始
            eventInfo.mouseMoveType = Begin;

            QStringList res = geneDataInForm();
            lastMouseMoveEvent_.lastRes = res;
            lastMouseMoveEvent_.type = QEvent::MouseMove;
            lastMouseMoveEvent_.timestamp = now;


            //qDebug() << "移动开始：" << res;

            // 保存移动开始数据
            emit userEvent(res);
        }
        break;
    }
    case QEvent::MouseButtonPress:
        // 单击点击时不记入，松开时记录数据
        break;
    case QEvent::MouseButtonRelease:
    case QEvent::MouseButtonDblClick: {
        QDateTime now = QDateTime::currentDateTime();
        QMouseEvent *mouseEvent = static_cast<QMouseEvent *>(event);

        if (mouseEvent == nullptr) {
            break;
        }

        // 重复事件不记录
        if (lastMouseClickEvent_.type == event->type() &&  mouseEvent->button() == lastMouseClickEvent_.button
                && std::llabs(now.msecsTo(lastMouseClickEvent_.timestamp)) < repeatMouseClickEventTimeoutMs) {
            break;
        }

        //qDebug() << "mouseEvent->type() = " << mouseEvent->type() << "clickPos = " << mouseEvent->globalPos();

        QPoint clickPos = mouseEvent->globalPos();

        eventInfo.obj = obj;
        eventInfo.event = event;
        eventInfo.globalPos = clickPos;

        eventInfo.type = MouseClick;
        eventInfo.mouseClickType = (event->type() == QEvent::MouseButtonDblClick ? Two : One);
        if (mouseEvent->button() == Qt::LeftButton) {
            eventInfo.mouseButtonType = Left;
        } else if (mouseEvent->button() == Qt::RightButton) {
            eventInfo.mouseButtonType = Right;
        }

        QStringList res = geneDataInForm();

        lastMouseClickEvent_.type = event->type();
        lastMouseClickEvent_.timestamp = now;
        lastMouseClickEvent_.button = mouseEvent->button();

        emit userEvent(res);
    }
    default:
        break;
    }

    return QObject::eventFilter(obj, event);
}

QStringList UserEventAnalyzer::geneDataInForm()
{
    QStringList res = {};

    // 公共参数
    if (eventInfo.type == KeyClick || eventInfo.type == MouseClick || eventInfo.type == MouseMove) {
        res.append(QString::number(eventInfo.type, 10));
        res.append(QString::number(QDateTime::currentDateTime().toMSecsSinceEpoch(), 10));
        res.append(QString("(%1,%2)").arg(eventInfo.globalPos.x()).arg(eventInfo.globalPos.y()));
    } else {
        res.append(QString());
        res.append(QString());
        res.append(QString());
    }


    switch (eventInfo.type) {
    case KeyClick: {
        res.append(QString());
        res.append(QString());
        res.append(QString());
        res.append(QString::number(eventInfo.keyClickType, 10));
        res.append(eventInfo.keyValue);

        break;
    }
    case MouseClick: {
        res.append(QString::number(eventInfo.mouseClickType, 10));
        res.append(QString::number(eventInfo.mouseButtonType, 10));
        res.append(QString());
        res.append(QString());
        res.append(QString());

        break;
    }
    case MouseMove: {
        res.append(QString());
        res.append(QString());
        res.append(QString::number(eventInfo.mouseMoveType, 10));
        res.append(QString());
        res.append(QString());
        break;
    }
    default:
        res.append(QString());
        res.append(QString());
        res.append(QString());
        res.append(QString());
        res.append(QString());
        break;
    }

    // 公共参数
    if (eventInfo.type == KeyClick || eventInfo.type == MouseClick || eventInfo.type == MouseMove) {

        QStringList c = geneComponent();
        if (c.length() < 3) {
            return res;
        }

        res.append(QString("%1").arg(c[0]));
        res.append(QString("%1").arg(c[1]));
        res.append(QString("%1").arg(c[2]));
    } else {
        res.append(QString());
        res.append(QString());
        res.append(QString());
    }

    return res;
}

QStringList UserEventAnalyzer::geneComponent() {
    QWidget *w = QApplication::widgetAt(eventInfo.globalPos);
    if (w == nullptr) {
        return QStringList();
    }

    QStringList res;
    for (auto &&fun : componentAnalyzer_) {
        res = fun(w);
        if (res.length() == 3) {
            return res;
        }
    }

    return QStringList();
}




