#pragma once
#include <QAbstractTableModel>
#include <QVector>
#include <QString>
#include <QDateTime>
#include <QStyledItemDelegate>
#include <QJsonArray>

struct TicketItem {
    QString id;
    QString title;
    QString status;
    QString priority;
    QString department;
    QString assignee;
    QDateTime createdAt;
    QDateTime updatedAt;
};

class TicketModel : public QAbstractTableModel {
    Q_OBJECT
public:
    explicit TicketModel(QObject *parent = nullptr);
    int rowCount(const QModelIndex &parent = QModelIndex()) const override;
    int columnCount(const QModelIndex &parent = QModelIndex()) const override;
    QVariant data(const QModelIndex &index, int role = Qt::DisplayRole) const override;
    QVariant headerData(int section, Qt::Orientation orientation, int role) const override;
    void setTickets(const QVector<TicketItem> &tickets);
    void loadTickets(const QJsonArray& array);
    TicketItem getTicket(int row) const;
private:
    QVector<TicketItem> m_tickets;
};

class TicketBadgeDelegate : public QStyledItemDelegate {
    Q_OBJECT
public:
    using QStyledItemDelegate::QStyledItemDelegate;
    void paint(QPainter *painter, const QStyleOptionViewItem &option, const QModelIndex &index) const override;
    QSize sizeHint(const QStyleOptionViewItem &option, const QModelIndex &index) const override;
}; 