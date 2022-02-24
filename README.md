# mt-set-time - MikroTik set time

The `mt-set-time` tool provides a way to set the time on MikroTik routers based on the current local time.

## Installation
Install the go package
```
# go install github.com/x-way/mt-set-time@latest
```

## Usage
Run the go binary from your local path
```
# mt-set-time -h host1
Device time: jan/01/2000 00:11:25 Europe/Zurich
Device time: feb/23/2022 19:49:37 Europe/Zurich
```

## Parameters
```
Usage of mt-set-time:
  -h string
    	host to set time for
  -i string
    	override IP address used to connect to host
  -m string
    	host-to-config-file mapping (default "mapping.json")
```

## Configuration
The `mapping.json` file contains the configuration details required for connecting to the MikroTik hosts (the `-h` parameter of `mt-set-time` selects which host to connect to)

Example `mapping.json` file:
```
{
  "hosts": [
    {
      "name": "host1",
      "ip": "198.51.100.123",
      "port": 8729,
      "username": "admin",
      "password": "secret"
    },
    {
      "name": "my-other-mikrotik-host",
      "ip": "198.51.100.100",
      "port": 8729,
      "username": "admin",
      "password": "someothersecret",
      "tlsInsecureSkipVerify": true
    }
  ]
}
```
