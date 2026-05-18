cmake_minimum_required(VERSION 3.22)
project(main LANGUAGES CXX)

# Sugeneruoja compile_commands.json, padeda IDE kad ne errorintu (vscode su clang extensionu)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON) 

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# If SFML is placed inside this repository as TCP/Dependencies/SFML,
# CMake will search there before looking in system install paths.
list(APPEND CMAKE_PREFIX_PATH "${CMAKE_CURRENT_SOURCE_DIR}/Dependencies/SFML")

find_package(SFML 3 REQUIRED COMPONENTS Graphics Window System)
find_package(Threads REQUIRED)

add_executable(main
    main.cpp
    client_ui/UI.cpp
    chess_logic/board.cpp
    client_net/ChessNetwork.cpp
)

target_include_directories(main PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/client_ui
    ${CMAKE_CURRENT_SOURCE_DIR}/chess_logic
    ${CMAKE_CURRENT_SOURCE_DIR}/client_net
)

target_link_libraries(main PRIVATE
    SFML::Graphics
    SFML::Window
    SFML::System
    Threads::Threads
)
