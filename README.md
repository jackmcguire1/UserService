# User Service

[![Build Status](https://travis-ci.org/jackmcguire1/UserService.svg?branch=main)](hhttps://travis-ci.org/jackmcguire1/UserService)
[![Go Report Card](https://goreportcard.com/badge/github.com/jackmcguire1/UserService)](https://goreportcard.com/report/github.com/jackmcguire1/UserService)

[git]:    https://git-scm.com/
[golang]: https://golang.org/
[dlv]:    https://github.com/go-delve/delve
[modules]: https://github.com/golang/go/wiki/Modules
[goLand]: https://www.jetbrains.com/go/
[golint]: https://github.com/golangci/golangci-lint
[aws-cli]: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-install.html
[aws-cli-config]: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html
[aws-sam-cli]: https://github.com/awslabs/aws-sam-cli
[localstack]: https://github.com/localstack/localstack


## ABOUT
> This repo contains a go module that exposes a User Microservice

### Prerequisites

- [Git][git]
- [Go 1.16][golang]+
- [golangCI-Lint][golint]
- [Delve Debugger][dlv]

<br>

### [golangCI-Lint][golint]
```shell
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $GOPATH/bin latest
```
<br>

```Shell 
golangci-lint run
```

### [Delve Debugger][dlv]
```shell
GOARCH=amd64 GOOS=linux go build -o ./dlv github.com/go-delve/delve/cmd/dlv
```

### [AWS CLI Configuration][aws-cli-config]
> Make sure you configure the AWS CLI
- AWS Access Key ID
- AWS Secret Access Key
- Default region 'us-east-1'
```shell
aws configure
```

## Thanks

This project exists thanks to **ME!**.

## Donations
All donations are appreciated!

[![Donate](https://img.shields.io/badge/Donate-PayPal-green.svg)](http://paypal.me/crazyjack12)