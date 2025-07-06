#pragma once
#include <QString>
#include <QObject>

class Config : public QObject {
    Q_OBJECT
public:
    static Config& instance();
    
    QString apiBaseUrl() const;
    QString apiVersion() const;
    QString fullApiUrl() const;
    
    void setApiBaseUrl(const QString& url);
    
private:
    Config();
    ~Config() = default;
    Config(const Config&) = delete;
    Config& operator=(const Config&) = delete;
    
    QString m_apiBaseUrl;
    QString m_apiVersion;
}; 