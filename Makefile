# Linux
main: main.cpp
	g++ main.cpp -o main $(shell pkg-config --cflags --libs sfml-all)
