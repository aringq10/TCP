cmake_minimum_required(VERSION 3.22)
project(main LANGUAGES CXX)

# Sugeneruoja compile_commands.json, padeda IDE kad ne errorintu (vscode su clang extensionu) Veikia tik specialuose cmake versijose (pvz. NINJA)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON) 

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

set(TCP_DEPENDENCIES_DIR "${CMAKE_CURRENT_SOURCE_DIR}/Dependencies")
file(GLOB TCP_LOCAL_PACKAGE_DIRS CONFIGURE_DEPENDS
    "${TCP_DEPENDENCIES_DIR}"
    "${TCP_DEPENDENCIES_DIR}/*"
    "${TCP_DEPENDENCIES_DIR}/*/*"
)

list(APPEND CMAKE_PREFIX_PATH ${TCP_LOCAL_PACKAGE_DIRS})

find_package(SFML 3 REQUIRED
    COMPONENTS Graphics Window System
    CONFIG
    HINTS ${TCP_LOCAL_PACKAGE_DIRS}
)
find_package(Threads REQUIRED)

find_path(Boost_INCLUDE_DIR
    NAMES boost/version.hpp
    HINTS ${TCP_LOCAL_PACKAGE_DIRS}
    PATH_SUFFIXES include
)

if(NOT Boost_INCLUDE_DIR)
    message(FATAL_ERROR "Boost not found. Put Boost inside ${TCP_DEPENDENCIES_DIR}.")
endif()

set(Boost_INCLUDE_DIRS "${Boost_INCLUDE_DIR}")

add_executable(main
    main.cpp
    client_ui/UI.cpp
    chess_logic/board.cpp
    chess_logic/knight.cpp
    chess_logic/pawn.cpp
    chess_logic/rook.cpp
    chess_logic/bishop.cpp
    chess_logic/queen.cpp
    chess_logic/king.cpp
    client_net/ChessNetwork.cpp
)

if(WIN32)
    target_compile_definitions(main PRIVATE _WIN32_WINNT=0x0A00)
endif()

target_include_directories(main PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/client_ui
    ${CMAKE_CURRENT_SOURCE_DIR}/chess_logic
    ${CMAKE_CURRENT_SOURCE_DIR}/client_net
    ${Boost_INCLUDE_DIRS}
)

target_link_libraries(main PRIVATE
    SFML::Graphics
    SFML::Window
    SFML::System
    Threads::Threads
)

if(WIN32)
    target_link_libraries(main PRIVATE ws2_32 mswsock)
endif()

if(WIN32)
    add_custom_command(TARGET main POST_BUILD
        COMMAND ${CMAKE_COMMAND} -E copy_if_different
            $<TARGET_RUNTIME_DLLS:main>
            $<TARGET_FILE_DIR:main>
        COMMAND_EXPAND_LISTS
    )
endif()
