# dizzy

> Note: this is still WIP. Use at your own peril.

Endpoint sniffer, take in a list of words and then make requests and notify if the request was successful against the given `subdomain` and or
`endpoint`.

## Getting started
### Dependencies
- golangci
- pre-commit
- go - v1.23.3 or later

### Running the application

You can either build a binary using the instructions below or run it by running main.go with the required args

## Installation

Even whilst it is still WIP, if you wish to install it, run the below commands to move it into your `bin`, this assumes that you have downloaded the binary into you `Downloads` folder:

```bash
  cd ~/Downloads
  chmod +x dizzy
  mv ./dizzy /usr/local/bin/
```

Once this is done, you should be able to use the `dizzy` command, in your terminal.
