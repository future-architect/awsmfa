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
	AWS serial number. `--serial-number` or `--account-id` and `--role` are required

--account-id string
	AWS account

--role string
	AWS role

--token-code string
	Device token codes issued by the MFA

--quiet bool
    if enabled, log is not printed in the console
```

### Example

```
$ awsmfa --serial-number arn:aws:iam::123456789012:mfa/d-tsuji --profile my-profile 123456
```

## Installation

```
# go get
$ go get github.com/d-tsuji/awsmfa/cmd/awsmfa
```
