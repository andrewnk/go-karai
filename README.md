# [go-karai]

Karai connection helper written in Go

## Usage

`./go-karai`

This will launch `go-karai`

`./go-karai /ip4/127.0.0.1/tcp/40485/p2p/QmZCddScZ82V8Y2dkU5o7Fvs8QLuFm1TbKepQWvrDGktgR`

With two instances of `go-karai` running, you can send a ping from node to node like this to verify that you're connected.

## Dependencies

-   Golang 1.10+

## Building

`git clone https://github.com/rocksteadytc/go-karai`

Clone the repository

`go mod init github.com/rocksteadytc/go-karai`

First run only: Initialize the go module

`GOPRIVATE='github.com/libp2p/*' go get ./...`

First run only: Look for available releases

`go build`

Compile to produce a binary `go-karai`

`go build -gcflags="-e" && ./go-karai`

**Optional:** Compile with all errors displayed, then run binary. Avoids "too many errors" from hiding error info.

## Contributing

-   `gofmt` is used on all files.
-   go modules are used to manage dependencies.
