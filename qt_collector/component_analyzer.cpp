#include "component_annlyzer.h"

// 同级中相同类型的序号
static QString numAmongOthersWithTheSameClass(const QObject &w)
{
    QObject *p = w.parent();
    if (p == nullptr)
        return QString();

    const QObjectList &childs = p->children();
    int order = 0;
    for (QObject *obj : childs) {
        if (obj == nullptr) {
            continue;
        }
        if (obj == &w) {
            if (order == 0)
                return QString();
            else
                return QString(",%1").arg(order);
            continue;
        }
        if (std::strcmp(obj->metaObject()->className(),w.metaObject()->className()) == 0)
            ++order;
    }
    return QString();
}

// 控件id： objectName ? <class_name=${className}[order]>
static QString qtObjectId(const QObject &w)
{
    const QString name = w.objectName();
    if (name.isEmpty()) {
        return QString("<class_name=%1%2>")
            .arg(w.metaObject()->className(),numAmongOthersWithTheSameClass(w));
    }
    return name;
}

static QString fullQtWidgetId(const QObject &w)
{
    QString res = qtObjectId(w);
    QObject *cur_obj = w.parent();
    while (cur_obj != nullptr) {
        res = qtObjectId(*cur_obj) + "." + res;
        cur_obj = cur_obj->parent();
    }
    return res;
}

QStringList qt_collector::geneQPushButton(QWidget *w) {
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    auto meta = w->metaObject()->className();
    qDebug() << "meta=" << meta;

    if (std::strcmp(meta, "QPushButton") != 0) {
        return res;
    }

    QPushButton* pb = qobject_cast<QPushButton*>(w);
    if (pb == nullptr) {
        return res;
    }

    c.type = qt_collector::Button;
    c.name = fullQtWidgetId(*w);
    c.desc = pb->text();

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}

QStringList qt_collector::geneQLabel(QWidget *w) {
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    auto meta = w->metaObject()->className();
    qDebug() << "meta=" << meta;

    if (std::strcmp(meta, "QLabel") != 0) {
        return res;
    }

    QLabel* pb = qobject_cast<QLabel*>(w);
    if (pb == nullptr) {
        return res;
    }

    c.type = qt_collector::Text;
    c.name = fullQtWidgetId(*w);
    c.desc = pb->text();

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}
