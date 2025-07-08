#pragma once
#include "models/ticket_model.h"
#include <QMainWindow>
#include <QString>
#include <QUrlQuery>
#include <QMetaType>
#include <QListWidget>
#include <QTreeView>
#include <QLabel>
#include <QNetworkAccessManager>

class QTableView;
class QToolBar;
class QAction;
class QStatusBar;
class QLineEdit;
class QPushButton;
class QSplitter;
class QTreeView;
class QStandardItemModel;
class QNetworkReply;
class QCloseEvent;
class QModelIndex;
class QJsonArray;

class MainWindow : public QMainWindow {
    Q_OBJECT
public:
    explicit MainWindow(const QString &jwt, QWidget *parent = nullptr);
    ~MainWindow();
protected:
    void closeEvent(QCloseEvent *event) override;
private:
    void setupUi();
    void loadDictionaries();
    void handleNetworkError(QNetworkReply* reply, const QString& context);

    // Auth & API
    QString jwtToken;
    QString userId;
    QString apiBaseUrl;

    // UI Widgets
    QSplitter *m_splitter;
    QTreeView *m_filterView;
    QStandardItemModel *m_filterModel;
    QTableView *m_tableView;
    TicketModel *m_ticketModel;
    QToolBar *m_toolBar;
    QAction *m_addAction;
    QAction *m_editAction;
    QAction *m_deleteAction;
    QAction *m_refreshAction;
    QLineEdit *m_searchEdit;
    QPushButton *m_searchButton;
    QStatusBar *m_statusBar;
    
    // State management
    QMap<QString, QString> m_currentQueryItems;
    int m_dictionariesToLoad = 3;
private slots:
    void loadTickets(); 
    void onAddTicket();
    void onEditTicket();
    void onDeleteTicket();
    void onFilterChanged(const QModelIndex &index);
    void onSearchTriggered();
    void onInitialDataLoaded();
    void populateDepartments(const QJsonArray &departments);
}; 