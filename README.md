# baud

A lightweight, high-performance CLI tool for discovering serial ports. Built in Pure Go with zero C-dependencies.

## Features
- **Fast Discovery:** Instantly list all connected serial devices.
- **Hardware ID Mapping:** Uses an embedded USB database to identify device manufacturers (Arduino, ESP32, etc.).
- **Smart Filtering:** Automatically hides system clutter (Bluetooth, debug ports) while allowing you to toggle them back with `-a`.
- **Zero Dependencies:** No need for `libusb` or `pkg-config`.

## Installation

### Using Go
```bash
go install [github.com/jakeswider/baud@latest](https://github.com/jakeswider/baud@latest)