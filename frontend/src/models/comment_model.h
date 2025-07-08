#pragma once

#include <QAbstractListModel>
#include <QVector>
#include <QString>
#include <QDateTime>
#include <QJsonObject>
#include <QJsonArray>
#include <QMap>

struct CommentItem {
    QString id;
    QString ticketId;
    QDateTime ticketCreatedAt;
    QString authorId;
    QString authorName;
    QString content;
    QDateTime createdAt;

    QJsonObject toJson() const;
};

Q_DECLARE_METATYPE(CommentItem)

class CommentModel : public QAbstractListModel {
    Q_OBJECT
public:
    enum CommentRoles {
        IdRole = Qt::UserRole + 1,
        TicketIdRole,
        TicketCreatedAtRole,
        AuthorIdRole,
        AuthorNameRole,
        ContentRole,
        CreatedAtRole
    };

    explicit CommentModel(QObject *parent = nullptr);
    int rowCount(const QModelIndex &parent = QModelIndex()) const override;
    QVariant data(const QModelIndex &index, int role = Qt::DisplayRole) const override;
    QVariant headerData(int section, Qt::Orientation orientation, int role = Qt::DisplayRole) const override;
    QHash<int, QByteArray> roleNames() const override;
    void loadComments(const QJsonArray& array);
    void addComment(const CommentItem& comment);
    CommentItem getComment(int row) const;
    void clearComments();
private:
    QVector<CommentItem> m_comments;
}; 