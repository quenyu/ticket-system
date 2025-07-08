#include "comment_model.h"
#include <QJsonDocument>
#include <QDebug>

QJsonObject CommentItem::toJson() const {
    QJsonObject obj;
    if (!id.isEmpty()) obj["comment_id"] = id;
    obj["ticket_id"] = ticketId;
    obj["ticket_created_at"] = ticketCreatedAt.toString(Qt::ISODate);
    obj["author_id"] = authorId;
    obj["content"] = content;
    obj["created_at"] = createdAt.toString(Qt::ISODate);
    return obj;
}

CommentModel::CommentModel(QObject *parent)
    : QAbstractListModel(parent) {}

int CommentModel::rowCount(const QModelIndex &parent) const {
    if (parent.isValid()) return 0;
    return m_comments.size();
}

QVariant CommentModel::data(const QModelIndex &index, int role) const {
    if (!index.isValid() || index.row() >= m_comments.size()) return QVariant();
    const CommentItem &comment = m_comments[index.row()];
    switch (role) {
        case IdRole: return comment.id;
        case TicketIdRole: return comment.ticketId;
        case TicketCreatedAtRole: return comment.ticketCreatedAt;
        case AuthorIdRole: return comment.authorId;
        case AuthorNameRole: return comment.authorName;
        case ContentRole: return comment.content;
        case CreatedAtRole: return comment.createdAt;
        case Qt::DisplayRole:
            return QString("[%1] %2: %3")
                   .arg(comment.createdAt.toString("hh:mm dd.MM.yyyy"))
                   .arg(comment.authorName)
                   .arg(comment.content);
        default: return QVariant();
    }
}

QVariant CommentModel::headerData(int section, Qt::Orientation orientation, int role) const {
    if (role == Qt::DisplayRole && orientation == Qt::Horizontal) {
    }
    return QVariant();
}

QHash<int, QByteArray> CommentModel::roleNames() const {
    QHash<int, QByteArray> roles;
    roles[IdRole] = "id";
    roles[TicketIdRole] = "ticketId";
    roles[TicketCreatedAtRole] = "ticketCreatedAt";
    roles[AuthorIdRole] = "authorId";
    roles[AuthorNameRole] = "authorName";
    roles[ContentRole] = "content";
    roles[CreatedAtRole] = "createdAt";
    return roles;
}

void CommentModel::loadComments(const QJsonArray& array) {
    beginResetModel();
    m_comments.clear();
    for (const QJsonValue &value : array) {
        QJsonObject obj = value.toObject();
        CommentItem c;
        c.id = obj.value("comment_id").toString();
        c.ticketId = obj.value("ticket_id").toString();
        c.ticketCreatedAt = QDateTime::fromString(obj.value("ticket_created_at").toString(), Qt::ISODate);
        c.authorId = obj.value("author_id").toString();
        c.authorName = obj.value("author_name").toString();
        c.content = obj.value("content").toString();
        c.createdAt = QDateTime::fromString(obj.value("created_at").toString(), Qt::ISODate);
        m_comments.append(c);
    }
    endResetModel();
}

void CommentModel::addComment(const CommentItem& comment) {
    CommentItem c = comment;
    if (c.authorId.isEmpty()) {
    }
    beginInsertRows(QModelIndex(), m_comments.size(), m_comments.size());
    m_comments.append(c);
    endInsertRows();
}

CommentItem CommentModel::getComment(int row) const {
    if (row < 0 || row >= m_comments.size()) {
        return CommentItem{};
    }
    return m_comments[row];
}

void CommentModel::clearComments() {
    beginResetModel();
    m_comments.clear();
    endResetModel();
}

void CommentModel::updateComment(int row, const CommentItem& updatedComment) {
    if (row < 0 || row >= m_comments.size()) return;
    m_comments[row] = updatedComment;
    emit dataChanged(index(row), index(row));
}

void CommentModel::removeComment(int row) {
    if (row < 0 || row >= m_comments.size()) return;
    beginRemoveRows(QModelIndex(), row, row);
    m_comments.removeAt(row);
    endRemoveRows();
} 