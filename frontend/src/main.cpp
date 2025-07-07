#include <memory>
#include <QApplication>
#include <QMetaType>
#include <QDebug>
#include "views/login_dialog.h"
#include "mainwindow.h"

int main(int argc, char *argv[]) {
    qRegisterMetaType<TicketItem>("TicketItem");
    QApplication app(argc, argv);
    LoginDialog login;
    QString token;
    std::unique_ptr<MainWindow> w;
    if (login.exec() == QDialog::Accepted && !(token = login.jwtToken()).isEmpty()) {
        qDebug() << "Login accepted, token:" << token;
        w = std::make_unique<MainWindow>(token);
        qDebug() << "MainWindow created, calling show()";
        w->show();
        qDebug() << "MainWindow shown, starting event loop";
        return app.exec();
    } else {
        qDebug() << "Login rejected or token empty";
    }
    return 0;
} 