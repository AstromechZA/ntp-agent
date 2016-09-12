package main

import (
    "os"
    "fmt"
    "flag"
    "errors"
    "time"
    "net"

    "github.com/AstromechZA/ntp-agent/packet"
    "github.com/AstromechZA/ntp-agent/translation"
    "github.com/AstromechZA/ntp-agent/constants"
)

const usageString =
`ntp-agent is a simple binary for demonstrating an NTP relationships

`

const ntpPort = 123

func getNTPPacket(server string) (*packet.RawPacket, error) {
    svrAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", server, ntpPort))
    if err != nil { return nil, err }

    conn, err := net.DialUDP("udp", nil, svrAddr)
    if err != nil { return nil, err }

    defer conn.Close()

    h := &packet.RawPacket{
        Version: 4,
        Mode: constants.ModeClient,
        OriginateTimestamp: translation.ConvertTimeToNTP(time.Now()),
    }
    buf, err := h.ToSlice()
    if err != nil { return nil, err }
    _, err = conn.Write(*buf)
    if err != nil { return nil, err }

    inbuf := make([]byte, 1024)
    n, _, err := conn.ReadFromUDP(inbuf)
    if err != nil { return nil, err }
    headerContent := inbuf[:n]

    return packet.ParseRaw(&headerContent)
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

        t1 := time.Now()
        h, err := getNTPPacket(server)
        if err != nil { return err }
        t4 := time.Now()
        t2 := translation.ConvertNTPToTime(h.ReceiveTimestamp)
        t3 := translation.ConvertNTPToTime(h.TransmitTimestamp)

        offset := (t2.Sub(t1) - t4.Sub(t3)) / 2
        delay := t4.Sub(t1) - t3.Sub(t2)

        fmt.Println(h)
        fmt.Println(offset, delay)
        fmt.Println(translation.ConvertNTPToTime(h.ReceiveTimestamp))
    }

    return nil
}

func main() {
    if err := mainInner(); err != nil {
        os.Stderr.WriteString(err.Error() + "\n")
        os.Exit(1)
    }
}
