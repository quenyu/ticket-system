#include "login_dialog.h"
#include "../config.h"
#include <QLineEdit>
#include <QPushButton>
#include <QVBoxLayout>
#include <QLabel>
#include <QNetworkAccessManager>
#include <QNetworkRequest>
#include <QNetworkReply>
#include <QJsonDocument>
#include <QJsonObject>
#include <QMessageBox>
#include <QSpacerItem>
#include <QSizePolicy>
#include <QHBoxLayout>
#include <QIcon>
#include <QDebug>
#include "register_dialog.h"

LoginDialog::LoginDialog(QWidget *parent) : QDialog(parent) {
    setWindowTitle("Login");
    setFixedSize(360, 260);
    auto *layout = new QVBoxLayout(this);
    layout->setAlignment(Qt::AlignCenter);
    layout->setContentsMargins(32, 32, 32, 32);
    layout->setSpacing(12);

    auto *title = new QLabel("<b>Ticket System Login</b>", this);
    title->setAlignment(Qt::AlignHCenter);
    layout->addWidget(title);

    // Username row
    auto *userRow = new QHBoxLayout();
    auto *userIcon = new QLabel(this);
    userIcon->setPixmap(QIcon::fromTheme("user-identity").pixmap(20, 20));
    userRow->addWidget(userIcon);
    usernameEdit = new QLineEdit(this);
    usernameEdit->setPlaceholderText("Username");
    usernameEdit->setMinimumWidth(180);
    userRow->addWidget(usernameEdit);
    layout->addLayout(userRow);

    // Password row
    auto *passRow = new QHBoxLayout();
    auto *passIcon = new QLabel(this);
    passIcon->setPixmap(QIcon::fromTheme("dialog-password").pixmap(20, 20));
    passRow->addWidget(passIcon);
    passwordEdit = new QLineEdit(this);
    passwordEdit->setPlaceholderText("Password");
    passwordEdit->setEchoMode(QLineEdit::Password);
    passwordEdit->setMinimumWidth(180);
    passRow->addWidget(passwordEdit);
    // Show/hide password button
    auto *showPassBtn = new QPushButton(this);
    showPassBtn->setIcon(QIcon::fromTheme("view-password"));
    showPassBtn->setCheckable(true);
    showPassBtn->setToolTip("Show/Hide Password");
    passRow->addWidget(showPassBtn);
    layout->addLayout(passRow);
    connect(showPassBtn, &QPushButton::toggled, this, [this](bool checked) {
        passwordEdit->setEchoMode(checked ? QLineEdit::Normal : QLineEdit::Password);
    });

    // Login/Register row
    auto *btnRow = new QHBoxLayout();
    loginButton = new QPushButton("Login", this);
    loginButton->setMinimumWidth(100);
    btnRow->addWidget(loginButton);
    btnRow->addSpacing(12);
    registerButton = new QPushButton("Register", this);
    registerButton->setMinimumWidth(100);
    btnRow->addWidget(registerButton);
    layout->addLayout(btnRow);

    // Forgot password link (stub)
    auto *forgot = new QPushButton("Forgot password?", this);
    forgot->setFlat(true);
    forgot->setStyleSheet("color: #0078d7; text-align: left;");
    layout->addWidget(forgot, 0, Qt::AlignLeft);
    connect(forgot, &QPushButton::clicked, this, [](){ QMessageBox::information(nullptr, "Not implemented", "Password recovery is not implemented yet."); });

    setLayout(layout);
    networkManager = new QNetworkAccessManager(this);
    connect(loginButton, &QPushButton::clicked, this, &LoginDialog::onLoginClicked);
    connect(registerButton, &QPushButton::clicked, this, &LoginDialog::onRegisterClicked);
    // Enter key submits login
    connect(usernameEdit, &QLineEdit::returnPressed, this, &LoginDialog::onLoginClicked);
    connect(passwordEdit, &QLineEdit::returnPressed, this, &LoginDialog::onLoginClicked);
    usernameEdit->setFocus();
}

QString LoginDialog::username() const { return usernameEdit->text(); }
QString LoginDialog::password() const { return passwordEdit->text(); }
QString LoginDialog::jwtToken() const { return m_jwtToken; }

void LoginDialog::onLoginClicked() {
    QString username = usernameEdit->text().trimmed();
    QString password = passwordEdit->text();
    qDebug() << "Login attempt: username=" << username;
    if (username.isEmpty() || password.isEmpty()) {
        qDebug() << "Username or password is empty";
        QMessageBox::warning(this, "Error", "Please enter both username and password.");
        return;
    }
    QUrl url(Config::instance().fullApiUrl() + "/auth/login");
    qDebug() << "Login URL:" << url.toString();
    QNetworkRequest req(url);
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    QJsonObject obj;
    obj["username"] = username;
    obj["password"] = password;
    QJsonDocument doc(obj);
    QByteArray data = doc.toJson();
    qDebug() << "Sending login request with data:" << data;
    QNetworkReply *reply = networkManager->post(req, data);
    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        qDebug() << "Login reply finished";
        if (reply->error() != QNetworkReply::NoError) {
            qDebug() << "Network error:" << reply->errorString();
            QMessageBox::warning(this, "Error", "Network error: " + reply->errorString());
            reply->deleteLater();
            return;
        }
        QByteArray response = reply->readAll();
        qDebug() << "Login response data:" << response;
        QJsonDocument doc = QJsonDocument::fromJson(response);
        if (!doc.isObject()) {
            qDebug() << "Invalid JSON response";
            QMessageBox::warning(this, "Error", "Invalid server response.");
            reply->deleteLater();
            return;
        }
        QJsonObject obj = doc.object();
        if (obj.contains("error")) {
            QJsonObject err = obj["error"].toObject();
            qDebug() << "Server error:" << err;
            QString msg = err.value("message").toString();
            QMessageBox::warning(this, "Login failed", msg.isEmpty() ? "Login failed" : msg);
            reply->deleteLater();
            return;
        }
        if (!obj.contains("token")) {
            qDebug() << "No token in response";
            QMessageBox::warning(this, "Error", "No token in server response.");
            reply->deleteLater();
            return;
        }
        QString token = obj["token"].toString();
        qDebug() << "Login successful, token length:" << token.length();
        m_jwtToken = token;
        accept();
        reply->deleteLater();
    });
}

void LoginDialog::onRegisterClicked() {
    RegisterDialog dlg(this);
    if (dlg.exec() == QDialog::Accepted) {
        QMessageBox::information(this, "Success", "Registration successful! Now you can log in.");
    }
} 