awsmfa (AWS MFA Refresher)
==========================

[![Test Status](https://github.com/future-architect/awsmfa/workflows/test/badge.svg?branch=master)][actions]
[![Apache-2.0 license](https://img.shields.io/badge/license-Apache2.0-blue.svg)][license]

[actions]: https://github.com/future-architect/awsmfa/actions?workflow=test
[license]: https://github.com/future-architect/awsmfa/blob/master/LICENSE

`awsmfa` replaces the config and credentials for MFA in AWS.

## Usage

```console
$ awsmfa [options] [token-code]
```

- The config file uses `~/.aws/config` by default. You can override it with the environment variable `AWS_CONFIG_FILE`.
- The credentials file uses `~/.aws/credentials` by default. You can override it with the environment variable `AWS_SHARED_CREDENTIALS_FILE`.

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

--quiet bool
	if enabled, log is not printed in the console. (default `false`)
```

### Example

```
$ awsmfa --serial-number arn:aws:iam::123456789012:mfa/future-architect --profile my-profile 123456
2021/02/28 11:01:49 Wrote session token for profile mfa, expiration: 2021-02-28 14:01:49.689 +0000 UTC
```

#### initial

- `.aws/config`

<table>
<thead><tr><th>Before</th><th>After</th></tr></thead>
<tbody>
<tr><td valign="top">

```ini
[default]
region = us-east-1
output = json
```

</td><td valign="top">

```ini
[default]
region = us-east-1
output = json

[profile mfa]

```
</td></tr>
</tbody></table>

- `.aws/credentials`

<table>
<thead><tr><th>Before</th><th>After</th></tr></thead>
<tbody>
<tr><td valign="top">

```ini
[default]
aws_access_key_id = ABCDEFGHIJKLMNOPQRST
aws_secret_access_key = ChcdJbC9kraRNW5iy8XgDyR4QNRT44kKRPmKEGQT
```

</td><td valign="top">

```ini
[default]
aws_access_key_id     = ABCDEFGHIJKLMNOPQRST
aws_secret_access_key = ChcdJbC9kraRNW5iy8XgDyR4QNRT44kKRPmKEGQT

[mfa]
aws_access_key_id     = AKIAIOSFODNN7EXAMPLE
aws_secret_access_key = wJalrXUtnFEMI/K7MDENG/bPxRfiCYzEXAMPLEKEY
aws_session_token     = AQoEXAMPLEH4aoA....

```
</td></tr>
</tbody></table>


#### next time

Update `aws_access_key_id`, `aws_secret_access_key` and `aws_session_token` in the target section (default is `mfa`) of `.aws/credentials`.

## Installation

- From binary

```
# binary
$ curl -sfL https://raw.githubusercontent.com/future-architect/awsmfa/master/install.sh | sudo sh -s -- -b /usr/local/bin
```

- From source code

```
# go get
$ go get github.com/future-architect/awsmfa/cmd/awsmfa
```
