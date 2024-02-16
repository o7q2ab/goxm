# goxm

## Install:

```sh
go install github.com/o7q2ab/goxm@latest
```

## Commands:

### `binary`

Examine Go binary file.

Aliases: `binary`, `bin`, `b`.

Example:

```sh
goxm b ~/go/bin/gopls
```

Flags:
```
  -b, --build    show the build settings used to build the binary
  -d, --deps     show all the dependency modules
  -h, --help     help for binary
      --latest   show latest versions for all the dependency modules
```

### `path`

Examine all Go binaries found in directories added to PATH environment variable.

Example:

```sh
goxm path
```

Flags:
```
  -b, --build   show the build settings used to build the binary
  -d, --deps    show all the dependency modules
  -h, --help    help for path
```

### `process`

Examine currently running Go processes.

Aliases: `process`, `proc`, `ps`, `p`.

Example:

```sh
goxm ps
```

Flags:
```
  -b, --build           show the build settings used to build the binary
      --conn            show all the connections (TCP, UDP, Unix) used by the process
  -d, --deps            show all the dependency modules
      --filter string   filter by the package name
  -h, --help            help for process
```

### `module`

Examine Go module.

Aliases: `module`, `mod`, `m`.

Example:

```sh
goxm m ./path/to/go.mod
```

Flags:
```
  -h, --help     help for module
```

