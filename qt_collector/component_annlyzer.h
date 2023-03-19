#include <QStringList>
#include <functional>
#include <QtWidgets>
#include <cstdlib>
#include <cstring>

namespace qt_collector {

enum CustomComponentType {
    None = -1, // 未定义组件
    Button = 1, // 按钮，QPushButton QToolButton QRadioButton QCheckBox
    Combo = 2, // 下拉组合框   QComboBox QFontComboBox
    Text = 3, // 文本编辑 QLineEdit QTextEdit QPlainTextEdit QLabel
    Spin = 4, // 滚轮 QSpinBox QDoubleSpinBox
    Slider = 5, // 滑块 QDial QSlider QScrollBar
    Calendar = 6, // 日历
    Lcd = 7, // lcd数字
    Progress = 8, // 进度条
    List = 9, // 列表视图
    Tree = 10, // 树状视图
    Table = 11, // 表视图
    Column = 12, // 列视图
    Action = 13, // 命令
    Container = 14, // 容器
};

struct CustomComponent {
    QString name; // 组件唯一标识
    CustomComponentType type;
    QString desc; // 组件描述
};

using ComponentAnalyzer = std::function<QStringList(QWidget *)>;

QStringList geneNone(QWidget *w);

QStringList geneQAbstractButton(QWidget *w);

QStringList geneQComboBox(QWidget *w);

QStringList geneText(QWidget *w);

QStringList geneSpin(QWidget *w);

QStringList geneSlider(QWidget *w);

QStringList geneCalendar(QWidget *w);

QStringList geneLcd(QWidget *w);

QStringList geneProgress(QWidget *w);

QStringList geneListView(QWidget *w);

QStringList geneTreeView(QWidget *w);

QStringList geneTableView(QWidget *w);

QStringList geneColumnView(QWidget *w);

QStringList geneAction(QWidget *w);

QStringList geneContainer(QWidget *w);
}
