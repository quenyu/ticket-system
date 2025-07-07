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
    return 9; // ID, Title, Description, Status, Priority, Department, Assignee, Created At, Updated At
}

QVariant TicketModel::data(const QModelIndex &index, int role) const {
    if (!index.isValid() || index.row() >= m_tickets.size()) return QVariant();
    const TicketItem &t = m_tickets[index.row()];
    if (role == Qt::DisplayRole) {
        switch (index.column()) {
        case 0: return t.id;
        case 1: return t.title;
        case 2: return t.description;
        case 3: return t.status;
        case 4: return t.priority;
        case 5: return t.department;
        case 6: return t.assignee;
        case 7: return t.createdAt.toString("dd:MM:yyyy");
        case 8: return t.updatedAt.toString("dd:MM:yyyy");
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
        case 2: return "Description";
        case 3: return "Status";
        case 4: return "Priority";
        case 5: return "Department";
        case 6: return "Assignee";
        case 7: return "Created At";
        case 8: return "Updated At";
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

QMap<int, QString> TicketItem::statusLabels;
QMap<int, QString> TicketItem::priorityLabels;
QMap<int, QString> TicketItem::departmentNames;

void TicketModel::loadTickets(const QJsonArray& array) {
    QVector<TicketItem> tickets;
    for (const QJsonValue& val : array) {
        if (!val.isObject()) continue;
        QJsonObject obj = val.toObject();
        TicketItem t;
        t.id = obj.value("ticket_id").toString();
        t.title = obj.value("title").toString();
        t.description = obj.value("description").toString();
        t.statusId = obj.value("status_id").toInt(-1);
        t.status = TicketItem::statusLabels.value(t.statusId, QString::number(t.statusId));
        t.priorityId = obj.value("priority_id").toInt(-1);
        t.priority = TicketItem::priorityLabels.value(t.priorityId, QString::number(t.priorityId));
        t.departmentId = obj.value("department_id").toInt(-1);
        t.department = TicketItem::departmentNames.value(t.departmentId, QString::number(t.departmentId));
        t.assignee = obj.value("assignee_name").toString();
        t.assigneeId = obj.value("assignee_id").toString();
        t.creatorId = obj.value("creator_id").toString();
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

QJsonObject TicketItem::toJson() const {
    QJsonObject obj;
    if (!id.isEmpty()) obj["ticket_id"] = id;
    obj["title"] = title;
    obj["description"] = description;
    obj["status_id"] = statusId;
    obj["priority_id"] = priorityId;
    obj["department_id"] = departmentId;
    obj["assignee_id"] = assigneeId;
    obj["creator_id"] = creatorId;
    return obj;
} 