#include "mainwindow.h"
#include "models/ticket_model.h"
#include "views/ticket_table_view.h"
#include "models/dictionary_model.h"
#include "views/ticket_dialog.h"
#include "config.h"
#include "ticket_table_view.h"

#include <QSplitter>
#include <QTreeView>
#include <QStandardItemModel>
#include <QHeaderView>
#include <QSettings>
#include <QCloseEvent>
#include <QToolBar>
#include <QTableView>
#include <QHBoxLayout>
#include <QLineEdit>
#include <QPushButton>
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
#include <QVBoxLayout>
#include <QWidget>

QNetworkAccessManager *networkManager = nullptr;

const int FilterTypeRole = Qt::UserRole + 1;
const int FilterValueRole = Qt::UserRole + 2;


MainWindow::MainWindow(const QString &jwt, QWidget *parent)
    : QMainWindow(parent), jwtToken(jwt)
{
    QStringList parts = jwtToken.split('.');
    if (parts.size() >= 2) {
        QByteArray payload = QByteArray::fromBase64(parts[1].toUtf8(), QByteArray::Base64UrlEncoding);
        QJsonDocument doc = QJsonDocument::fromJson(payload);
        if (doc.isObject()) {
            userId = doc.object().value("user_id").toString();
        }
    }
    qDebug() << "Extracted userId =" << userId;
    
    apiBaseUrl = Config::instance().fullApiUrl();
    networkManager = new QNetworkAccessManager(this);

    setupUi();
    loadDictionaries();
}

void MainWindow::setupUi() {
    setWindowTitle("Ticket System");
    
    QSettings settings("MyCompany", "TicketSystem");
    if (settings.contains("geometry")) {
        restoreGeometry(settings.value("geometry").toByteArray());
    } else {
        resize(1024, 768);
    }

    // --- Toolbar ---
    m_toolBar = new QToolBar(this);
    addToolBar(m_toolBar);
    m_addAction = m_toolBar->addAction("➕ Add");
    m_editAction = m_toolBar->addAction("✏️ Edit");
    m_deleteAction = m_toolBar->addAction("🗑️ Delete");
    m_refreshAction = m_toolBar->addAction("🔄 Refresh");
    m_toolBar->addSeparator();

    m_searchEdit = new QLineEdit(this);
    m_searchEdit->setPlaceholderText("Search tickets...");
    m_searchEdit->setMinimumWidth(200);
    m_toolBar->addWidget(m_searchEdit);

    m_searchButton = new QPushButton("Search", this);
    m_toolBar->addWidget(m_searchButton);

    // --- Main Layout (Splitter) ---
    m_splitter = new QSplitter(Qt::Horizontal, this);
    setCentralWidget(m_splitter);

    // --- Left Panel (Filters) ---
    m_filterView = new QTreeView(m_splitter);
    m_filterView->setHeaderHidden(true);
    m_filterModel = new QStandardItemModel(this);
    m_filterView->setModel(m_filterModel);
    
    QStandardItem* rootNode = m_filterModel->invisibleRootItem();
    
    QStandardItem* allTicketsItem = new QStandardItem("All Tickets");
    allTicketsItem->setData("all_tickets", FilterTypeRole);
    rootNode->appendRow(allTicketsItem);
    
    QStandardItem* myTicketsItem = new QStandardItem("My Tickets");
    myTicketsItem->setData("my_tickets", FilterTypeRole);
    rootNode->appendRow(myTicketsItem);
    
    QStandardItem* departmentsItem = new QStandardItem("Departments");
    departmentsItem->setData("department_header", FilterTypeRole);
    departmentsItem->setSelectable(false);
    rootNode->appendRow(departmentsItem);

    // --- Right Panel (Table) ---
    m_tableView = new TicketTableView(this);
    m_ticketModel = new TicketModel(this);
    m_tableView->setModel(m_ticketModel);
    m_tableView->setSelectionBehavior(QAbstractItemView::SelectRows);
    m_tableView->setSelectionMode(QAbstractItemView::SingleSelection);
    m_tableView->horizontalHeader()->setStretchLastSection(true);
    m_tableView->setAlternatingRowColors(true);
    m_tableView->verticalHeader()->setSectionResizeMode(QHeaderView::ResizeToContents);
    // ВФЫвфывфывфывфыв
    m_tableView->setStyleSheet("QTableView { gridline-color: #e0e0e0; font-size: 13px } QHeaderView::section { font-weight: bold; font-size: 14px; text-align: center; }");
    m_tableView->horizontalHeader()->setDefaultAlignment(Qt::AlignCenter);
    m_tableView->setWordWrap(false);
    m_tableView->setSelectionMode(QAbstractItemView::SingleSelection);
    m_tableView->setSelectionBehavior(QAbstractItemView::SelectRows);
    m_tableView->setContentsMargins(8, 4, 8, 4);
    m_tableView->verticalHeader()->setDefaultSectionSize(36);
    // ВФЫвфывфывфывфыв
    
    m_tableView->setItemDelegateForColumn(2, new TicketBadgeDelegate(this));
    m_tableView->setItemDelegateForColumn(3, new TicketPriorityDelegate(this));
    m_tableView->setColumnHidden(6, true);
    m_tableView->setColumnHidden(0, true);

    QWidget *tableViewContainer = new QWidget(this);
    QVBoxLayout *tableViewLayout = new QVBoxLayout(tableViewContainer);
    tableViewLayout->addWidget(m_tableView);
    tableViewLayout->setContentsMargins(0, 0, 0, 0);
    m_splitter->addWidget(tableViewContainer);
    m_splitter->setStretchFactor(0, 0);
    m_splitter->setStretchFactor(1, 1);
    m_splitter->setSizes({200, 800});

    // --- Status Bar ---
    m_statusBar = new QStatusBar(this);
    setStatusBar(m_statusBar);
    m_statusBar->showMessage("Ready");
    
    // --- Connections ---
    connect(m_addAction, &QAction::triggered, this, &MainWindow::onAddTicket);
    connect(m_editAction, &QAction::triggered, this, &MainWindow::onEditTicket);
    connect(m_deleteAction, &QAction::triggered, this, &MainWindow::onDeleteTicket);
    connect(m_refreshAction, &QAction::triggered, this, &MainWindow::loadTickets);
    connect(m_searchButton, &QPushButton::clicked, this, &MainWindow::onSearchTriggered);
    connect(m_searchEdit, &QLineEdit::returnPressed, this, &MainWindow::onSearchTriggered);
    connect(m_filterView, &QTreeView::clicked, this, &MainWindow::onFilterChanged);

    m_filterView->setCurrentIndex(m_filterModel->index(0, 0));
}

void MainWindow::loadDictionaries() {
    m_statusBar->showMessage("Loading initial data...");

    // Statuses
    QNetworkRequest statusReq(QUrl(apiBaseUrl + "/ticket_statuses"));
    QNetworkReply *statusReply = networkManager->get(statusReq);
    connect(statusReply, &QNetworkReply::finished, this, [this, statusReply](){
        if (statusReply->error() == QNetworkReply::NoError) {
            QJsonDocument doc = QJsonDocument::fromJson(statusReply->readAll());
            if (doc.isArray()) {
                for (const QJsonValue &v : doc.array()) {
                    QJsonObject o = v.toObject();
                    TicketItem::statusLabels[o.value("id").toInt()] = o.value("label").toString();
                }
            }
        } else { handleNetworkError(statusReply, "ticket statuses"); }
        onInitialDataLoaded();
        statusReply->deleteLater();
    });

    // Priorities
    QNetworkRequest prioReq(QUrl(apiBaseUrl + "/ticket_priorities"));
    QNetworkReply *prioReply = networkManager->get(prioReq);
    connect(prioReply, &QNetworkReply::finished, this, [this, prioReply](){
        if (prioReply->error() == QNetworkReply::NoError) {
            QJsonDocument doc = QJsonDocument::fromJson(prioReply->readAll());
             if (doc.isArray()) {
                for (const QJsonValue &v : doc.array()) {
                    QJsonObject o = v.toObject();
                    TicketItem::priorityLabels[o.value("id").toInt()] = o.value("label").toString();
                }
            }
        } else { handleNetworkError(prioReply, "ticket priorities"); }
        onInitialDataLoaded();
        prioReply->deleteLater();
    });
    
    // Departments
    QNetworkRequest depReq(QUrl(apiBaseUrl + "/departments"));
    QNetworkReply *depReply = networkManager->get(depReq);
    connect(depReply, &QNetworkReply::finished, this, [this, depReply](){
        if (depReply->error() == QNetworkReply::NoError) {
            QJsonDocument doc = QJsonDocument::fromJson(depReply->readAll());
            if (doc.isArray()) {
                populateDepartments(doc.array());
            }
        } else { handleNetworkError(depReply, "departments"); }
        onInitialDataLoaded();
        depReply->deleteLater();
    });
}

void MainWindow::populateDepartments(const QJsonArray &departments) {
    QStandardItem *departmentsItem = m_filterModel->item(2);
    if (!departmentsItem) return;

    for (const QJsonValue &v : departments) {
        QJsonObject o = v.toObject();
        int id = o.value("id").toInt(-1);
        QString name = o.value("name").toString();
        if (id > 0 && !name.isEmpty()) {
            TicketItem::departmentNames[id] = name;
            QStandardItem* item = new QStandardItem(name);
            item->setData("department", FilterTypeRole);
            item->setData(id, FilterValueRole);
            departmentsItem->appendRow(item);
        }
    }
}

void MainWindow::onInitialDataLoaded() {
    m_dictionariesToLoad--;
    if (m_dictionariesToLoad == 0) {
        m_statusBar->showMessage("Ready");
        loadTickets();
    }
}

void MainWindow::loadTickets() {
    if (userId.isEmpty()) {
        QMessageBox::warning(this, "Error", "User ID is not available. Please re-login.");
        return;
    }
    m_statusBar->showMessage("Loading tickets...");

    QUrl url(apiBaseUrl + "/tickets");
    QUrlQuery query;

    for (auto it = m_currentQueryItems.constBegin(); it != m_currentQueryItems.constEnd(); ++it) {
        query.addQueryItem(it.key(), it.value());
    }

    url.setQuery(query);

    QNetworkRequest req(url);
    req.setRawHeader("Authorization", "Bearer " + jwtToken.toUtf8());
    QNetworkReply *reply = networkManager->get(req);

    connect(reply, &QNetworkReply::finished, this, [this, reply]() {
        if (reply->error() == QNetworkReply::NoError) {
            QJsonDocument doc = QJsonDocument::fromJson(reply->readAll());
            if (doc.isArray()) {
                m_ticketModel->loadTickets(doc.array());
                m_statusBar->showMessage(QString("Loaded %1 tickets").arg(m_ticketModel->rowCount()));
                m_tableView->resizeColumnsToContents();
                m_tableView->viewport()->update();
                m_tableView->update();
            } else {
                qDebug() << "Invalid JSON response for tickets:" << doc;
                m_statusBar->showMessage("Error: Invalid response from server.");
            }
        } else {
            handleNetworkError(reply, "loading tickets");
        }
        reply->deleteLater();
    });
}

void MainWindow::onFilterChanged(const QModelIndex &index) {
    if (!index.isValid()) return;

    const QStandardItem *item = m_filterModel->itemFromIndex(index);
    QString filterType = item->data(FilterTypeRole).toString();

    m_currentQueryItems.remove("assignee_id");
    m_currentQueryItems.remove("department_id");

    if (filterType == "my_tickets") {
        m_currentQueryItems.insert("assignee_id", userId);
        setWindowTitle(QString("Ticket System - My Tickets"));
    } else if (filterType == "department") {
        QString deptId = item->data(FilterValueRole).toString();
        m_currentQueryItems.insert("department_id", deptId);
        setWindowTitle(QString("Ticket System - %1").arg(item->text()));
    } else if (filterType == "all_tickets") {
         setWindowTitle("Ticket System - All Tickets");
    } else {
        return;
    }

    loadTickets();
}

void MainWindow::onSearchTriggered() {
    QString searchText = m_searchEdit->text().trimmed();
    if (!searchText.isEmpty()) {
        m_currentQueryItems.insert("q", searchText);
    } else {
        m_currentQueryItems.remove("q");
    }
    loadTickets();
}


void MainWindow::onAddTicket() {
    TicketItem t;
    TicketDialog dlg(t, jwtToken, this, TicketDialog::Create);
    connect(&dlg, &TicketDialog::ticketSaved, this, &MainWindow::loadTickets);
    dlg.exec();
}

void MainWindow::onEditTicket() {
    QModelIndex currentIndex = m_tableView->currentIndex();
    if (!currentIndex.isValid()) {
        QMessageBox::information(this, "Info", "Select a ticket to edit");
        return;
    }

    TicketItem ticket = m_ticketModel->getTicket(currentIndex.row());
    TicketDialog dlg(ticket, jwtToken, this, TicketDialog::Edit);
    connect(&dlg, &TicketDialog::ticketSaved, this, &MainWindow::loadTickets);
    dlg.exec();
}

void MainWindow::onDeleteTicket() {
    QModelIndex currentIndex = m_tableView->currentIndex();
    if (!currentIndex.isValid()) {
        QMessageBox::information(this, "Info", "Select a ticket to delete");
        return;
    }

    TicketItem ticket = m_ticketModel->getTicket(currentIndex.row());
    QMessageBox::StandardButton reply = QMessageBox::question(this, "Confirmation", 
        QString("Are you sure you want to delete ticket '%1'?").arg(ticket.title),
        QMessageBox::Yes | QMessageBox::No);

    if (reply == QMessageBox::Yes) {
        QUrl url(apiBaseUrl + "/tickets/" + ticket.id);
        QNetworkRequest req(url);
        req.setRawHeader("Authorization", "Bearer " + jwtToken.toUtf8());
        QNetworkReply *reply = networkManager->sendCustomRequest(req, "DELETE");

        connect(reply, &QNetworkReply::finished, this, [this, reply]() {
            if (reply->error() == QNetworkReply::NoError) {
                loadTickets();
                m_statusBar->showMessage("Ticket deleted successfully");
            } else {
                handleNetworkError(reply, "deleting ticket");
            }
            reply->deleteLater();
        });
    }
}

void MainWindow::handleNetworkError(QNetworkReply *reply, const QString& context) {
    QString errorMsg = QString("Error %1: %2").arg(context, reply->errorString());
    qDebug() << errorMsg;
    qDebug() << "Response:" << reply->readAll();
    m_statusBar->showMessage(QString("Error %1").arg(context));
    QMessageBox::warning(this, "Network Error", errorMsg);
}

void MainWindow::closeEvent(QCloseEvent *event) {
    QSettings settings("MyCompany", "TicketSystem");
    settings.setValue("geometry", saveGeometry());
    QMainWindow::closeEvent(event);
}

MainWindow::~MainWindow() {}