#pragma once
#include <QDialog>
#include <QNetworkAccessManager>

QT_BEGIN_NAMESPACE
class QLineEdit;
class QPushButton;
QT_END_NAMESPACE

class LoginDialog : public QDialog {
    Q_OBJECT
public:
    explicit LoginDialog(QWidget *parent = nullptr);
    QString username() const;
    QString password() const;
    QString jwtToken() const;
private slots:
    void onLoginClicked();
    void onRegisterClicked();
    void onRegisterReply(QNetworkReply* reply);
signals:
    void loginSuccess(const QString &token);
private:
    QLineEdit *usernameEdit;
    QLineEdit *passwordEdit;
    QPushButton *loginButton;
    QPushButton *registerButton;
    QNetworkAccessManager* networkManager;
    QString m_jwtToken;
}; 