cmake_minimum_required(VERSION 3.22)
project(main LANGUAGES CXX)

# Sugeneruoja compile_commands.json, padeda IDE kad ne errorintu (vscode su clang extensionu)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON) 

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

if(WIN32)
    list(APPEND CMAKE_PREFIX_PATH "C:/local")
    file(GLOB TCP_LOCAL_PACKAGE_DIRS
        "C:/local/*"
        "C:/local/*/*"
    )
else()
    list(APPEND CMAKE_PREFIX_PATH "/usr/local")
    file(GLOB TCP_LOCAL_PACKAGE_DIRS
        "/usr/local/*"
        "/usr/local/*/*"
    )
endif()

list(APPEND CMAKE_PREFIX_PATH ${TCP_LOCAL_PACKAGE_DIRS})

find_package(SFML 3 REQUIRED COMPONENTS Graphics Window System)
find_package(Threads REQUIRED)

set(Boost_NO_BOOST_CMAKE ON)
find_package(Boost REQUIRED)

add_executable(main
    main.cpp
    client_ui/UI.cpp
    chess_logic/board.cpp
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
