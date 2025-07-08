#pragma once

#include <QAbstractListModel>
#include <QVector>
#include <QString>
#include <QDateTime>
#include <QJsonObject>
#include <QJsonArray>

struct AttachmentItem {
    QString id;
    QString ticketId;
    QDateTime ticketCreatedAt;
    QString filename;
    QString filePath;
    QString uploadedBy;
    QDateTime uploadedAt;

    QJsonObject toJson() const;
    bool isImage() const {
        QString lower = filename.toLower();
        return lower.endsWith(".jpg") || lower.endsWith(".jpeg") || lower.endsWith(".png") || lower.endsWith(".gif");
    }
};

Q_DECLARE_METATYPE(AttachmentItem)

class AttachmentModel : public QAbstractListModel {
    Q_OBJECT
public:
    enum AttachmentRoles {
        IdRole = Qt::UserRole + 1,
        TicketIdRole,
        TicketCreatedAtRole,
        FilenameRole,
        FilePathRole,
        UploadedByRole,
        UploadedAtRole
    };

    explicit AttachmentModel(QObject *parent = nullptr);
    int rowCount(const QModelIndex &parent = QModelIndex()) const override;
    QVariant data(const QModelIndex &index, int role = Qt::DisplayRole) const override;
    QHash<int, QByteArray> roleNames() const override;
    void loadAttachments(const QJsonArray& array);
    void addAttachment(const AttachmentItem& attachment);
    AttachmentItem getAttachment(int row) const;
    void clearAttachments();
    void removeAttachment(int row);
private:
    QVector<AttachmentItem> m_attachments;
}; 