#include "attachment_model.h"
#include <QJsonDocument>
#include <QDebug>

QJsonObject AttachmentItem::toJson() const {
    QJsonObject obj;
    if (!id.isEmpty()) obj["attachment_id"] = id;
    obj["ticket_id"] = ticketId;
    obj["ticket_created_at"] = ticketCreatedAt.toString(Qt::ISODate);
    obj["filename"] = filename;
    obj["file_path"] = filePath;
    obj["uploaded_by"] = uploadedBy;
    obj["uploaded_at"] = uploadedAt.toString(Qt::ISODate);
    return obj;
}

AttachmentModel::AttachmentModel(QObject *parent)
    : QAbstractListModel(parent) {}

int AttachmentModel::rowCount(const QModelIndex &parent) const {
    if (parent.isValid()) return 0;
    return m_attachments.size();
}

QVariant AttachmentModel::data(const QModelIndex &index, int role) const {
    if (!index.isValid() || index.row() >= m_attachments.size()) return QVariant();
    const AttachmentItem &att = m_attachments[index.row()];
    switch (role) {
        case IdRole: return att.id;
        case TicketIdRole: return att.ticketId;
        case TicketCreatedAtRole: return att.ticketCreatedAt;
        case FilenameRole: return att.filename;
        case FilePathRole: return att.filePath;
        case UploadedByRole: return att.uploadedBy;
        case UploadedAtRole: return att.uploadedAt;
        case Qt::DisplayRole:
            return QString("[%1] %2 (%3)")
                   .arg(att.uploadedAt.toString("hh:mm dd.MM.yyyy"))
                   .arg(att.filename)
                   .arg(att.filePath);
        default: return QVariant();
    }
}

QHash<int, QByteArray> AttachmentModel::roleNames() const {
    QHash<int, QByteArray> roles;
    roles[IdRole] = "id";
    roles[TicketIdRole] = "ticketId";
    roles[TicketCreatedAtRole] = "ticketCreatedAt";
    roles[FilenameRole] = "filename";
    roles[FilePathRole] = "filePath";
    roles[UploadedByRole] = "uploadedBy";
    roles[UploadedAtRole] = "uploadedAt";
    return roles;
}

void AttachmentModel::loadAttachments(const QJsonArray& array) {
    beginResetModel();
    m_attachments.clear();
    for (const QJsonValue &value : array) {
        QJsonObject obj = value.toObject();
        AttachmentItem a;
        a.id = obj.value("attachment_id").toString();
        a.ticketId = obj.value("ticket_id").toString();
        a.ticketCreatedAt = QDateTime::fromString(obj.value("ticket_created_at").toString(), Qt::ISODate);
        a.filename = obj.value("filename").toString();
        a.filePath = obj.value("file_path").toString();
        a.uploadedBy = obj.value("uploaded_by").toString();
        a.uploadedAt = QDateTime::fromString(obj.value("uploaded_at").toString(), Qt::ISODate);
        m_attachments.append(a);
    }
    endResetModel();
}

void AttachmentModel::addAttachment(const AttachmentItem& attachment) {
    beginInsertRows(QModelIndex(), m_attachments.size(), m_attachments.size());
    m_attachments.append(attachment);
    endInsertRows();
}

AttachmentItem AttachmentModel::getAttachment(int row) const {
    if (row < 0 || row >= m_attachments.size()) {
        return AttachmentItem{};
    }
    return m_attachments[row];
}

void AttachmentModel::clearAttachments() {
    beginResetModel();
    m_attachments.clear();
    endResetModel();
}

void AttachmentModel::removeAttachment(int row) {
    if (row < 0 || row >= m_attachments.size()) return;
    beginRemoveRows(QModelIndex(), row, row);
    m_attachments.removeAt(row);
    endRemoveRows();
} 