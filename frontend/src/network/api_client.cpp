#include "api_client.h"
#include "../config.h"
#include <QJsonDocument>
#include <QJsonObject>
#include <QNetworkRequest>

APIClient::APIClient(QObject *parent) : QObject(parent) {}

void APIClient::login(const QString &username, const QString &password) {
    QUrl url(Config::instance().fullApiUrl() + "/auth/login");
    QNetworkRequest req(url);
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    QJsonObject obj{{"username", username}, {"password", password}};
    QNetworkReply *reply = manager.post(req, QJsonDocument(obj).toJson());
    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        if (reply->error() == QNetworkReply::NoError) {
            QByteArray data = reply->readAll();
            QJsonDocument doc = QJsonDocument::fromJson(data);
            if (doc.isObject() && doc.object().contains("token")) {
                emit loginSuccess(doc.object()["token"].toString());
            } else {
                emit loginFailed("Invalid response format");
            }
        } else {
            emit loginFailed(reply->errorString());
        }
        reply->deleteLater();
    });
}

void APIClient::getDepartments(const QString &token) {
    QUrl url(Config::instance().fullApiUrl() + "/departments");
    QNetworkRequest req(url);
    req.setRawHeader("Authorization", "Bearer " + token.toUtf8());
    QNetworkReply *reply = manager.get(req);
    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        if (reply->error() == QNetworkReply::NoError) {
            emit departmentsReceived(reply->readAll());
        } else {
            emit apiError(reply->errorString());
        }
        reply->deleteLater();
    });
} 