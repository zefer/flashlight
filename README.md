# flashlight

An LCD display for an MPD music server on the network.

Will work for any MPD server, but designed to compliment
[mothership][mothership].

## Usage

  ./flashlight -mpdaddr=somehost:6600

## GPIO pins

See `lcd/lcd.go` for GPIO pin mappings.

## Examples

![](https://dl.dropboxusercontent.com/u/89410/project_images/mpdlcd-1.jpg)
![](https://dl.dropboxusercontent.com/u/89410/project_images/mpdlcd-2.jpg)

## License

This project uses the MIT License. See [LICENSE](LICENSE).

[mothership]: https://github.com/zefer/mothership
