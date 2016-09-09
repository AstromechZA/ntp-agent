package main

import (
    "os"
    "fmt"
    "flag"
    "errors"
    "net"
)

const usageString =
`ntp-taste is a simple binary for demonstrating an NTP relationships

`

const ntpPort = 123

func mainInner() error {

    // set a more verbose usage message.
    flag.Usage = func() {
        os.Stderr.WriteString(usageString)
        flag.PrintDefaults()
    }

    // parse them
    flag.Parse()

    // expect at least one time server
    if flag.NArg() == 0 {
        return errors.New("Expected at least one upstream NTP server as an argument.")
    }

    svrAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", flag.Args()[0], ntpPort))
    if err != nil { return err }

    conn, err := net.DialUDP("udp", nil, svrAddr)
    if err != nil { return err }

    defer conn.Close()

    buf := make([]byte, 48)
    buf[0] = 0x1B
    _, err = conn.Write(buf)
    if err != nil { return err }

    inbuf := make([]byte, 1024)
    n, addr, err := conn.ReadFromUDP(inbuf)

    fmt.Println("UDP Server", addr)
    fmt.Println("Received", n, "bytes")
    fmt.Printf("Bytes %x\n", inbuf[:n])

    return err
}

func main() {
    if err := mainInner(); err != nil {
        os.Stderr.WriteString(err.Error() + "\n")
        os.Exit(1)
    }
}
