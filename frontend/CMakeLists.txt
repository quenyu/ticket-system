cmake_minimum_required(VERSION 3.16)
project(ticket_frontend LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)
set(CMAKE_AUTOUIC ON)

find_package(Qt6 REQUIRED COMPONENTS Widgets Network Gui)

set(SOURCES
    src/main.cpp
    src/views/login_dialog.cpp
    src/views/ticket_dialog.cpp
    src/views/register_dialog.cpp
    src/views/ticket_table_view.cpp
    src/network/api_client.cpp
    src/models/ticket_model.cpp
    src/models/dictionary_model.cpp
    src/models/comment_model.cpp
    src/models/attachment_model.cpp
    src/mainwindow.cpp
    src/config.cpp
)

set(HEADERS
    src/views/login_dialog.h
    src/views/ticket_dialog.h
    src/views/register_dialog.h
    src/views/ticket_table_view.h
    src/network/api_client.h
    src/models/ticket_model.h
    src/models/dictionary_model.h
    src/mainwindow.h
    src/config.h
)

add_executable(ticket_frontend
    ${SOURCES}
)

target_include_directories(ticket_frontend PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}/src
    ${CMAKE_CURRENT_SOURCE_DIR}/src/models
    ${CMAKE_CURRENT_SOURCE_DIR}/src/views
    ${CMAKE_CURRENT_SOURCE_DIR}/src/network
)

target_link_libraries(ticket_frontend
    Qt6::Widgets
    Qt6::Network
    Qt6::Gui
) 