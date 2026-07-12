# baud

A command-line utility for discovering and monitoring serial port connections.

## Features

* Lists physical USB serial devices.
* Identifies device manufacturers and products using an embedded USB database.
* Monitors for hardware connections and disconnections in real-time.
* Ignores internal system ports, virtual ports, and Bluetooth connections by default.
* Runs natively as a single binary without requiring `libusb` or external system libraries.

---

## Installation

Download the latest compiled binaries for macOS, Linux, or Windows from the [Releases](https://github.com/jakeswider/baud/releases) page.

### macOS (Apple Silicon & Intel)

```bash
# 1. Make the downloaded file executable
chmod +x ~/Downloads/baud-mac-arm64

# 2. Bypass the macOS Gatekeeper security check
xattr -d com.apple.quarantine ~/Downloads/baud-mac-arm64

# 3. Move the binary to your system path
sudo mv ~/Downloads/baud-mac-arm64 /usr/local/bin/baud
```

### Windows

Extract the downloaded executable and run it directly from your terminal, or add it to your system's `PATH` variable for global access:

```powershell
.\baud-windows-amd64.exe list
```

### Go Developers

If you have the Go toolchain installed, you can compile and install the utility directly:

```bash
go install [github.com/jakeswider/baud/cmd@latest](https://github.com/jakeswider/baud/cmd@latest)
```
*Note: Ensure `$(go env GOPATH)/bin` is included in your system `PATH`.*

---

## Usage

### list

List physical USB hardware connected to the machine:
```bash
baud list
```

Include the manufacturer and product names in the output table:
```bash
baud list --name
```

Override the default filter to display all internal system ports and virtual connections:
```bash
baud list --showall
```

### watch

Monitor the system for any serial devices being plugged in or unplugged. Set an optional timer in seconds using the `--time` or `-t` flag:
```bash
baud watch --time 30
```

---

## Development

To build the project from the source code:

```bash
# Clone the repository
git clone [https://github.com/jakeswider/baud.git](https://github.com/jakeswider/baud.git)

# Navigate to the directory and compile the binary
cd baud
go build -o baud ./cmd
```

---

## License

Distributed under the MIT License.