#include "ticket_model.h"
#include <QPainter>
#include <QApplication>
#include <QBrush>
#include <QColor>
#include <QJsonArray>
#include <QJsonObject>
#include <QPainterPath>

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
        case 6: return t.createdAt.toString("MM/dd/yyyy");
        case 7: return t.updatedAt.toString("MM/dd/yyyy");
        }
    }
    if (role == Qt::TextAlignmentRole) {
        return Qt::AlignCenter;
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

QSize TicketBadgeDelegate::sizeHint(const QStyleOptionViewItem &option, const QModelIndex &index) const {
    QString value = index.data(Qt::DisplayRole).toString();

    QFont f = option.font;
    f.setBold(true);
    f.setPointSize(9);
    QFontMetrics fm(f);

    int textWidth = fm.horizontalAdvance(value);
    int textHeight = fm.height();

    int horizontalPadding = 10;
    int verticalPadding = 3;

    int badgeWidth = textWidth + 2 * horizontalPadding;
    int badgeHeight = textHeight + 2 * verticalPadding;

    int cellPadding = 2;

    QSize size(badgeWidth + 2 * cellPadding, badgeHeight + 2 * cellPadding);

    return size;
}


void TicketBadgeDelegate::paint(QPainter *painter, const QStyleOptionViewItem &option, const QModelIndex &index) const {
    if (index.column() == 2) {
        QString value = index.data(Qt::DisplayRole).toString();
        QColor fg;
        QColor fill;

        if (value.compare("Opened", Qt::CaseInsensitive) == 0 || value.compare("Open", Qt::CaseInsensitive) == 0) {
            fill = QColor("#D32F2F");
            value = "Open";
        } else if (value.compare("Closed", Qt::CaseInsensitive) == 0) {
            fill = QColor("#4CAF50");
            value = "Closed";
        } else {
            fill = Qt::gray;
        }

        painter->save();
        painter->setRenderHint(QPainter::Antialiasing, true);

        if (option.state & QStyle::State_Selected) {
            painter->fillRect(option.rect, option.palette.highlight());
            fg = option.palette.highlightedText().color();
        } else {
            fg = Qt::white;
        }

        QFont f = option.font;
        f.setBold(true);
        f.setPointSize(9);
        painter->setFont(f);
        QFontMetrics fm(f);

        int textWidth = fm.horizontalAdvance(value);
        int textHeight = fm.height();
        int horizontalPadding = 10;
        int verticalPadding = 3;
        int badgeWidth = textWidth + 2 * horizontalPadding;
        int badgeHeight = textHeight + 2 * verticalPadding;
        int badgeX = option.rect.x() + (option.rect.width() - badgeWidth) / 2;
        int badgeY = option.rect.y() + (option.rect.height() - badgeHeight) / 2;

        QRectF badgeRect(badgeX, badgeY, badgeWidth, badgeHeight);

        painter->setPen(Qt::NoPen);
        painter->setBrush(fill);
        painter->drawRoundedRect(badgeRect, badgeRect.height() / 2, badgeRect.height() / 2);

        painter->setPen(fg);
        painter->drawText(badgeRect, Qt::AlignCenter, value);

        painter->restore();
    } else {
        QStyledItemDelegate::paint(painter, option, index);
    }
}


void TicketPriorityDelegate::paint(QPainter *painter, const QStyleOptionViewItem &option, const QModelIndex &index) const {
    QString value = index.data(Qt::DisplayRole).toString();

    painter->save();
    painter->setRenderHint(QPainter::Antialiasing, true);

    if (option.state & QStyle::State_Selected) {
        painter->fillRect(option.rect, option.palette.highlight());
    }
    
    int iconSize = 10;
    int iconPadding = 5;
    int iconY = option.rect.y() + (option.rect.height() - iconSize) / 2;
    int iconX = option.rect.x() + iconPadding;

    QPainterPath path;
    path.moveTo(iconX, iconY + iconSize);
    path.lineTo(iconX + iconSize, iconY + iconSize);
    path.lineTo(iconX + (iconSize / 2.0), iconY);
    path.closeSubpath();
    
    painter->setPen(Qt::NoPen);
    painter->setBrush(QColor("#FFA500"));
    painter->drawPath(path);

    QRect textRect = option.rect;
    textRect.setX(iconX + iconSize + iconPadding);
    
    QColor textColor = (option.state & QStyle::State_Selected) ? option.palette.highlightedText().color() : option.palette.text().color();
    painter->setPen(textColor);
    
    painter->drawText(textRect, Qt::AlignVCenter | Qt::AlignLeft, value);
    painter->restore();
}

QSize TicketPriorityDelegate::sizeHint(const QStyleOptionViewItem &option, const QModelIndex &index) const {
    QSize size = QStyledItemDelegate::sizeHint(option, index);
    int iconSize = 10;
    int iconPadding = 5;
    return QSize(size.width() + iconSize + iconPadding * 2, size.height());
}


void TicketModel::loadTickets(const QJsonArray& array) {
    beginResetModel();
    m_tickets.clear();
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
        t.createdAtRaw = obj.value("created_at").toString();
        t.createdAt = QDateTime::fromString(t.createdAtRaw, Qt::ISODateWithMs);
        t.updatedAt = QDateTime::fromString(obj.value("updated_at").toString(), Qt::ISODate);
        m_tickets.append(t);
    }
    endResetModel();
}

QMap<int, QString> TicketItem::statusLabels;
QMap<int, QString> TicketItem::priorityLabels;
QMap<int, QString> TicketItem::departmentNames;

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

void TicketModel::setTickets(const QVector<TicketItem> &tickets) {
    beginResetModel();
    m_tickets = tickets;
    endResetModel();
}