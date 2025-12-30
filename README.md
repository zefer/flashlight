# flashlight

An LCD display for an MPD music server on the network.

Will work for any MPD server, but designed to compliment
[mothership][mothership].

## Requirements

- Go 1.23 or later
- Raspberry Pi (2, 3, 4, Zero, or compatible)
- HD44780-compatible 16x2 I2C character LCD display
- I2C wiring as specified below

## I2C LCD Wiring

The LCD uses an I2C backpack (PCF8574) for 4-wire connection via **software I2C on bus 4**:

| LCD Pin | Wire Color | Purpose | GPIO | Physical Pin |
|---------|------------|---------|------|--------------|
| VSS | Black | Ground | - | Pin 6 (or any GND) |
| VDD | Red | 5V Power | - | Pin 2 or 4 |
| SDA | White | I2C Data | GPIO 23 | Pin 16 |
| SCL | Brown | I2C Clock | GPIO 24 | Pin 18 |

I2C address: `0x27` (bus 4)

### Required Configuration

Add this line to `/boot/config.txt` to enable software I2C on GPIO 23/24:

```
dtoverlay=i2c-gpio,bus=4,i2c_gpio_sda=23,i2c_gpio_scl=24
```

Reboot after adding this configuration.

See `lcd/lcd.go` for the complete I2C configuration.

## HiFiBerry Compatibility

This project uses **software I2C on bus 4** (GPIO 23/24) instead of the standard hardware I2C bus 1 (GPIO 2/3). This allows the LCD to coexist with the HiFiBerry DAC+ which uses bus 1 for its PCM512x audio codec (address `0x4d`). Both devices can now operate simultaneously without conflicts.

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
