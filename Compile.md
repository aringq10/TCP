# How to Compile (FOR MY TEAM MEMBERS)

This is for building the C++ client.

## 1. Install CMake

Make sure to have cmake installed.

Check if you have it:

```bash
cmake --version
```

## 2. Install Dependencies

Install Boost and SFML 3.0.2:

Linux users can install Boost with one command:

```bash
sudo apt install libboost-all-dev
```

Windows users should download the library from https://www.boost.org/releases/latest/ 

## 3. Put in correct directories

Make sure to install dependencies in the correct directories where cmake can find them.

For this project you can save the dependencies inside Dependencies folder

Example:
```text
Dependencies/
  SFML-3.0.2
  boost_1_91_0
```

The project uses SFML 3.0.2.

CMake is set up to pick the right folder automatically.

(It will still search paths such as `usr/` for SFML)

## 4. Configure the Project

Make sure to be in `TCP/` then run:

```bash
cmake -B build
```

This creates the `build` folder and makes CMake find the libraries needed by the project, like SFML.

## 5. Build the Executable

After CMake configures successfully, run:

```bash
cmake --build build
```

This actually compiles the code and creates the executable.

On Linux, the executable should be:

```text
build/main
```

On Windows it should be:

```text
build\Debug\main.exe
```

## VS Code Note

!!! ONLY WORKS ON LINUX (OR NINJA CMAKE)!!!

Cmake is configured to create this file:

```text
build/compile_commands.json
```

VS Code can use it so includes are found correctly and there are fewer false red squiggly lines.

For `clangd`, add this to `.vscode/settings.json`:

```json
{
  "clangd.arguments": [
    "--compile-commands-dir=build"
  ]
}
```
