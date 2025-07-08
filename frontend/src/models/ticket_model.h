#pragma once

#include <QAbstractTableModel>
#include <QStyledItemDelegate>
#include <QVector>
#include <QDateTime>
#include <QJsonObject>
#include <QJsonArray>
#include <QPainter>
#include <QApplication>

struct TicketItem {
    QString id;
    QString title;
    QString description;
    int statusId;
    QString status;
    int priorityId;
    QString priority;
    int departmentId;
    QString department;
    QString assignee;
    QString assigneeId;
    QString creatorId;
    QDateTime createdAt;
    QString createdAtRaw;
    QDateTime updatedAt;

    static QMap<int, QString> statusLabels;
    static QMap<int, QString> priorityLabels;
    static QMap<int, QString> departmentNames;

    QJsonObject toJson() const;
};

Q_DECLARE_METATYPE(TicketItem)

class TicketModel : public QAbstractTableModel {
    Q_OBJECT

public:
    explicit TicketModel(QObject *parent = nullptr);

    int rowCount(const QModelIndex &parent = QModelIndex()) const override;
    int columnCount(const QModelIndex &parent = QModelIndex()) const override;
    QVariant data(const QModelIndex &index, int role = Qt::DisplayRole) const override;
    QVariant headerData(int section, Qt::Orientation orientation, int role = Qt::DisplayRole) const override;

    void loadTickets(const QJsonArray& array);
    TicketItem getTicket(int row) const;
    void setTickets(const QVector<TicketItem> &tickets);

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
class TicketPriorityDelegate : public QStyledItemDelegate {
    Q_OBJECT
public:
    using QStyledItemDelegate::QStyledItemDelegate;
    void paint(QPainter *painter, const QStyleOptionViewItem &option, const QModelIndex &index) const override;
    QSize sizeHint(const QStyleOptionViewItem &option, const QModelIndex &index) const override;
};
