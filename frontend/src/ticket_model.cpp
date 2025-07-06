#include "ticket_model.h"
#include <QPainter>
#include <QApplication>
#include <QBrush>
#include <QColor>
#include <QJsonArray>
#include <QJsonObject>

TicketModel::TicketModel(QObject *parent)
    : QAbstractTableModel(parent) {}

int TicketModel::rowCount(const QModelIndex &parent) const {
    if (parent.isValid()) return 0;
    return m_tickets.size();
}

int TicketModel::columnCount(const QModelIndex &parent) const {
    if (parent.isValid()) return 0;
    return 8; // ID, Title, Status, Priority, Department, Assignee, Created At, Updated At
}

QVariant TicketModel::data(const QModelIndex &index, int role) const {
    if (!index.isValid() || index.row() >= m_tickets.size()) return QVariant();
    const TicketItem &t = m_tickets[index.row()];
    if (role == Qt::DisplayRole) {
        switch (index.column()) {
        case 0: return t.id;
        case 1: return t.title;
        case 2: return t.status;
        case 3: return t.priority;
        case 4: return t.department;
        case 5: return t.assignee;
        case 6: return t.createdAt.toString(Qt::ISODate);
        case 7: return t.updatedAt.toString(Qt::ISODate);
        }
    }
    if (role == Qt::UserRole) {
        return QVariant::fromValue(t);
    }
    return QVariant();
}

QVariant TicketModel::headerData(int section, Qt::Orientation orientation, int role) const {
    if (orientation == Qt::Horizontal && role == Qt::DisplayRole) {
        switch (section) {
        case 0: return "ID";
        case 1: return "Title";
        case 2: return "Status";
        case 3: return "Priority";
        case 4: return "Department";
        case 5: return "Assignee";
        case 6: return "Created At";
        case 7: return "Updated At";
        }
    }
    return QVariant();
}

void TicketModel::setTickets(const QVector<TicketItem> &tickets) {
    beginResetModel();
    m_tickets = tickets;
    endResetModel();
}

void TicketBadgeDelegate::paint(QPainter *painter, const QStyleOptionViewItem &option, const QModelIndex &index) const {
    QString value = index.data(Qt::DisplayRole).toString();
    QRect r = option.rect;
    QColor bg = Qt::white;
    QColor fg = Qt::black;
    if (index.column() == 2) { // Status
        if (value == "Opened" || value == "Open") { bg = QColor("#e74c3c"); fg = Qt::white; value = "Open"; }
        else if (value == "Closed") { bg = QColor("#27ae60"); fg = Qt::white; }
        else if (value == "In Progress") { bg = QColor("#2980b9"); fg = Qt::white; }
    } else if (index.column() == 3) { // Priority
        if (value == "Low") { bg = QColor("#f1c40f"); fg = Qt::black; }
        else if (value == "Medium") { bg = QColor("#f39c12"); fg = Qt::black; }
        else if (value == "High") { bg = QColor("#e67e22"); fg = Qt::white; }
    }
    painter->save();
    painter->setRenderHint(QPainter::Antialiasing, true);
    painter->setPen(Qt::NoPen);
    painter->setBrush(bg);
    QRect badgeRect = r.adjusted(12, 8, -12, -8);
    painter->drawRoundedRect(badgeRect, 6, 6);
    painter->setPen(fg);
    QFont f = option.font;
    f.setBold(true);
    f.setPointSize(13);
    painter->setFont(f);
    painter->drawText(badgeRect, Qt::AlignCenter, value);
    painter->restore();
}

QSize TicketBadgeDelegate::sizeHint(const QStyleOptionViewItem &option, const QModelIndex &index) const {
    QSize sz = QStyledItemDelegate::sizeHint(option, index);
    sz.setHeight(28);
    return sz;
}

void TicketModel::loadTickets(const QJsonArray& array) {
    QVector<TicketItem> tickets;
    for (const QJsonValue& val : array) {
        if (!val.isObject()) continue;
        QJsonObject obj = val.toObject();
        TicketItem t;
        t.id = obj.value("ticket_id").toString();
        t.title = obj.value("title").toString();
        t.status = obj.value("status_label").toString();
        t.priority = obj.value("priority_label").toString();
        t.department = obj.value("department_name").toString();
        t.assignee = obj.value("assignee_name").toString();
        t.createdAt = QDateTime::fromString(obj.value("created_at").toString(), Qt::ISODate);
        t.updatedAt = QDateTime::fromString(obj.value("updated_at").toString(), Qt::ISODate);
        tickets.append(t);
    }
    setTickets(tickets);
}

TicketItem TicketModel::getTicket(int row) const {
    if (row < 0 || row >= m_tickets.size()) return TicketItem{};
    return m_tickets[row];
} 