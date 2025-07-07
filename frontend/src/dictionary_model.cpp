#include "dictionary_model.h"

DictionaryModel::DictionaryModel(QObject *parent)
    : QAbstractTableModel(parent) {}

int DictionaryModel::rowCount(const QModelIndex &parent) const {
    if (parent.isValid()) return 0;
    return m_items.size();
}

int DictionaryModel::columnCount(const QModelIndex &parent) const {
    if (parent.isValid()) return 0;
    return m_headers.size();
}

QVariant DictionaryModel::data(const QModelIndex &index, int role) const {
    if (!index.isValid() || index.row() >= m_items.size() || index.column() >= m_headers.size()) return QVariant();
    const DictionaryItem &item = m_items[index.row()];
    if (role == Qt::DisplayRole) {
        switch (index.column()) {
        case 0: return item.id;
        case 1: return item.col1;
        case 2: return item.col2;
        case 3: return item.col3;
        }
    }
    return QVariant();
}

QVariant DictionaryModel::headerData(int section, Qt::Orientation orientation, int role) const {
    if (orientation == Qt::Horizontal && role == Qt::DisplayRole && section < m_headers.size()) {
        return m_headers[section];
    }
    return QVariant();
}

void DictionaryModel::setItems(const QVector<DictionaryItem> &items, const QStringList &headers) {
    beginResetModel();
    m_items = items;
    m_headers = headers;
    endResetModel();
} 