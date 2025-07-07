#include "mainwindow.h"
#include "ticket_model.h"
#include "views/ticket_dialog.h"
#include "config.h"

#include <QHeaderView>
#include <QMenuBar>
#include <QToolBar>
#include <QTableView>
#include <QHBoxLayout>
#include <QLineEdit>
#include <QPushButton>
#include <QStandardItemModel>
#include <QStandardItem>
#include <QStatusBar>
#include <QNetworkAccessManager>
#include <QNetworkRequest>
#include <QNetworkReply>
#include <QUrlQuery>
#include <QMessageBox>
#include <QJsonDocument>
#include <QJsonArray>
#include <QJsonObject>
#include <QDebug>

TicketModel *ticketModel = nullptr;
QNetworkAccessManager *network = nullptr;

MainWindow::MainWindow(const QString &jwt, QWidget *parent)
    : QMainWindow(parent), jwtToken(jwt)
{
    // --- Title and size ---
    setWindowTitle("Ticket System - All Tickets");
    resize(1200, 800);

    // --- JWT → userId from payload ---
    QStringList parts = jwtToken.split('.');
    if (parts.size() >= 2) {
        QByteArray payload = QByteArray::fromBase64(parts[1].toUtf8());
        QJsonDocument doc = QJsonDocument::fromJson(payload);
        if (doc.isObject()) {
            userId = doc.object().value("user_id").toString();
        }
    }
    qDebug() << "Extracted userId =" << userId;

    apiBaseUrl = Config::instance().fullApiUrl();

    // --- Menu and toolbar ---
    menuBar = new QMenuBar(this);
    setMenuBar(menuBar);
    menuBar->addMenu("File");
    menuBar->addMenu("Help");

    toolBar = new QToolBar(this);
    addToolBar(toolBar);
    addAction     = toolBar->addAction("➕ Add");
    editAction    = toolBar->addAction("✏️ Edit");
    deleteAction  = toolBar->addAction("🗑️ Delete");
    refreshAction = toolBar->addAction("🔄 Refresh");

    // Search
    QLineEdit *searchEdit = new QLineEdit(this);
    searchEdit->setPlaceholderText("Search tickets...");
    searchEdit->setMinimumWidth(200);
    toolBar->addWidget(searchEdit);

    // --- Table ---
    ticketModel = new TicketModel(this);
    tableView = new QTableView(this);
    tableView->setModel(ticketModel);
    tableView->setSelectionBehavior(QAbstractItemView::SelectRows);
    tableView->setSelectionMode(QAbstractItemView::SingleSelection);
    tableView->horizontalHeader()->setStretchLastSection(true);
    tableView->setAlternatingRowColors(true);
    setCentralWidget(tableView);

    // Set custom column widths
    tableView->setColumnWidth(2, 300); // Description column wider
    tableView->setColumnWidth(8, 110); // Updated At column narrower

    // --- Status bar ---
    m_statusBar = new QStatusBar(this);
    setStatusBar(m_statusBar);
    m_statusBar->showMessage("Ready");

    // --- Network manager ---
    network = new QNetworkAccessManager(this);

    // --- Signals ---
    connect(addAction, &QAction::triggered, this, &MainWindow::onAddTicket);
    connect(editAction, &QAction::triggered, this, &MainWindow::onEditTicket);
    connect(deleteAction, &QAction::triggered, this, &MainWindow::onDeleteTicket);
    connect(refreshAction, &QAction::triggered, this, [this]() { loadTickets(); });
    connect(searchEdit, &QLineEdit::textChanged, this, &MainWindow::onSearchChanged);

    // --- Data loading ---
    loadTickets();

    // Load dictionaries
    QNetworkAccessManager *dictManager = new QNetworkAccessManager(this);
    // Statuses
    QNetworkRequest statusReq(QUrl(Config::instance().fullApiUrl() + "/ticket_statuses"));
    QNetworkReply *statusReply = dictManager->get(statusReq);
    connect(statusReply, &QNetworkReply::finished, this, [statusReply]{
        if (statusReply->error() == QNetworkReply::NoError) {
            QJsonDocument doc = QJsonDocument::fromJson(statusReply->readAll());
            if (doc.isArray()) {
                for (const QJsonValue &v : doc.array()) {
                    QJsonObject o = v.toObject();
                    int id = o.value("id").toInt(-1);
                    QString label = o.value("label").toString();
                    if (id > 0 && !label.isEmpty()) {
                        TicketItem::statusLabels[id] = label;
                    }
                }
            }
        }
        statusReply->deleteLater();
    });
    // Priorities
    QNetworkRequest prioReq(QUrl(Config::instance().fullApiUrl() + "/ticket_priorities"));
    QNetworkReply *prioReply = dictManager->get(prioReq);
    connect(prioReply, &QNetworkReply::finished, this, [prioReply]{
        if (prioReply->error() == QNetworkReply::NoError) {
            QJsonDocument doc = QJsonDocument::fromJson(prioReply->readAll());
            if (doc.isArray()) {
                for (const QJsonValue &v : doc.array()) {
                    QJsonObject o = v.toObject();
                    int id = o.value("id").toInt(-1);
                    QString label = o.value("label").toString();
                    if (id > 0 && !label.isEmpty()) {
                        TicketItem::priorityLabels[id] = label;
                    }
                }
            }
        }
        prioReply->deleteLater();
    });
    // Departments
    QNetworkRequest depReq(QUrl(Config::instance().fullApiUrl() + "/departments"));
    QNetworkReply *depReply = dictManager->get(depReq);
    connect(depReply, &QNetworkReply::finished, this, [depReply]{
        if (depReply->error() == QNetworkReply::NoError) {
            QJsonDocument doc = QJsonDocument::fromJson(depReply->readAll());
            if (doc.isArray()) {
                for (const QJsonValue &v : doc.array()) {
                    QJsonObject o = v.toObject();
                    int id = o.value("id").toInt(-1);
                    QString name = o.value("name").toString();
                    if (id > 0 && !name.isEmpty()) {
                        TicketItem::departmentNames[id] = name;
                    }
                }
            }
        }
        depReply->deleteLater();
    });
}

void MainWindow::loadTickets() {
    if (userId.isEmpty()) {
        qDebug() << "ERROR: userId is empty!";
        QMessageBox::warning(this, "Error", "Failed to determine user. Please re-login.");
        return;
    }

    QUrl url(apiBaseUrl + "/tickets");
    QUrlQuery query;
    query.addQueryItem("assignee_id", userId);
    url.setQuery(query);

    QNetworkRequest req(url);
    req.setRawHeader("Authorization", "Bearer " + jwtToken.toUtf8());
    QNetworkReply *reply = network->get(req);

    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        if (reply->error() == QNetworkReply::NoError) {
            QByteArray data = reply->readAll();
            QJsonDocument doc = QJsonDocument::fromJson(QString::fromUtf8(data).toUtf8());
            if (doc.isArray()) {
                ticketModel->loadTickets(doc.array());
                m_statusBar->showMessage(QString("Loaded %1 tickets").arg(ticketModel->rowCount()));
            } else {
                qDebug() << "Invalid JSON response format";
                m_statusBar->showMessage("Error: Invalid response format");
            }
        } else {
            QMessageBox::warning(this, "Error", "Error loading tickets: " + reply->errorString());
            m_statusBar->showMessage("Error loading tickets");
        }
        reply->deleteLater();
    });
}

void MainWindow::onAddTicket() {
    qDebug() << "=== onAddTicket() START ===";
    qDebug() << "userId:" << userId;
    qDebug() << "jwtToken length:" << jwtToken.length();
    qDebug() << "currentFilter:" << currentFilter;
    
    if (userId.isEmpty()) {
        qDebug() << "ERROR: userId is empty!";
        QMessageBox::warning(this, "Error", "Failed to determine user. Please re-login.");
        return;
    }
    
    qDebug() << "Creating TicketItem...";
    TicketItem t;
    qDebug() << "TicketItem created successfully";
    
    qDebug() << "Creating TicketDialog...";
    try {
        TicketDialog dlg(t, jwtToken, this, TicketDialog::Create);
        qDebug() << "TicketDialog created successfully";
        
        qDebug() << "Connecting ticketSaved signal...";
        connect(&dlg, &TicketDialog::ticketSaved, this, [this]() { loadTickets(); });
        qDebug() << "Signal connected successfully";
        
        qDebug() << "Executing TicketDialog...";
        dlg.exec();
        qDebug() << "TicketDialog finished execution";
    } catch (const std::exception& e) {
        qDebug() << "EXCEPTION in TicketDialog:" << e.what();
        QMessageBox::critical(this, "Error", QString("Error creating dialog: %1").arg(e.what()));
    } catch (...) {
        qDebug() << "UNKNOWN EXCEPTION in TicketDialog";
        QMessageBox::critical(this, "Error", "Unknown error while creating dialog");
    }
    
    qDebug() << "=== onAddTicket() END ===";
}

void MainWindow::onEditTicket() {
    QModelIndex currentIndex = tableView->currentIndex();
    if (!currentIndex.isValid()) {
        QMessageBox::information(this, "Info", "Select a ticket to edit");
        return;
    }

    TicketItem ticket = ticketModel->getTicket(currentIndex.row());
    TicketDialog dlg(ticket, jwtToken, this, TicketDialog::Edit);
    connect(&dlg, &TicketDialog::ticketSaved, this, [this]() { loadTickets(); });
    dlg.exec();
}

void MainWindow::onDeleteTicket() {
    QModelIndex currentIndex = tableView->currentIndex();
    if (!currentIndex.isValid()) {
        QMessageBox::information(this, "Info", "Select a ticket to delete");
        return;
    }

    TicketItem ticket = ticketModel->getTicket(currentIndex.row());
    QMessageBox::StandardButton reply = QMessageBox::question(this, "Confirmation", 
        QString("Are you sure you want to delete ticket '%1'?").arg(ticket.title),
        QMessageBox::Yes | QMessageBox::No);

    if (reply == QMessageBox::Yes) {
        QUrl url(apiBaseUrl + "/tickets/" + ticket.id);
        QNetworkRequest req(url);
        req.setRawHeader("Authorization", "Bearer " + jwtToken.toUtf8());
        QNetworkReply *reply = network->sendCustomRequest(req, "DELETE");

        connect(reply, &QNetworkReply::finished, this, [this, reply]() {
            if (reply->error() == QNetworkReply::NoError) {
                loadTickets();
                m_statusBar->showMessage("Ticket deleted successfully");
            } else {
                QMessageBox::warning(this, "Error", "Error deleting ticket: " + reply->errorString());
                m_statusBar->showMessage("Error deleting ticket");
            }
            reply->deleteLater();
        });
    }
}

void MainWindow::onSearchChanged(const QString &text) {
    currentFilter = text;
    // Здесь можно добавить фильтрацию тикетов
    qDebug() << "Search filter changed to:" << text;
}

MainWindow::~MainWindow() {}

void MainWindow::onRefreshTickets() {}
