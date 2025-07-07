#include <QHBoxLayout>
#include <QFont>
#include "ticket_dialog.h"
#include <QTabWidget>
#include <QVBoxLayout>
#include <QLineEdit>
#include <QTextEdit>
#include <QTableView>
#include <QLabel>
#include <QPushButton>
#include <QMessageBox>
#include <QJsonDocument>
#include <QJsonObject>
#include <QNetworkRequest>
#include <QNetworkAccessManager>
#include <QNetworkReply>
#include <QComboBox>
#include <QJsonArray>
#include <QDebug>
#include <QTimer>
#include "../config.h"
#include "../network/api_client.h"

TicketDialog::TicketDialog(const TicketItem &ticket, const QString &jwtToken, QWidget *parent, Mode mode)
    : QDialog(parent), m_ticket(ticket), m_jwtToken(jwtToken), m_mode(mode) {
    qDebug() << "=== TicketDialog constructor START ===";
    qDebug() << "Mode:" << (mode == Create ? "Create" : "Edit");
    qDebug() << "Ticket ID:" << ticket.id;
    qDebug() << "JWT token length:" << jwtToken.length();
    
    try {
        qDebug() << "Setting window title...";
        setWindowTitle(mode == Create ? "Create Ticket" : "Edit Ticket");
        qDebug() << "Setting fixed size...";
        setFixedSize(500, 420);
        
        qDebug() << "Creating main layout...";
        QVBoxLayout *mainLayout = new QVBoxLayout(this);
        
        qDebug() << "Creating title label...";
        auto *titleLbl = new QLabel(mode == Create ? "<b>Create New Ticket</b>" : "<b>Edit Ticket</b>", this);
        titleLbl->setAlignment(Qt::AlignHCenter);
        mainLayout->addWidget(titleLbl);
        
        qDebug() << "Creating title edit...";
        titleEdit = new QLineEdit(ticket.title, this);
        titleEdit->setPlaceholderText("Title");
        mainLayout->addWidget(titleEdit);
        
        qDebug() << "Creating description edit...";
        descEdit = new QTextEdit("", this);
        descEdit->setPlaceholderText("Description");
        descEdit->setMinimumHeight(80);
        descEdit->setMaximumHeight(160);
        mainLayout->addWidget(descEdit);
        
        qDebug() << "Creating department combo...";
        departmentCombo = new QComboBox(this);
        departmentCombo->addItem("Loading departments...", -1);
        mainLayout->addWidget(new QLabel("Department", this));
        mainLayout->addWidget(departmentCombo);
        
        qDebug() << "Creating status combo...";
        statusCombo = new QComboBox(this);
        statusCombo->addItem("Loading statuses...", -1);
        mainLayout->addWidget(new QLabel("Status", this));
        mainLayout->addWidget(statusCombo);
        
        qDebug() << "Creating priority combo...";
        priorityCombo = new QComboBox(this);
        priorityCombo->addItem("Loading priorities...", -1);
        mainLayout->addWidget(new QLabel("Priority", this));
        mainLayout->addWidget(priorityCombo);
        
        qDebug() << "Creating assignee combo...";
        assigneeCombo = new QComboBox(this);
        assigneeCombo->addItem("Loading assignees...", "");
        mainLayout->addWidget(new QLabel("Assignee", this));
        mainLayout->addWidget(assigneeCombo);
        
        qDebug() << "Creating save button...";
        saveButton = new QPushButton(mode == Create ? "Create" : "Save", this);
        saveButton->setEnabled(true);
        mainLayout->addWidget(saveButton);
        
        qDebug() << "Connecting signals...";
        connect(saveButton, &QPushButton::clicked, this, &TicketDialog::onSaveClicked);
        connect(titleEdit, &QLineEdit::returnPressed, this, &TicketDialog::onSaveClicked);
        
        qDebug() << "Creating network manager...";
        network = new QNetworkAccessManager(this);
        
        qDebug() << "Setting focus...";
        titleEdit->setFocus();
        
        // Deferred loading of departments
        QTimer::singleShot(200, this, &TicketDialog::loadDepartments);
        // Deferred loading of statuses
        QTimer::singleShot(150, this, [this]{ loadStatuses(); });
        // Deferred loading of priorities
        QTimer::singleShot(180, this, [this]{ loadPriorities(); });
        
        connect(departmentCombo, QOverload<int>::of(&QComboBox::currentIndexChanged), this, [this]{
            int deptId = departmentCombo->currentData().toInt();
            filterAssigneesByDepartment(deptId);
        });
        QTimer::singleShot(220, this, [this]{ loadUsers(); });
        
        qDebug() << "=== TicketDialog constructor SUCCESS ===";
    } catch (const std::exception& e) {
        qDebug() << "EXCEPTION in TicketDialog constructor:" << e.what();
        throw;
    } catch (...) {
        qDebug() << "UNKNOWN EXCEPTION in TicketDialog constructor";
        throw;
    }
}

void TicketDialog::loadDepartments() {
    qDebug() << "=== loadDepartments() START ===";
    qDebug() << "Network manager:" << (network ? "valid" : "null");
    
    if (!network) {
        qDebug() << "ERROR: Network manager is null!";
        departmentCombo->clear();
        departmentCombo->addItem("Network manager not initialized", -1);
        saveButton->setEnabled(false);
        return;
    }
    
    qDebug() << "Creating network request...";
    
    try {
        QUrl url(Config::instance().fullApiUrl() + "/departments");
        qDebug() << "URL:" << url.toString();
        
        QNetworkRequest req(url);
        req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
        qDebug() << "Network request created";
        
        qDebug() << "Sending GET request...";
        QNetworkReply *reply = nullptr;
        
        try {
            reply = network->get(req);
        } catch (const std::exception& e) {
            qDebug() << "EXCEPTION during network->get():" << e.what();
            departmentCombo->clear();
            departmentCombo->addItem("Network exception", -1);
            saveButton->setEnabled(false);
            return;
        }
        
        qDebug() << "GET request sent, reply:" << (reply ? "valid" : "null");
        
        if (!reply) {
            qDebug() << "ERROR: reply is null!";
            departmentCombo->clear();
            departmentCombo->addItem("Network error", -1);
            saveButton->setEnabled(false);
            return;
        }
        
        qDebug() << "Connecting finished signal...";
        connect(reply, &QNetworkReply::finished, this, [this, reply]() {
            qDebug() << "Department reply finished";
            if (reply->error() == QNetworkReply::NoError) {
                QByteArray data = reply->readAll();
                qDebug() << "Department data received:" << data.length() << "bytes";
                
                QJsonDocument doc = QJsonDocument::fromJson(data);
                if (doc.isNull()) {
                    qDebug() << "ERROR: Invalid JSON in department response";
                    departmentCombo->clear();
                    departmentCombo->addItem("Invalid JSON response", -1);
                    saveButton->setEnabled(false);
                } else if (!doc.isArray()) {
                    qDebug() << "ERROR: Department response is not an array";
                    departmentCombo->clear();
                    departmentCombo->addItem("Invalid response format", -1);
                    saveButton->setEnabled(false);
                } else {
                    QJsonArray arr = doc.array();
                    qDebug() << "Department array size:" << arr.size();
                    
                    departmentCombo->clear();
                    bool hasValid = false;
                    
                    for (const QJsonValue &v : arr) {
                        if (!v.isObject()) {
                            qDebug() << "WARNING: Department item is not an object";
                            continue;
                        }
                        
                        QJsonObject o = v.toObject();
                        if (!o.contains("id") || !o.contains("name")) {
                            qDebug() << "WARNING: Department object missing required fields";
                            continue;
                        }
                        
                        int id = o["id"].toInt();
                        QString name = o["name"].toString();
                        
                        if (id > 0 && !name.isEmpty()) {
                            hasValid = true;
                            departmentCombo->addItem(name, id);
                            qDebug() << "Added department:" << name << "with ID:" << id;
                        }
                    }
                    
                    if (hasValid) {
                        int idx = 0;
                        for (int i = 0; i < departmentCombo->count(); ++i) {
                            if (departmentCombo->itemData(i).toInt() > 0) { 
                                idx = i; 
                                break; 
                            }
                        }
                        departmentCombo->setCurrentIndex(idx);
                        qDebug() << "Set department combo to index:" << idx;
                    } else {
                        departmentCombo->addItem("No departments available", -1);
                        qDebug() << "No valid departments found";
                    }
                }
            } else {
                qDebug() << "Network error in department request:" << reply->errorString();
                departmentCombo->clear();
                departmentCombo->addItem("Failed to load departments", -1);
            }
            reply->deleteLater();
        });
        
        qDebug() << "=== loadDepartments() END ===";
    } catch (const std::exception& e) {
        qDebug() << "EXCEPTION in loadDepartments:" << e.what();
        departmentCombo->clear();
        departmentCombo->addItem("Exception occurred", -1);
        saveButton->setEnabled(false);
    } catch (...) {
        qDebug() << "UNKNOWN EXCEPTION in loadDepartments";
        departmentCombo->clear();
        departmentCombo->addItem("Unknown error", -1);
        saveButton->setEnabled(false);
    }
}

void TicketDialog::loadStatuses() {
    qDebug() << "=== loadStatuses() START ===";
    if (!network) return;
    QNetworkRequest req(QUrl(Config::instance().fullApiUrl() + "/ticket_statuses"));
    QNetworkReply *reply = network->get(req);
    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        if (reply->error() == QNetworkReply::NoError) {
            QByteArray data = reply->readAll();
            QJsonDocument doc = QJsonDocument::fromJson(data);
            if (doc.isNull() || !doc.isArray()) {
                qDebug() << "ERROR: Invalid JSON in status response";
                statusCombo->clear();
                statusCombo->addItem("Invalid response format", -1);
                return;
            }
            
            QJsonArray arr = doc.array();
            statusCombo->clear();
            bool hasValid = false;
            for (const QJsonValue &v : arr) {
                if (!v.isObject()) continue;
                QJsonObject o = v.toObject();
                if (!o.contains("id") || !o.contains("label")) continue;
                
                int id = o["id"].toInt();
                QString name = o["label"].toString();
                if (id > 0 && !name.isEmpty()) {
                    hasValid = true;
                    statusCombo->addItem(name, id);
                }
            }
            if (hasValid) {
                int idx = 0;
                for (int i = 0; i < statusCombo->count(); ++i) {
                    if (statusCombo->itemData(i).toInt() > 0) { idx = i; break; }
                }
                statusCombo->setCurrentIndex(idx);
            } else {
                statusCombo->addItem("No statuses available", -1);
            }
        } else {
            statusCombo->clear();
            statusCombo->addItem("Failed to load statuses", -1);
        }
        reply->deleteLater();
    });
    qDebug() << "=== loadStatuses() END ===";
}

void TicketDialog::loadPriorities() {
    qDebug() << "=== loadPriorities() START ===";
    if (!network) return;
    QNetworkRequest req(QUrl(Config::instance().fullApiUrl() + "/ticket_priorities"));
    QNetworkReply *reply = network->get(req);
    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        if (reply->error() == QNetworkReply::NoError) {
            QByteArray data = reply->readAll();
            QJsonDocument doc = QJsonDocument::fromJson(data);
            if (doc.isNull() || !doc.isArray()) {
                qDebug() << "ERROR: Invalid JSON in priority response";
                priorityCombo->clear();
                priorityCombo->addItem("Invalid response format", -1);
                return;
            }
            
            QJsonArray arr = doc.array();
            priorityCombo->clear();
            bool hasValid = false;
            for (const QJsonValue &v : arr) {
                if (!v.isObject()) continue;
                QJsonObject o = v.toObject();
                if (!o.contains("id") || !o.contains("label")) continue;
                
                int id = o["id"].toInt();
                QString name = o["label"].toString();
                if (id > 0 && !name.isEmpty()) {
                    hasValid = true;
                    priorityCombo->addItem(name, id);
                }
            }
            if (hasValid) {
                int idx = 0;
                for (int i = 0; i < priorityCombo->count(); ++i) {
                    if (priorityCombo->itemData(i).toInt() > 0) { idx = i; break; }
                }
                priorityCombo->setCurrentIndex(idx);
            } else {
                priorityCombo->addItem("No priorities available", -1);
            }
        } else {
            priorityCombo->clear();
            priorityCombo->addItem("Failed to load priorities", -1);
        }
        reply->deleteLater();
    });
    qDebug() << "=== loadPriorities() END ===";
}

void TicketDialog::loadUsers() {
    if (!network) return;
    QNetworkRequest req(QUrl(Config::instance().fullApiUrl() + "/users"));
    QNetworkReply *reply = network->get(req);
    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        if (reply->error() == QNetworkReply::NoError) {
            QByteArray data = reply->readAll();
            QJsonDocument doc = QJsonDocument::fromJson(data);
            if (doc.isNull() || !doc.isArray()) {
                assigneeCombo->clear();
                assigneeCombo->addItem("Invalid user response", "");
                return;
            }
            QJsonArray arr = doc.array();
            users.clear();
            assigneeCombo->clear();
            for (const QJsonValue &v : arr) {
                if (!v.isObject()) continue;
                QJsonObject o = v.toObject();
                QString userId = o.value("user_id").toString();
                if (userId.isEmpty()) userId = o.value("id").toString();
                QString username = o.value("username").toString();
                int deptId = o.value("department_id").toInt(-1);
                if (!userId.isEmpty() && !username.isEmpty() && deptId > 0) {
                    users.append({username, userId, deptId});
                }
            }
            qDebug() << "Loaded users:";
            for (const auto &u : users) {
                qDebug() << u.username << u.userId << u.departmentId;
            }
            if (users.isEmpty())
                assigneeCombo->addItem("No users available", "");
            // Show all users, no filtering
            for (const auto &user : users) {
                assigneeCombo->addItem(user.username, user.userId);
            }
            if (!users.isEmpty()) {
                assigneeCombo->setCurrentIndex(0);
                saveButton->setEnabled(true);
            } else {
                saveButton->setEnabled(false);
            }
        } else {
            assigneeCombo->clear();
            assigneeCombo->addItem("Failed to load users", "");
        }
        reply->deleteLater();
    });
}

void TicketDialog::filterAssigneesByDepartment(int departmentId) {
    assigneeCombo->clear();
    QVector<int> validIndexes;
    for (int i = 0; i < users.size(); ++i) {
        const auto &user = users[i];
        if (departmentId <= 0 || user.departmentId == departmentId) {
            assigneeCombo->addItem(user.username, user.userId);
            validIndexes.append(i);
        }
    }
    if (validIndexes.isEmpty() && !users.isEmpty()) {
        const auto &user = users[0];
        assigneeCombo->addItem(user.username, user.userId);
        assigneeCombo->setCurrentIndex(0);
        saveButton->setEnabled(true);
    } else if (validIndexes.isEmpty()) {
        assigneeCombo->addItem("No assignees for department", "");
        assigneeCombo->setCurrentIndex(0);
        saveButton->setEnabled(false);
    } else {
        assigneeCombo->setCurrentIndex(0);
        saveButton->setEnabled(true);
    }
}

void TicketDialog::onSaveClicked() {
    qDebug() << "=== onSaveClicked() START ===";
    
    // Validation
    if (titleEdit->text().trimmed().isEmpty()) {
        QMessageBox::warning(this, "Error", "Title cannot be empty");
        titleEdit->setFocus();
        return;
    }
    
    int deptId = departmentCombo->currentData().toInt();
    int statusId = statusCombo->currentData().toInt();
    int priorityId = priorityCombo->currentData().toInt();
    
    if (deptId <= 0 || statusId <= 0 || priorityId <= 0) {
        QMessageBox::warning(this, "Error", "Please select all required fields");
        return;
    }
    
    QString assigneeId = assigneeCombo->currentData().toString();
    qDebug() << "Assignee combo count:" << assigneeCombo->count();
    qDebug() << "Current index:" << assigneeCombo->currentIndex();
    qDebug() << "Current text:" << assigneeCombo->currentText();
    qDebug() << "Current data:" << assigneeId;
    if (assigneeId.isEmpty()) {
        QMessageBox::warning(this, "Error", "No available assignees");
        return;
    }
    
    qDebug() << "Validation passed";
    qDebug() << "Department ID:" << deptId;
    qDebug() << "Status ID:" << statusId;
    qDebug() << "Priority ID:" << priorityId;
    
    // Form JSON
    QJsonObject obj;
    obj["title"] = titleEdit->text().trimmed();
    obj["description"] = descEdit->toPlainText().trimmed();
    obj["department_id"] = deptId;
    obj["status_id"] = statusId;
    obj["priority_id"] = priorityId;
    obj["assignee_id"] = assigneeId;
    
    qDebug() << "JSON object created:" << QJsonDocument(obj).toJson();
    
    // Use APIClient to create ticket
    APIClient *api = new APIClient(this);
    connect(api, &APIClient::ticketCreated, this, [this](const QByteArray &data){
        emit ticketSaved();
        accept();
    });
    connect(api, &APIClient::apiError, this, [this](const QString &err){
        QMessageBox::warning(this, "Error", err);
    });
    api->createTicket(m_jwtToken, obj);
}

void TicketDialog::setCurrentTab(int index) {
    if (tabs && index >= 0 && index < tabs->count()) {
        tabs->setCurrentIndex(index);
    }
} 