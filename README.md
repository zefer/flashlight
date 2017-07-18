# flashlight

An LCD display for an MPD music server on the network.

Will work for any MPD server, but designed to compliment
[mothership][mothership].

## Usage

  export GO15VENDOREXPERIMENT=1
  go build
  ./flashlight -mpdaddr=somehost:6600

## GPIO pins

See `lcd/lcd.go` for GPIO pin mappings.

## Examples

![](https://user-images.githubusercontent.com/101193/28341797-0bc76894-6c0d-11e7-83cb-b49263554768.jpg)
![](https://user-images.githubusercontent.com/101193/28341800-1317cc7e-6c0d-11e7-89da-a37be3303e06.jpg)

## Deploy

Deployment is simple, transfer the binary & run it. A complete example is
provided below:

* [Example server configuration](https://github.com/zefer/ansible/tree/master/roles/flashlight)
  (using Ansible)
* [Example deploy script](bin/deploy)

## License

This project uses the MIT License. See [LICENSE](LICENSE).

[mothership]: https://github.com/zefer/mothership
