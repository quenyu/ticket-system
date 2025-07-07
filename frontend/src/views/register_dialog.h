#pragma once
#include <QDialog>
#include <QLineEdit>
#include <QComboBox>
#include <QPushButton>
#include "../network/api_client.h"

class RegisterDialog : public QDialog {
    Q_OBJECT
public:
    explicit RegisterDialog(QWidget *parent = nullptr);
private slots:
    void onRegisterClicked();
    void onDepartmentsReceived(const QByteArray &data);
    void onRolesReceived(const QByteArray &data);
    void onRegistrationSuccess(const QString &userId);
    void onRegistrationFailed(const QString &error);
private:
    QLineEdit *usernameEdit;
    QLineEdit *emailEdit;
    QLineEdit *passwordEdit;
    QComboBox *departmentCombo;
    QComboBox *roleCombo;
    QPushButton *registerButton;
    APIClient *apiClient;
    void loadDepartments();
    void loadRoles();
}; 