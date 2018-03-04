# Helm Local Chart Version Plugin

This is a Helm plugin that helps you manipulate your chart version.
It is meant to be used in integration pipelines.

## Usage

see [the documentation](docs/local-chart-version.md)

## Install

### Using Helm plugin manager (> 2.3.x)

```shell
helm plugin install https://github.com/mbenabda/helm-local-chart-version
```

### Pre Helm 2.3.0 Installation
Pick a release tarball from the [releases](https://github.com/mbenabda/helm-local-chart-version/releases) page.

Unpack the tarball in your helm plugins directory (`$(helm home)/plugins`).

E.g.
```
curl -L $TARBALL_URL | tar -C $(helm home)/plugins -xzv
```

## Build

Clone the repository into your `$GOPATH` and then build it.

```
$ mkdir -p $GOPATH/src/github.com/mbenabda/
$ cd $GOPATH/src/github.com/mbenabda/
$ git clone https://github.com/mbenabda/helm-local-chart-version.git
$ cd helm-local-chart-version
$ make install
```

The above will install this plugin into your `$HELM_HOME/plugins` directory.

### Prerequisites

- You need to have [Go](http://golang.org) installed. Make sure to set `$GOPATH`
- If you don't have [Glide](http://glide.sh) installed, this will install it into
  `$GOPATH/bin` for you.