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
    // --- Заголовок и размер ---
    setWindowTitle("Ticket System - All Tickets");
    resize(1200, 800);

    // --- JWT → userId из payload ---
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

    // --- Меню и тулбар ---
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

    // Поиск
    QLineEdit *searchEdit = new QLineEdit(this);
    searchEdit->setPlaceholderText("Search tickets...");
    searchEdit->setMinimumWidth(200);
    toolBar->addWidget(searchEdit);

    // --- Таблица ---
    ticketModel = new TicketModel(this);
    tableView = new QTableView(this);
    tableView->setModel(ticketModel);
    tableView->setSelectionBehavior(QAbstractItemView::SelectRows);
    tableView->setSelectionMode(QAbstractItemView::SingleSelection);
    tableView->horizontalHeader()->setStretchLastSection(true);
    tableView->setAlternatingRowColors(true);
    setCentralWidget(tableView);

    // --- Статус бар ---
    m_statusBar = new QStatusBar(this);
    setStatusBar(m_statusBar);
    m_statusBar->showMessage("Ready");

    // --- Сетевой менеджер ---
    network = new QNetworkAccessManager(this);

    // --- Сигналы ---
    connect(addAction, &QAction::triggered, this, &MainWindow::onAddTicket);
    connect(editAction, &QAction::triggered, this, &MainWindow::onEditTicket);
    connect(deleteAction, &QAction::triggered, this, &MainWindow::onDeleteTicket);
    connect(refreshAction, &QAction::triggered, this, [this]() { loadTickets(); });
    connect(searchEdit, &QLineEdit::textChanged, this, &MainWindow::onSearchChanged);

    // --- Загрузка данных ---
    loadTickets();
}

void MainWindow::loadTickets() {
    if (userId.isEmpty()) {
        qDebug() << "ERROR: userId is empty!";
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
            QJsonDocument doc = QJsonDocument::fromJson(data);
            if (doc.isArray()) {
                ticketModel->loadTickets(doc.array());
                m_statusBar->showMessage(QString("Loaded %1 tickets").arg(ticketModel->rowCount()));
            } else {
                qDebug() << "Invalid JSON response format";
                m_statusBar->showMessage("Error: Invalid response format");
            }
        } else {
            QMessageBox::warning(this, "Ошибка", "Ошибка загрузки тикетов: " + reply->errorString());
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
        QMessageBox::warning(this, "Ошибка", "Не удалось определить пользователя. Попробуйте перезайти в систему.");
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
        QMessageBox::critical(this, "Ошибка", QString("Ошибка создания диалога: %1").arg(e.what()));
    } catch (...) {
        qDebug() << "UNKNOWN EXCEPTION in TicketDialog";
        QMessageBox::critical(this, "Ошибка", "Неизвестная ошибка при создании диалога");
    }
    
    qDebug() << "=== onAddTicket() END ===";
}

void MainWindow::onEditTicket() {
    QModelIndex currentIndex = tableView->currentIndex();
    if (!currentIndex.isValid()) {
        QMessageBox::information(this, "Информация", "Выберите тикет для редактирования");
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
        QMessageBox::information(this, "Информация", "Выберите тикет для удаления");
        return;
    }

    TicketItem ticket = ticketModel->getTicket(currentIndex.row());
    QMessageBox::StandardButton reply = QMessageBox::question(this, "Подтверждение", 
        QString("Вы уверены, что хотите удалить тикет '%1'?").arg(ticket.title),
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
                QMessageBox::warning(this, "Ошибка", "Ошибка удаления тикета: " + reply->errorString());
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
