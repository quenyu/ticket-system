#include "config.h"
#include <QCoreApplication>
#include <QSettings>
#include <QDebug>

Config& Config::instance() {
    static Config instance;
    return instance;
}

Config::Config() {
    m_apiVersion = "v1";
    
    QString envUrl = qEnvironmentVariable("API_BASE_URL");
    if (!envUrl.isEmpty()) {
        m_apiBaseUrl = envUrl;
        qDebug() << "Using API URL from environment:" << m_apiBaseUrl;
    } else {
        QSettings settings(QCoreApplication::applicationDirPath() + "/config.ini", QSettings::IniFormat);
        m_apiBaseUrl = settings.value("api/base_url", "http://localhost:8080").toString();
        qDebug() << "Using API URL from config file:" << m_apiBaseUrl;
    }
    
    if (m_apiBaseUrl.endsWith('/')) {
        m_apiBaseUrl.chop(1);
    }
}

QString Config::apiBaseUrl() const {
    return m_apiBaseUrl;
}

QString Config::apiVersion() const {
    return m_apiVersion;
}

QString Config::fullApiUrl() const {
    return m_apiBaseUrl + "/api/" + m_apiVersion;
}

void Config::setApiBaseUrl(const QString& url) {
    m_apiBaseUrl = url;
    if (m_apiBaseUrl.endsWith('/')) {
        m_apiBaseUrl.chop(1);
    }
    
    QSettings settings(QCoreApplication::applicationDirPath() + "/config.ini", QSettings::IniFormat);
    settings.setValue("api/base_url", m_apiBaseUrl);
    settings.sync();
} 