cmake_minimum_required(VERSION 3.16)
project(ticket_frontend LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)
set(CMAKE_AUTOUIC ON)

find_package(Qt6 REQUIRED COMPONENTS Widgets Network)

set(SOURCES
    src/main.cpp
    src/views/login_dialog.cpp
    src/network/api_client.cpp
    src/mainwindow.cpp
    src/ticket_model.cpp
    src/views/ticket_dialog.cpp
    src/views/register_dialog.cpp
    src/dictionary_model.cpp
    src/config.cpp
)

set(HEADERS
    src/views/login_dialog.h
    src/network/api_client.h
    src/mainwindow.h
    src/ticket_model.h
    src/views/ticket_dialog.h
    src/views/register_dialog.h
    src/dictionary_model.h
    src/config.h
)

add_executable(ticket_frontend
    src/main.cpp
    src/network/api_client.cpp
    src/views/login_dialog.cpp
    src/mainwindow.cpp
    src/ticket_model.cpp
    src/views/ticket_dialog.cpp
    src/views/register_dialog.cpp
    src/dictionary_model.cpp
    src/config.cpp
)

target_link_libraries(ticket_frontend
    Qt6::Widgets
    Qt6::Network
) 