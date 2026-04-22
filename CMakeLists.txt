cmake_minimum_required(VERSION 3.22)
project(main LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

find_package(SFML 3 REQUIRED COMPONENTS Graphics Window System)

add_executable(main main.cpp)
target_link_libraries(main PRIVATE SFML::Graphics SFML::Window SFML::System)
