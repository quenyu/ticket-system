#pragma once
#include <QAbstractTableModel>
#include <QVector>
#include <QString>
#include <QVariant>

struct DictionaryItem {
    QString id;
    QString col1;
    QString col2;
    QString col3;
};

Q_DECLARE_METATYPE(DictionaryItem)

class DictionaryModel : public QAbstractTableModel {
    Q_OBJECT
public:
    DictionaryModel(QObject *parent = nullptr);
    int rowCount(const QModelIndex &parent = QModelIndex()) const override;
    int columnCount(const QModelIndex &parent = QModelIndex()) const override;
    QVariant data(const QModelIndex &index, int role = Qt::DisplayRole) const override;
    QVariant headerData(int section, Qt::Orientation orientation, int role) const override;
    void setItems(const QVector<DictionaryItem> &items, const QStringList &headers);
private:
    QVector<DictionaryItem> m_items;
    QStringList m_headers;
}; 