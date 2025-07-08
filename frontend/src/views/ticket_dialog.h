#pragma once
#include <QDialog>
#include "models/ticket_model.h"
#include <QPushButton>
#include <QNetworkAccessManager>
#include <QNetworkReply>
#include <QComboBox>
#include "models/comment_model.h"
#include <QListView>
#include "models/attachment_model.h"
#include <QFileDialog>
#include <QNetworkAccessManager>
#include <QMap>
#include <QPixmap>
#include <QLabel>
#include <QVBoxLayout>
#include <QScreen>
#include <QMouseEvent>
#include <QPainter>

class QTabWidget;
class QWidget;
class QLineEdit;
class QTextEdit;
class QTableView;

class ImagePreviewDialog;

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
    void requestAttachmentImage(const QString &attId, const QString &ticketId);
    void setCurrentTab(int index);
    void showImagePreview(const QPixmap &pixmap);
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
    QListView *m_attachmentsListView = nullptr;
    AttachmentModel *m_attachmentModel = nullptr;
    QPushButton *m_deleteAttachmentBtn = nullptr;
    QPushButton *m_uploadAttachmentBtn = nullptr;
    QNetworkAccessManager *m_attachmentImageManager = nullptr;
    QMap<QString, QPixmap> m_attachmentPixmaps;
    void onAttachmentImageLoaded(const QString &attId, const QPixmap &pixmap);
    void loadHistory();
    void loadComments();
    void loadAttachments();
    void loadDepartments();
    void loadStatuses();
    void loadPriorities();
    void loadUsers();
    void filterAssigneesByDepartment(int departmentId);
    void postNewComment();
    void decodeJwtToken();
    QString m_userRole;
    QString m_userId;
    void deleteSelectedAttachment();
    void uploadAttachment();
    ImagePreviewDialog *m_imagePreviewDialog = nullptr;
};

class ImagePreviewDialog : public QDialog {
    Q_OBJECT
public:
    ImagePreviewDialog(QWidget *parent = nullptr) : QDialog(parent) {
        setWindowFlags(Qt::Dialog | Qt::FramelessWindowHint);
        setAttribute(Qt::WA_TranslucentBackground);
        setModal(true);
        QVBoxLayout *layout = new QVBoxLayout(this);
        layout->setContentsMargins(0,0,0,0);
        layout->setAlignment(Qt::AlignCenter);
        m_label = new QLabel(this);
        m_label->setAlignment(Qt::AlignCenter);
        layout->addWidget(m_label);
        setLayout(layout);
    }
    void setPixmap(const QPixmap &pix) {
        m_label->setPixmap(pix.scaled(qApp->primaryScreen()->size() * 0.7, Qt::KeepAspectRatio, Qt::SmoothTransformation));
    }
protected:
    void paintEvent(QPaintEvent *event) override {
        QPainter p(this);
        p.setRenderHint(QPainter::Antialiasing);
        p.fillRect(rect(), QColor(0,0,0,180));
        QDialog::paintEvent(event);
    }
    void mousePressEvent(QMouseEvent *event) override {
        if (!m_label->geometry().contains(event->pos())) {
            close();
        }
    }
private:
    QLabel *m_label;
}; 