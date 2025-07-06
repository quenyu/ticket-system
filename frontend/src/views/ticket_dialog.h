#pragma once
#include <QDialog>
#include "../ticket_model.h"
#include <QPushButton>
#include <QNetworkAccessManager>
#include <QNetworkReply>
#include <QComboBox>

class QTabWidget;
class QWidget;
class QLineEdit;
class QTextEdit;
class QTableView;

class TicketDialog : public QDialog {
    Q_OBJECT
public:
    enum Mode { Create, Edit };
    TicketDialog(const TicketItem &ticket, const QString &jwt, QWidget *parent = nullptr, Mode mode = Edit);
    void setCurrentTab(int index);
signals:
    void ticketSaved();
private slots:
    void onSaveClicked();
    void onReplyFinished(QNetworkReply* reply);
private:
    TicketItem m_ticket;
    QString m_jwtToken;
    Mode m_mode;
    QTabWidget *tabs;
    QWidget *overviewTab;
    QWidget *historyTab;
    QWidget *commentsTab;
    QWidget *attachmentsTab;
    // Overview fields
    QLineEdit *titleEdit;
    QTextEdit *descEdit;
    // History, Comments, Attachments: QTableView или QListView
    QTableView *historyView;
    QTableView *commentsView;
    QTableView *attachmentsView;
    QPushButton *saveButton;
    QNetworkAccessManager *network;
    QComboBox *departmentCombo;
    QComboBox *statusCombo;
    QComboBox *priorityCombo;
    void loadHistory();
    void loadComments();
    void loadAttachments();
    void loadDepartments();
    void loadStatuses();
    void loadPriorities();
}; 