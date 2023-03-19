#include <QStringList>
#include <functional>
#include <QtWidgets>
#include <cstdlib>
#include <cstring>

namespace qt_collector {

enum CustomComponentType {
    Button = 1, // 按钮，QPushButton QToolButton
    Combo = 2, // 下拉组合框 QRadioButton QCheckBox QComboBox QFontComboBox
    Text = 3, // 文本编辑 QLineEdit QTextEdit QPlainTextEdit QLabel
    List = 4, // 列表视图
    Tree = 5, // 树状视图
    Table = 6, // 表视图
    Action = 7, // 命令 QAction
};

struct CustomComponent {
    QString name; // 组件唯一标识
    CustomComponentType type;
    QString desc; // 组件描述
};

using ComponentAnalyzer = std::function<QStringList(QWidget *)>;

QStringList geneQPushButton(QWidget *w);

QStringList geneQLabel(QWidget *w);

}
