#pragma once
#include <QTableView>
#include <QPainter>

class TicketTableView : public QTableView {
    Q_OBJECT
public:
    using QTableView::QTableView;
protected:
    void paintEvent(QPaintEvent *event) override {
        QTableView::paintEvent(event);
        qDebug() << "TicketTableView paintEvent, rowCount:" << (model() ? model()->rowCount(rootIndex()) : -1);
        if (model() && model()->rowCount(rootIndex()) == 0) {
            QPainter p(viewport());
            p.setPen(Qt::gray);
            QFont f = font();
            f.setBold(true);
            f.setPointSize(18);
            p.setFont(f);
            QString text = "No tickets in this department";
            QRect r = viewport()->rect();
            p.drawText(r, Qt::AlignCenter, text);
        }
    }
}; 