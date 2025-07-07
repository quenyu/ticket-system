#include "register_dialog.h"
#include <QVBoxLayout>
#include <QLabel>
#include <QMessageBox>
#include <QJsonDocument>
#include <QJsonArray>
#include <QJsonObject>

RegisterDialog::RegisterDialog(QWidget *parent)
    : QDialog(parent), apiClient(new APIClient(this)) {
    setWindowTitle("User Registration");
    setFixedSize(400, 350);
    QVBoxLayout *layout = new QVBoxLayout(this);
    layout->addWidget(new QLabel("Username:"));
    usernameEdit = new QLineEdit(this);
    layout->addWidget(usernameEdit);
    layout->addWidget(new QLabel("Email:"));
    emailEdit = new QLineEdit(this);
    layout->addWidget(emailEdit);
    layout->addWidget(new QLabel("Password:"));
    passwordEdit = new QLineEdit(this);
    passwordEdit->setEchoMode(QLineEdit::Password);
    layout->addWidget(passwordEdit);
    layout->addWidget(new QLabel("Department:"));
    departmentCombo = new QComboBox(this);
    layout->addWidget(departmentCombo);
    layout->addWidget(new QLabel("Role:"));
    roleCombo = new QComboBox(this);
    layout->addWidget(roleCombo);
    registerButton = new QPushButton("Register", this);
    layout->addWidget(registerButton);
    connect(registerButton, &QPushButton::clicked, this, &RegisterDialog::onRegisterClicked);
    connect(apiClient, &APIClient::departmentsReceived, this, &RegisterDialog::onDepartmentsReceived);
    connect(apiClient, &APIClient::rolesReceived, this, &RegisterDialog::onRolesReceived);
    connect(apiClient, &APIClient::registrationSuccess, this, &RegisterDialog::onRegistrationSuccess);
    connect(apiClient, &APIClient::registrationFailed, this, &RegisterDialog::onRegistrationFailed);
    loadDepartments();
    loadRoles();
}

void RegisterDialog::loadDepartments() {
    apiClient->getDepartments(""); // No token needed for public dictionaries
}

void RegisterDialog::loadRoles() {
    QJsonArray fakeRoles;
    QJsonObject role1;
    role1["id"] = "00000000-0000-0000-0000-000000000001";
    role1["name"] = "User";
    fakeRoles.append(role1);
    QJsonObject role2;
    role2["id"] = "00000000-0000-0000-0000-000000000002";
    role2["name"] = "Admin";
    fakeRoles.append(role2);
    QJsonDocument doc(fakeRoles);
    onRolesReceived(doc.toJson());
}

void RegisterDialog::onDepartmentsReceived(const QByteArray &data) {
    departmentCombo->clear();
    QJsonDocument doc = QJsonDocument::fromJson(data);
    if (doc.isArray()) {
        for (const QJsonValue &v : doc.array()) {
            QJsonObject o = v.toObject();
            int id = o.value("id").toInt(-1);
            QString name = o.value("name").toString();
            if (id > 0 && !name.isEmpty()) {
                departmentCombo->addItem(name, id);
            }
        }
    }
}

void RegisterDialog::onRolesReceived(const QByteArray &data) {
    roleCombo->clear();
    QJsonDocument doc = QJsonDocument::fromJson(data);
    if (doc.isArray()) {
        for (const QJsonValue &v : doc.array()) {
            QJsonObject o = v.toObject();
            QString id = o.value("id").toString();
            QString name = o.value("name").toString();
            if (!id.isEmpty() && !name.isEmpty()) {
                roleCombo->addItem(name, id);
            }
        }
    }
}

void RegisterDialog::onRegisterClicked() {
    QString username = usernameEdit->text().trimmed();
    QString email = emailEdit->text().trimmed();
    QString password = passwordEdit->text();
    int departmentId = departmentCombo->currentData().toInt();
    QString roleId = roleCombo->currentData().toString();
    if (username.isEmpty() || email.isEmpty() || password.isEmpty() || departmentId <= 0 || roleId.isEmpty()) {
        QMessageBox::warning(this, "Error", "Please fill in all fields.");
        return;
    }
    apiClient->registerUser(username, email, password, roleId, departmentId);
    registerButton->setEnabled(false);
}

void RegisterDialog::onRegistrationSuccess(const QString &userId) {
    QMessageBox::information(this, "Success", "User registered successfully! ID: " + userId);
    accept();
}

void RegisterDialog::onRegistrationFailed(const QString &error) {
    QMessageBox::warning(this, "Registration Error", error);
    registerButton->setEnabled(true);
} 