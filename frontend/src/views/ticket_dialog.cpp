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
        setFixedSize(400, 300);
        
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
        descEdit->setMaximumHeight(100);
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
        
        qDebug() << "Creating save button...";
        saveButton = new QPushButton(mode == Create ? "Create" : "Save", this);
        saveButton->setEnabled(false);
        mainLayout->addWidget(saveButton);
        
        qDebug() << "Connecting signals...";
        connect(saveButton, &QPushButton::clicked, this, &TicketDialog::onSaveClicked);
        connect(titleEdit, &QLineEdit::returnPressed, this, &TicketDialog::onSaveClicked);
        
        qDebug() << "Creating network manager...";
        network = new QNetworkAccessManager(this);
        
        qDebug() << "Setting focus...";
        titleEdit->setFocus();
        
        // Отложенная загрузка департаментов
        QTimer::singleShot(200, this, &TicketDialog::loadDepartments);
        // Отложенная загрузка статусов
        QTimer::singleShot(150, this, [this]{ loadStatuses(); });
        // Отложенная загрузка приоритетов
        QTimer::singleShot(180, this, [this]{ loadPriorities(); });
        
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

void TicketDialog::onSaveClicked() {
    qDebug() << "=== onSaveClicked() START ===";
    
    // Валидация
    if (titleEdit->text().trimmed().isEmpty()) {
        QMessageBox::warning(this, "Ошибка", "Заголовок не может быть пустым");
        titleEdit->setFocus();
        return;
    }
    
    int deptId = departmentCombo->currentData().toInt();
    int statusId = statusCombo->currentData().toInt();
    int priorityId = priorityCombo->currentData().toInt();
    
    if (deptId <= 0 || statusId <= 0 || priorityId <= 0) {
        QMessageBox::warning(this, "Ошибка", "Пожалуйста, выберите все обязательные поля");
        return;
    }
    
    qDebug() << "Validation passed";
    qDebug() << "Department ID:" << deptId;
    qDebug() << "Status ID:" << statusId;
    qDebug() << "Priority ID:" << priorityId;
    
    // Формируем JSON
    QJsonObject obj;
    obj["title"] = titleEdit->text().trimmed();
    obj["description"] = descEdit->toPlainText().trimmed();
    obj["department_id"] = deptId;
    obj["status_id"] = statusId;
    obj["priority_id"] = priorityId;
    
    qDebug() << "JSON object created:" << QJsonDocument(obj).toJson();
    
    // Отправляем запрос
    QNetworkRequest req;
    req.setHeader(QNetworkRequest::ContentTypeHeader, "application/json");
    req.setRawHeader("Authorization", "Bearer " + m_jwtToken.toUtf8());
    
    QNetworkReply *reply = nullptr;
    
    if (m_mode == Create) {
        req.setUrl(QUrl(Config::instance().fullApiUrl() + "/tickets"));
        reply = network->post(req, QJsonDocument(obj).toJson());
        qDebug() << "Sending POST request to create ticket";
    } else {
        req.setUrl(QUrl(Config::instance().fullApiUrl() + "/tickets/" + m_ticket.id));
        reply = network->sendCustomRequest(req, "PATCH", QJsonDocument(obj).toJson());
        qDebug() << "Sending PATCH request to update ticket";
    }
    connect(reply, &QNetworkReply::finished, this, [this, reply]() { onReplyFinished(reply); });
}

void TicketDialog::onReplyFinished(QNetworkReply* reply) {
    if (reply->error() == QNetworkReply::NoError) {
        QByteArray data = reply->readAll();
        QJsonDocument doc = QJsonDocument::fromJson(data);
        
        // Проверяем валидность ответа
        if (doc.isNull()) {
            QMessageBox::warning(this, "API Error", "Получен невалидный JSON ответ");
            reply->deleteLater();
            return;
        }
        
        emit ticketSaved();
        accept();
    } else {
        QByteArray data = reply->readAll();
        QJsonDocument doc = QJsonDocument::fromJson(data);
        QString msg = reply->errorString();
        
        // Пытаемся извлечь сообщение об ошибке из JSON
        if (!doc.isNull() && doc.isObject()) {
            QJsonObject obj = doc.object();
            if (obj.contains("error") && obj["error"].isObject()) {
                QJsonObject errorObj = obj["error"].toObject();
                if (errorObj.contains("message")) {
                    msg = errorObj["message"].toString();
                }
            }
        }
        
        QMessageBox::warning(this, "API Error", msg);
    }
    reply->deleteLater();
}

void TicketDialog::setCurrentTab(int index) {
    if (tabs && index >= 0 && index < tabs->count()) {
        tabs->setCurrentIndex(index);
    }
} 