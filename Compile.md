# How to Compile (FOR MY TEAM MEMBERS)

This is for building the C++ client.

## 1. Install CMake

Make sure to have cmake installed.

Check if you have it:

```bash
cmake --version
```

## 2. Install SFML

The project uses SFML 3.0.2.

Put SFML inside the project like this:

```text
TCP/
  Dependencies/
    SFML/
```

CMake is set up to look for SFML there automatically.

(It will still search paths such as `usr/local/` for SFML)

## 3. Configure the Project

Make sure to be in `TCP/` then run:

```bash
cmake -B build
```

This creates the `build` folder and makes CMake find the libraries needed by the project, like SFML.

## 4. Build the Executable

After CMake configures successfully, run:

```bash
cmake --build build
```

This actually compiles the code and creates the executable.

On Linux, the executable should be:

```text
build/main
```

On Windows, it may be:

```text
build\Debug\main.exe
```

or:

```text
build\Release\main.exe
```

## VS Code Note

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
