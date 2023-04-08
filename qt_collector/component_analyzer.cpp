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

// 组件id：. 分割
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

static QWidget* searchThroghSuperClassesAndParents(QWidget *widget, const char *wname/*, size_t limit = size_t(-1)*/)
{
    for (size_t i = 0; widget != nullptr /*&& i < limit*/; ++i) {
        const QMetaObject *mo = widget->metaObject();
        while (mo != nullptr && std::strcmp(mo->className(), wname) != 0) {
            mo = mo->superClass();
        }

        if (mo != nullptr) {
            return widget;
        } else {
            widget = qobject_cast<QWidget *>(widget->parent());
        }
    }
    return nullptr;
}

QStringList qt_collector::geneNone(QWidget *w)
{
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    c.type = qt_collector::None;
    c.name = fullQtWidgetId(*w);
    c.desc = "";

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}

QStringList qt_collector::geneQAbstractButton(QWidget *w)
{
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    //qDebug() << fullQtWidgetId(*w);

    QAbstractButton *b = qobject_cast<QAbstractButton *>(w);

    if (b == nullptr) {
        return res;
    }

    c.type = qt_collector::Button;
    c.name = fullQtWidgetId(*w);
    c.desc = b->text();

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}

QStringList qt_collector::geneQComboBox(QWidget *w)
{
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    QWidget *tmp =  searchThroghSuperClassesAndParents(w, "QComboBox");

    if (tmp == nullptr) {
        return res;
    }

    QComboBox *b = qobject_cast<QComboBox *>(tmp);

    if (b == nullptr) {
        return res;
    }

    c.type = qt_collector::Combo;
    c.name = fullQtWidgetId(*tmp);
    c.desc = b->currentText();

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}

QStringList qt_collector::geneText(QWidget *w)
{
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    QLineEdit *e1;
    QTextEdit *e2;
    QPlainTextEdit *e3;
    QLabel *l;

    QWidget * tmp1;
    QWidget * tmp2;

    // QLineEdit
    e1 = qobject_cast<QLineEdit *>(w);
    if (e1 != nullptr) {
        c.desc = e1->text();
        c.name = fullQtWidgetId(*w);
        goto end;
    }

    // QTextEdit
    tmp1 =  searchThroghSuperClassesAndParents(w, "QTextEdit");
    if (tmp1 != nullptr) {
        e2 = qobject_cast<QTextEdit *>(tmp1);
        if (e2 != nullptr) {
            c.desc = e2->toPlainText();
            c.name = fullQtWidgetId(*tmp1);
            goto end;
        }
    }

    // QPlainTextEdit
    tmp2 =  searchThroghSuperClassesAndParents(w, "QPlainTextEdit");
    if (tmp2 != nullptr) {
        e3 = qobject_cast<QPlainTextEdit *>(tmp2);
        if (e3 != nullptr) {
            c.desc = e3->toPlainText();
            c.name = fullQtWidgetId(*tmp2);
            goto end;
        }
    }

    // QLabel
    l = qobject_cast<QLabel *>(w);
    if (l != nullptr) {
        c.desc = l->text();
        c.name = fullQtWidgetId(*w);
        goto end;
    }

    return res;

end:
    c.type = qt_collector::Text;

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}

QStringList qt_collector::geneSpin(QWidget *w)
{
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    QWidget *tmp =  searchThroghSuperClassesAndParents(w, "QAbstractSpinBox");
    if (tmp == nullptr) {
        return res;
    }

    QAbstractSpinBox *b = qobject_cast<QAbstractSpinBox *>(tmp);

    if (b == nullptr) {
        return res;
    }

    c.type = qt_collector::Spin;
    c.name = fullQtWidgetId(*tmp);
    c.desc = b->text();

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}

QStringList qt_collector::geneSlider(QWidget *w)
{
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    QAbstractSlider *b = qobject_cast<QAbstractSlider *>(w);

    if (b == nullptr) {
        return res;
    }

    c.type = qt_collector::Slider;
    c.name = fullQtWidgetId(*w);
    c.desc = QString::number(b->value());

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}

QStringList qt_collector::geneCalendar(QWidget *w)
{
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    QWidget * calendar =  searchThroghSuperClassesAndParents(w, "QCalendarWidget");
    if (calendar == nullptr) {
        return res;
    }

    QCalendarWidget *b = qobject_cast<QCalendarWidget *>(calendar);

    if (b == nullptr) {
        return res;
    }

    c.type = qt_collector::Calendar;
    c.name = fullQtWidgetId(*calendar);
    c.desc = b->selectedDate().toString();

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}

QStringList qt_collector::geneLcd(QWidget *w){
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    QLCDNumber *b = qobject_cast<QLCDNumber *>(w);

    if (b == nullptr) {
        return res;
    }

    c.type = qt_collector::Lcd;
    c.name = fullQtWidgetId(*w);
    c.desc = QString("%1").arg(b->value());

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}

QStringList qt_collector::geneProgress(QWidget *w){
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    QProgressBar *b = qobject_cast<QProgressBar *>(w);

    if (b == nullptr) {
        return res;
    }

    c.type = qt_collector::Progress;
    c.name = fullQtWidgetId(*w);
    c.desc = QString("%1").arg(b->value());

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}

QStringList qt_collector::geneListView(QWidget *w)
{
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    QWidget *listView = searchThroghSuperClassesAndParents(w, "QListView");

    if (listView == nullptr) {
        return res;
    }

    QListView *b = qobject_cast<QListView *>(listView);

    if (b == nullptr) {
        return res;
    }

    QPoint globalPos = QCursor::pos();
    QPoint pos = w->mapFromGlobal(globalPos);

    QModelIndex idx = b->indexAt(pos);

    c.type = qt_collector::List;
    c.name = fullQtWidgetId(*listView);
    c.desc = idx.data().toString();

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}

QStringList qt_collector::geneTreeView(QWidget *w)
{
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    QWidget *tree_view = searchThroghSuperClassesAndParents(w, "QTreeView");

    if (tree_view == nullptr) {
        return res;
    }

    QTreeView *b = qobject_cast<QTreeView *>(tree_view);

    if (b == nullptr) {
        return res;
    }

    QPoint globalPos = QCursor::pos();
    QPoint pos = w->mapFromGlobal(globalPos);

    QModelIndex idx = b->indexAt(pos);

    c.type = qt_collector::Tree;
    c.name = fullQtWidgetId(*tree_view);
    c.desc = idx.data().toString();

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}

QStringList qt_collector::geneTableView(QWidget *w)
{
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    QWidget *table_view = searchThroghSuperClassesAndParents(w, "QTableView");

    if (table_view == nullptr) {
        return res;
    }

    QTableView *b = qobject_cast<QTableView *>(table_view);

    if (b == nullptr) {
        return res;
    }

    QPoint globalPos = QCursor::pos();
    QPoint pos = w->mapFromGlobal(globalPos);

    QModelIndex idx = b->indexAt(pos);

    c.type = qt_collector::Table;
    c.name = fullQtWidgetId(*table_view);
    c.desc = idx.data().toString();

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}

QStringList qt_collector::geneColumnView(QWidget *w)
{
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    QWidget *column_view = searchThroghSuperClassesAndParents(w, "QColumnView");

    if (column_view == nullptr) {
        return res;
    }

    QColumnView *b = qobject_cast<QColumnView *>(column_view);

    if (b == nullptr) {
        return res;
    }

    QPoint globalPos = QCursor::pos();
    QPoint pos = w->mapFromGlobal(globalPos);

    QModelIndex idx = b->indexAt(pos);

    c.type = qt_collector::Column;
    c.name = fullQtWidgetId(*column_view);
    c.desc = idx.data().toString();

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}

QStringList qt_collector::geneAction(QWidget *w)
{
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    QMenu *b = qobject_cast<QMenu *>(w);

    if (b == nullptr) {
        return res;
    }

    QPoint globalPos = QCursor::pos();
    QPoint pos = w->mapFromGlobal(globalPos);

    QAction *act = b->actionAt(pos);
    if (act == nullptr) {
        return res;
    }

    c.type = qt_collector::Action;
    c.name = fullQtWidgetId(*w) + "." + act->objectName();
    c.desc = act->text();

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}

QStringList qt_collector::geneContainer(QWidget *w)
{
    QStringList res = {};
    struct CustomComponent c;

    if (w == nullptr) {
        return res;
    }

    QGroupBox * b1;
    QToolBox * b2;
    QTabWidget * b3;
    QStackedWidget * b4;

    QWidget *group_box;
    QWidget *tool_box;
    QWidget *tab_box;
    QWidget *stack_box;


    group_box = searchThroghSuperClassesAndParents(w, "QGroupBox");
    if (group_box != nullptr) {
        b1 = qobject_cast<QGroupBox *>(group_box);
        if (b1 != nullptr) {
            c.desc = b1->title();
            c.name = fullQtWidgetId(*group_box);
            goto end;
        }
    }

    tool_box = searchThroghSuperClassesAndParents(w, "QToolBox");
    if (tool_box != nullptr) {
        b2 = qobject_cast<QToolBox *>(tool_box);
        if (b2 != nullptr) {
            c.desc = b2->itemText(b2->currentIndex());
            c.name = QString("%1.%2").arg(fullQtWidgetId(*tool_box), QString::number(b2->currentIndex()));
            goto end;
        }
    }

    tab_box = searchThroghSuperClassesAndParents(w, "QTabWidget");
    if (tab_box != nullptr) {
        b3 = qobject_cast<QTabWidget *>(tab_box);
        if (b3 != nullptr) {
            c.desc = b3->tabText(b3->currentIndex());
            c.name = QString("%1.%2").arg(fullQtWidgetId(*tab_box), QString::number(b3->currentIndex()));
            goto end;
        }
    }

    stack_box = searchThroghSuperClassesAndParents(w, "QStackedWidget");
    b4 = qobject_cast<QStackedWidget *>(stack_box);
    if (b4 != nullptr) {
        c.desc = b4->currentWidget()->objectName();
        c.name = QString("%1.%2").arg(fullQtWidgetId(*stack_box), QString::number(b4->currentIndex()));
        goto end;
    }

    return res;
end:

    c.type = qt_collector::Container;

    res.append(QString("%1").arg(c.name));
    res.append(QString("%1").arg(c.type));
    res.append(QString("%1").arg(c.desc));

    return res;
}
