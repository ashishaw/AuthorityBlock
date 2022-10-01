# Authority Block

A general purpose blockchain inspired by Nano and VeChain and highly compatible with Ethereum's ecosystem.

This is the first implementation written in golang.

[![Go](https://img.shields.io/badge/golang-%3E%3D1.16-orange.svg)](https://golang.org)

## Installation

### Requirements

Ablock requires `Go` 1.16+ and `C` compiler to build. To install `Go`, follow this [link](https://golang.org/doc/install).

### Getting the source

Clone this repo:

```shell
git clone https://github.com/ashishaw/authorityblock.git
cd authorityblock
```

### Building

To build the main app `ablock`, just run

```shell
make
```

or build the full suite:

```shell
make all
```

If no errors are reported, all built executable binaries will appear in folder *bin*.

## Running

Connect Development network

```shell
bin/ablock --network genesis/genesis.json
```

### Sub-commands

* `solo`                client runs in solo mode for test & dev

```shell
# create new block when there is pending transaction
bin/ablock solo --on-demand

# save blockchain data to disk(default to memory)
bin/ablock solo --persist

# two options can work together
bin/ablock solo --persist --on-demand
```

* `master-key`          master key management

```shell
# print the master address
bin/ablock master-key

# export master key to keystore
bin/ablock master-key --export > keystore.json


# import master key from keystore
cat keystore.json | bin/ablock master-key --import
```

## Acknowledgement

A special shout out to following projects:

* [Ethereum](https://github.com/ethereum)
* [VeChain](https://github.com/vechain)
* [Nano](https://github.com/nanocurrency)
* [Swagger](https://github.com/swagger-api)

## License

Authority Block is licensed under the
[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.html)