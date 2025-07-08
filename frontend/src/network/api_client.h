#pragma once
#include <QObject>
#include <QNetworkAccessManager>
#include <QNetworkReply>

class APIClient : public QObject {
    Q_OBJECT
public:
    explicit APIClient(QObject *parent = nullptr);
    void login(const QString &username, const QString &password);
    void getDepartments(const QString &token);
    void getRoles(const QString &token);
    void registerUser(const QString &username, const QString &email, const QString &password, const QString &roleId, int departmentId);
    void getUsers(const QString &token);
    void getStatuses(const QString &token);
    void getPriorities(const QString &token);
    void createTicket(const QString &token, const QJsonObject &ticketData);
    void updateTicket(const QString &token, const QString &ticketId, const QJsonObject &ticketData);
    void loadHistory(const QString &token, const QString &ticketId);
    void deleteAttachment(const QString &token, const QString &ticketId, const QString &attachmentId);
    void uploadAttachment(const QString &token, const QString &ticketId, const QString &filePath);
signals:
    void loginSuccess(const QString &token);
    void loginFailed(const QString &error);
    void departmentsReceived(const QByteArray &data);
    void rolesReceived(const QByteArray &data);
    void registrationSuccess(const QString &userId);
    void registrationFailed(const QString &error);
    void usersReceived(const QByteArray &data);
    void statusesReceived(const QByteArray &data);
    void prioritiesReceived(const QByteArray &data);
    void ticketCreated(const QByteArray &data);
    void apiError(const QString &error);
    void historyLoaded(const QByteArray &data);
    void attachmentDeleted(const QString &attachmentId);
    void attachmentUploaded();
private:
    QNetworkAccessManager manager;
}; 