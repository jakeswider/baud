# baud

A lightweight, high-performance CLI tool for discovering serial ports. Built in **Pure Go** with zero C-dependencies (CGO-free for Windows/Linux), making it incredibly fast and portable.

## Features
- **Fast Discovery:** Instantly list all connected serial devices.
- **Hardware ID Mapping:** Uses an embedded USB database to identify device manufacturers (Arduino, ESP32, etc.).
- **Smart Filtering:** Automatically hides system clutter (Bluetooth, debug ports) to keep your view clean.
- **Zero Dependencies:** No need for `libusb` or `pkg-config`.

---

## Installation

### For macOS (Apple Silicon & Intel)
1. Download the latest binary from the [Releases](https://github.com/jakeswider/baud/releases) page.
2. Open your terminal and run:
   ```bash
   # 1. Make the file executable
   chmod +x ~/Downloads/baud-mac-arm64

   # 2. Bypass macOS security check
   xattr -d com.apple.quarantine ~/Downloads/baud-mac-arm64

   # 3. Move it to your system path
   sudo mv ~/Downloads/baud-mac-arm64 /usr/local/bin/baud


3. Type `baud list` to start!

### For Windows

1. Download `baud-windows.exe` from the [Releases](https://www.google.com/url?sa=E&source=gmail&q=https://github.com/jakeswider/baud/releases) page.
2. Run it directly from your terminal:
```powershell
.\baud-windows.exe list

```



### For Developers (using Go)

If you have Go installed, you can compile and install it directly:

```bash
go install [github.com/jakeswider/baud/cmd@latest](https://github.com/jakeswider/baud/cmd@latest)

```

*Note: Ensure your `$(go env GOPATH)/bin` is in your system `PATH`.*

---

## Usage

### Basic List

Show only physical USB hardware (filters out Bluetooth and virtual ports):

```bash
baud list

```

### Identify Devices

Show the manufacturer and product names (e.g., "Arduino Uno"):

```bash
baud list --name

```

### Show All

Include internal system ports and Bluetooth:

```bash
baud list --showall

```

---

## Development

To build from source:

1. **Clone the repo:**
```bash
git clone [https://github.com/jakeswider/baud.git](https://github.com/jakeswider/baud.git)

```


2. **Build:**
```bash
go build -o baud ./cmd

```



---

## License
Distributed under the MIT License.