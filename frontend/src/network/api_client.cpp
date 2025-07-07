#include "api_client.h"
#include "../config.h"
#include <QJsonDocument>
#include <QJsonObject>
#include <QNetworkRequest>
#include <QDebug>

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

void APIClient::getUsers(const QString &token) {
    QUrl url(Config::instance().fullApiUrl() + "/users");
    QNetworkRequest req(url);
    req.setRawHeader("Authorization", "Bearer " + token.toUtf8());
    QNetworkReply *reply = manager.get(req);
    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        if (reply->error() == QNetworkReply::NoError) {
            emit usersReceived(reply->readAll());
        } else {
            emit apiError(reply->errorString());
        }
        reply->deleteLater();
    });
}

void APIClient::getStatuses(const QString &token) {
    QUrl url(Config::instance().fullApiUrl() + "/ticket_statuses");
    QNetworkRequest req(url);
    req.setRawHeader("Authorization", "Bearer " + token.toUtf8());
    QNetworkReply *reply = manager.get(req);
    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        if (reply->error() == QNetworkReply::NoError) {
            emit statusesReceived(reply->readAll());
        } else {
            emit apiError(reply->errorString());
        }
        reply->deleteLater();
    });
}

void APIClient::getPriorities(const QString &token) {
    QUrl url(Config::instance().fullApiUrl() + "/ticket_priorities");
    QNetworkRequest req(url);
    req.setRawHeader("Authorization", "Bearer " + token.toUtf8());
    QNetworkReply *reply = manager.get(req);
    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        if (reply->error() == QNetworkReply::NoError) {
            emit prioritiesReceived(reply->readAll());
        } else {
            emit apiError(reply->errorString());
        }
        reply->deleteLater();
    });
}

void APIClient::createTicket(const QString &token, const QJsonObject &ticketData) {
    QUrl url(Config::instance().fullApiUrl() + "/tickets");
    QNetworkRequest req(url);
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    req.setRawHeader("Authorization", "Bearer " + token.toUtf8());
    QNetworkReply *reply = manager.post(req, QJsonDocument(ticketData).toJson());
    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        if (reply->error() == QNetworkReply::NoError) {
            emit ticketCreated(reply->readAll());
        } else {
            emit apiError(reply->errorString());
        }
        reply->deleteLater();
    });
}

void APIClient::getRoles(const QString &token) {
    QUrl url(Config::instance().fullApiUrl() + "/roles");
    QNetworkRequest req(url);
    req.setRawHeader("Authorization", "Bearer " + token.toUtf8());
    QNetworkReply *reply = manager.get(req);
    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        if (reply->error() == QNetworkReply::NoError) {
            emit rolesReceived(reply->readAll());
        } else {
            emit apiError(reply->errorString());
        }
        reply->deleteLater();
    });
}

void APIClient::registerUser(const QString &username, const QString &email, const QString &password, const QString &roleId, int departmentId) {
    QUrl url(Config::instance().fullApiUrl() + "/auth/register");
    QNetworkRequest req(url);
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    QJsonObject obj;
    obj["username"] = username;
    obj["email"] = email;
    obj["password"] = password;
    obj["role_id"] = roleId;
    obj["department_id"] = departmentId;
    QNetworkReply *reply = manager.post(req, QJsonDocument(obj).toJson());
    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        if (reply->error() == QNetworkReply::NoError) {
            QJsonDocument doc = QJsonDocument::fromJson(reply->readAll());
            if (doc.isObject() && doc.object().contains("user_id")) {
                emit registrationSuccess(doc.object()["user_id"].toString());
            } else {
                emit registrationFailed("Invalid response format");
            }
        } else {
            emit registrationFailed(reply->errorString());
        }
        reply->deleteLater();
    });
}

void APIClient::updateTicket(const QString &token, const QString &ticketId, const QJsonObject &ticketData) {
    QUrl url(Config::instance().fullApiUrl() + "/tickets/" + ticketId);
    QNetworkRequest req(url);
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    req.setRawHeader("Authorization", "Bearer " + token.toUtf8());
    QNetworkReply *reply = manager.sendCustomRequest(req, "PATCH", QJsonDocument(ticketData).toJson());
    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        if (reply->error() == QNetworkReply::NoError) {
            emit ticketCreated(reply->readAll());
        } else {
            emit apiError(reply->errorString());
        }
        reply->deleteLater();
    });
}

void APIClient::loadHistory(const QString &token, const QString &ticketId) {
    QUrl url(Config::instance().fullApiUrl() + "/tickets/" + ticketId + "/history");
    QNetworkRequest req(url);
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    req.setRawHeader("Authorization", "Bearer " + token.toUtf8());
    qDebug() << "JWT for history:" << token;
    QNetworkReply *reply = manager.get(req);
    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        if (reply->error() == QNetworkReply::NoError) {
            emit historyLoaded(reply->readAll());
        } else {
            emit apiError(reply->errorString());
        }
        reply->deleteLater();
    });
} 