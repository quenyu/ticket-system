#pragma once
#include <QDialog>
#include "models/ticket_model.h"
#include <QPushButton>
#include <QNetworkAccessManager>
#include <QNetworkReply>
#include <QComboBox>
#include "models/comment_model.h"
#include <QListView>

class QTabWidget;
class QWidget;
class QLineEdit;
class QTextEdit;
class QTableView;

struct UserInfo {
    QString username;
    QString userId;
    int departmentId;
};

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
    // History, Comments, Attachments: QTableView or QListView
    QTableView *historyView;
    QTableView *commentsView;
    QTableView *attachmentsView;
    QPushButton *saveButton;
    QNetworkAccessManager *network;
    QComboBox *departmentCombo;
    QComboBox *statusCombo;
    QComboBox *priorityCombo;
    QComboBox *assigneeCombo;
    QVector<UserInfo> users;
    QPushButton *overviewSaveBtn = nullptr;
    QPushButton *overviewCancelBtn = nullptr;
    QListView *m_commentsListView = nullptr;
    CommentModel *m_commentModel = nullptr;
    QTextEdit *m_newCommentEdit = nullptr;
    QPushButton *m_postCommentBtn = nullptr;
    void loadHistory();
    void loadComments();
    void loadAttachments();
    void loadDepartments();
    void loadStatuses();
    void loadPriorities();
    void loadUsers();
    void filterAssigneesByDepartment(int departmentId);
    void postNewComment();
}; 