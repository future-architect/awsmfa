awsmfa (AWS MFA Refresher)
==========================

[![Test Status](https://github.com/d-tsuji/awsmfa/workflows/test/badge.svg?branch=master)][actions]
[![MIT License](http://img.shields.io/badge/license-MIT-blue.svg?style=flat-square)][license]

[actions]: https://github.com/d-tsuji/awsmfa/actions?workflow=test
[license]: https://github.com/d-tsuji/awsmfa/blob/main/LICENSE

`awsmfa` replaces the config(`~/.aws/config`) and credentials(`~/.aws/credentials`) for MFA in AWS.

## Usage

```console
$ awsmfa [options] [token-code]
```

### Options

```
--profile string
	The name of the profile from which the session can be obtained (default `default`)

--mfa-profile-name string
	MFA profile name (default `mfa`)

--duration-seconds int64
	Session expiration duration secounds (default `43200`)

--serial-number string
	AWS serial number. `--serial-number` is required

--token-code string
	Device token codes issued by the MFA. `--token-code` option or `token-code` is required

--quiet bool
	if enabled, log is not printed in the console. (default `false`)
```

### Example

```
$ awsmfa --serial-number arn:aws:iam::123456789012:mfa/d-tsuji --profile my-profile 123456
```

## Installation

- From binary

```
# binary
$ curl -sfL https://raw.githubusercontent.com/d-tsuji/awsmfa/master/install.sh | sudo sh -s -- -b /usr/local/bin
```

- From source code

```
# go get
$ go get github.com/d-tsuji/awsmfa/cmd/awsmfa
```
