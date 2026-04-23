cmake_minimum_required(VERSION 3.22)
project(main LANGUAGES CXX)

# Sugeneruoja compile_commands.json, padeda IDE kad ne errorintu (vscode su clang extensionu)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON) 

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(SFML 3 REQUIRED COMPONENTS Graphics Window System)

add_executable(main main.cpp)
target_link_libraries(main PRIVATE SFML::Graphics SFML::Window SFML::System)
