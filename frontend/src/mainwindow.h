#pragma once
#include "ticket_model.h"
#include <QMainWindow>
#include <QString>
#include <QUrlQuery>
#include <QMetaType>
#include <QListWidget>
#include <QTreeView>
#include <QNetworkAccessManager>

class QMenuBar;
class QToolBar;
class QTableView;
class QTreeView;
class QStatusBar;

class MainWindow : public QMainWindow {
    Q_OBJECT
public:
    explicit MainWindow(const QString &jwt, QWidget *parent = nullptr);
    void loadTickets(const QJsonArray& array);
    TicketItem getTicket(int row) const;
    ~MainWindow();
private:
    QString jwtToken;
    QString userId;
    QString apiBaseUrl;
    QMenuBar *menuBar;
    QToolBar *toolBar;
    QTreeView *leftPanel;
    QTableView *ticketTable;
    QTableView* tableView;
    TicketModel *ticketModel = nullptr;
    QNetworkAccessManager *network = nullptr;
    QString searchQuery;
    QString currentFilter;
    QAction *addAction;
    QAction *editAction;
    QAction *deleteAction;
    QAction *refreshAction;
    QStatusBar *m_statusBar;
private slots:
    void onAddTicket();
    void onEditTicket();
    void onDeleteTicket();
    void onRefreshTickets();
    void onSearchChanged(const QString &text);
    void loadTickets();
}; 