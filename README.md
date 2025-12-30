# flashlight

An LCD display for an MPD music server on the network.

Will work for any MPD server, but designed to compliment
[mothership][mothership].

## Requirements

- Go 1.23 or later
- Raspberry Pi (2, 3, 4, Zero, or compatible)
- HD44780-compatible 16x2 character LCD display
- GPIO wiring as specified below

## GPIO Pin Mappings

Pin configuration (BCM GPIO numbers):

- RS (Register Select): GPIO 7
- EN (Enable): GPIO 8
- D4-D7 (Data pins): GPIO 25, 24, 23, 17
- Backlight: GPIO 11

See `lcd/lcd.go` for the complete pin configuration.

## Building

### Cross-compile for Raspberry Pi

```bash
make build-arm
```

This creates a statically-linked ARM binary for Raspberry Pi 2/Zero (ARMv7).

### Local build

Note: Local builds on macOS/Windows will fail due to Linux-specific GPIO code. Always cross-compile for the target platform.

## Usage

```bash
./flashlight -mpdaddr=somehost:6600
```

### Command-line flags

- `-mpdaddr` - MPD server address (default: 127.0.0.1:6600)
- `-abprojectid` - Airbrake project ID (optional)
- `-abapikey` - Airbrake API key (optional)
- `-abenv` - Airbrake environment name (default: development)

## Examples

![](https://user-images.githubusercontent.com/101193/28341797-0bc76894-6c0d-11e7-83cb-b49263554768.jpg)
![](https://user-images.githubusercontent.com/101193/28341800-1317cc7e-6c0d-11e7-89da-a37be3303e06.jpg)

## Deployment

### Using Make

The Makefile provides automated deployment to your Raspberry Pi:

```bash
make deploy
```

This will:
1. Cross-compile the ARM binary
2. Upload it to the server via SCP
3. Stop the systemd service
4. Install the new binary
5. Restart the service

### Configuration

Set environment variables to customize deployment:

```bash
SERVER_HOST=mypi SERVER_USER=pi make deploy
```

Default values:
- `SERVER_HOST`: music
- `SERVER_USER`: joe

### Check service status

```bash
make status
```

### Additional resources

* [Example server configuration](https://github.com/zefer/ansible/tree/master/roles/flashlight)
  (using Ansible)

## Dependencies

This project uses:

- [periph.io](https://periph.io/) - Modern Go library for peripheral I/O (GPIO control)
- [gompd](https://github.com/zefer/gompd) - MPD protocol client
- [mothership](https://github.com/zefer/mothership) - MPD client wrapper

## License

This project uses the MIT License. See [LICENSE](LICENSE).

[mothership]: https://github.com/zefer/mothership
