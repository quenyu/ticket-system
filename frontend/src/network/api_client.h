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
signals:
    void loginSuccess(const QString &token);
    void loginFailed(const QString &error);
    void departmentsReceived(const QByteArray &data);
    void apiError(const QString &error);
private:
    QNetworkAccessManager manager;
}; 