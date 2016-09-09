package main

import (
    "os"
    "fmt"
    "flag"
    "errors"
    "net"

    "github.com/AstromechZA/ntp-taste/header"
    "github.com/AstromechZA/ntp-taste/constants"
)

const usageString =
`ntp-taste is a simple binary for demonstrating an NTP relationships

`

const ntpPort = 123

func getNTPHeader(server string) (*header.RawHeader, error) {
    svrAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", server, ntpPort))
    if err != nil { return nil, err }

    conn, err := net.DialUDP("udp", nil, svrAddr)
    if err != nil { return nil, err }

    defer conn.Close()

    h := &header.RawHeader{Version: 3, Mode: constants.ModeClient}
    buf, err := h.ToSlice()
    if err != nil { return nil, err }
    _, err = conn.Write(*buf)
    if err != nil { return nil, err }

    inbuf := make([]byte, 1024)
    n, _, err := conn.ReadFromUDP(inbuf)
    if err != nil { return nil, err }
    headerContent := inbuf[:n]
    return header.ParseRaw(&headerContent)
}


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

    for _, server := range flag.Args() {
        h, err := getNTPHeader(server)
        if err != nil { return err }

        fmt.Println(header.ConvertNTPToTime(h.ReceiveTimestamp), header.ConvertNTPToTime(h.TransmitTimestamp))
    }

    return nil
}

func main() {
    if err := mainInner(); err != nil {
        os.Stderr.WriteString(err.Error() + "\n")
        os.Exit(1)
    }
}
