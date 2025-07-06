#include "config.h"
#include <QCoreApplication>
#include <QSettings>
#include <QDebug>

Config& Config::instance() {
    static Config instance;
    return instance;
}

Config::Config() {
    // Загружаем настройки из переменных окружения или файла конфигурации
    m_apiVersion = "v1";
    
    // Сначала пробуем переменную окружения
    QString envUrl = qEnvironmentVariable("API_BASE_URL");
    if (!envUrl.isEmpty()) {
        m_apiBaseUrl = envUrl;
        qDebug() << "Using API URL from environment:" << m_apiBaseUrl;
    } else {
        // Затем пробуем файл конфигурации
        QSettings settings(QCoreApplication::applicationDirPath() + "/config.ini", QSettings::IniFormat);
        m_apiBaseUrl = settings.value("api/base_url", "http://localhost:8080").toString();
        qDebug() << "Using API URL from config file:" << m_apiBaseUrl;
    }
    
    // Убираем trailing slash если есть
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
    
    // Сохраняем в файл конфигурации
    QSettings settings(QCoreApplication::applicationDirPath() + "/config.ini", QSettings::IniFormat);
    settings.setValue("api/base_url", m_apiBaseUrl);
    settings.sync();
} 