# Haptic

## Installation

```
go get nanocloud.com/zeroinstall/haptic
```

## Configuration

If you want to develop on **Haptic**, we recommand you to export **NANOCONF**
environment variable in order to use your own version of configuration file, for
example:

```
NANOCONF="config-local.json"
go build
./haptic serve
```
