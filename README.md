# Sms Service

This is the Sms service

Generated with

```
micro new github.com/xionglongjun/micro-mall/sms/srv --namespace=go.micro --type=srv
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.srv.sms
- Type: srv
- Alias: sms

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./github.com/xionglongjun/micro-mall/sms/srv
```

Build a docker image
```
make docker
```