package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/AstromechZA/ntp-agent/constants"
	"github.com/AstromechZA/ntp-agent/packet"
	"github.com/AstromechZA/ntp-agent/translation"
)

const usageString = `ntp-agent is a simple binary for pulling and setting a
more accurate time.

Although not as accurate as true NTP, it may be effective enough for some
use cases.

Given a number of remote NTP servers, this application will calculate an
average clock offset and if you approve, set the current date and time
accordingly.

See www.ntp.org for a list of useful ntp servers to pull from.

`

// Version is the version string
// format should be 'X.YZ'
// Set this at build time using the -ldflags="-X main.Version=X.YZ"
var Version = "<unofficial build>"

const ntpPort = 123

func getNTPPacket(server string) (*packet.RawPacket, error) {
	svrAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", server, ntpPort))
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, svrAddr)
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	h := &packet.RawPacket{
		Version:            4,
		Mode:               constants.ModeClient,
		OriginateTimestamp: translation.ConvertTimeToNTP(time.Now()),
	}
	buf, err := h.ToSlice()
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(*buf)
	if err != nil {
		return nil, err
	}

	inbuf := make([]byte, 1024)
	conn.SetDeadline(time.Now().Add(1 * time.Second))
	n, _, err := conn.ReadFromUDP(inbuf)
	if err != nil {
		return nil, err
	}
	headerContent := inbuf[:n]

	return packet.ParseRaw(&headerContent)
}

func mainInner() error {
	// some flag args
	assumeYesFlag := flag.Bool("assume-yes", false, "Don't prompt for sync")
    versionFlag := flag.Bool("version", false, "Print the version string")

	// set a more verbose usage message.
	flag.Usage = func() {
		os.Stderr.WriteString(usageString)
		flag.PrintDefaults()
	}

	// parse them
	flag.Parse()

    // check for the version option
    if *versionFlag {
        fmt.Println("Version: " + Version)
        fmt.Println("Project: https://github.com/AstromechZA/ntp-agent")
        return nil
    }

	// expect at least one time server
	if flag.NArg() == 0 {
		return errors.New("Expected at least one upstream NTP server as an argument.")
	}

	offsetTimes := []time.Duration{}
	for _, server := range flag.Args() {

		// t1 == time we sent the request
		t1 := time.Now()
		h, err := getNTPPacket(server)
		if err != nil {
			return err
		}
		// t4 == time we received the response
		t4 := time.Now()
		// t2 == time server received our packet
		t2 := translation.ConvertNTPToTime(h.ReceiveTimestamp)
		// t3 == time server sent its response
		t3 := translation.ConvertNTPToTime(h.TransmitTimestamp)

		// time difference between client and server
		offset := (t2.Sub(t1) - t4.Sub(t3)) / 2
		// round trip delay
		// delay := t4.Sub(t1) - t3.Sub(t2)

		// add to offset list
		offsetTimes = append(offsetTimes, offset)
	}

	totalOffset := offsetTimes[0]
	for i := 1; i < len(offsetTimes); i++ {
		totalOffset += offsetTimes[i]
	}
	avgOffset := totalOffset / time.Duration(len(offsetTimes))

	fmt.Printf("Clock offset seems to be about %s\n", avgOffset)
	first := time.Now()
	fixed := first.Add(avgOffset)
	fmt.Printf("This would change the current time %s -> %s\n", first, fixed)

	if avgOffset > (time.Duration(60) * time.Minute) {
		return fmt.Errorf("Refusing to change time by more than 1 hour. You'll need to fix this manually.")
	}

	assumeYes := *assumeYesFlag
	if !assumeYes {
		fmt.Println("Is this ok? (yes/no)")
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		text = strings.Trim(text, " \n\t")
		if text != "yes" {
			fmt.Println("Not setting time.")
			return nil
		}
	}
	fmt.Println("Attempting to set time..")
	first = time.Now()
	fixed = first.Add(avgOffset)

	cmd := exec.Command("date", fixed.Format("010215042006.05"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				return fmt.Errorf("Failed with code %d %s\n", status.ExitStatus(), err.Error())
			}
		} else {
			return fmt.Errorf("Failed with code 127 %s\n", err.Error())
		}
	}

	return nil
}

func main() {
	if err := mainInner(); err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
