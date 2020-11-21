# VDF
![Status Shield](https://img.shields.io/badge/status-release-brightgreen?style=for-the-badge)
[![Test Status](https://img.shields.io/github/workflow/status/wakeful-cloud/vdf/Tests?label=Tests&style=for-the-badge&logo=github-actions)](https://github.com/wakeful-cloud/vdf/actions)

Golang Binary **Valve Data Format** implementation

## Install
```console
go get -u github.com/wakeful-cloud/vdf
```

## Features
* Binary VDF support
* Fully unit tested
* No dependencies

## Limitations
* No support for non-binary VDF's
* Order of key-value's are not preserved (Steam doesn't care though)

## Docs

* `vdf.Map`: structure use to represent a parsed/read VDF
  * Signature: `map[string]interface{}`

* `vdf.ReadVdf`: function that reads bytes to a `vdf.Map`
  * Signature: `ReadVdf([]byte): (vdf.Map, error)`

* `vdf.WriteVdf`: function that writes a `vdf.Map` to bytes
  * Signature: `WriteVdf(vdf.Map): ([]byte, error)`

## Example
See [example/main.go](./example/main.go)

## Credit
Heavily based on [Corecii's Steam Binary VSF TS Package](https://github.com/Corecii/steam-binary-vdf-ts).