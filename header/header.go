package header

import (
    "errors"
    "fmt"
)

type RawHeader struct {
    Leap byte
    Version byte
    Mode byte
    Stratum byte
    Poll byte
    Precision byte
    RootDelay int32
    RootDispersion int32
    ReferenceID int32
    ReferenceTimestamp int64
    OriginateTimestamp int64
    ReceiveTimestamp int64
    TransmitTimestamp int64
}


/*
2 bits --- LI
3 bits --- VN
3 bits --- Mode
1 byte --- Stratum
1 byte --- Poll
1 byte --- Precision

4 byte --- Root delay
4 byte --- Root disperse
4 byte --- ref id
8 byte --- ref time
8 byte --- origin time
8 byte --- recv time
8 byte --- trasmit time
*/

func (h *RawHeader) ToSlice() (*[]byte, error) {
    output := make([]byte, 48)
    err := h.IntoSlice(&output)
    return &output, err
}

func (h *RawHeader) IntoSlice(output *[]byte) error {

    outputData := *output

    var b1 byte
    b1 = b1 | ((h.Leap & 0x2) << 6)
    b1 = b1 | ((h.Version & 0x7) << 3)
    b1 = b1 | (h.Mode & 0x7)
    outputData[0] = b1

    outputData[1] = h.Stratum
    outputData[2] = h.Poll
    outputData[3] = h.Precision

    if err := putInt32(output, 4, h.RootDelay); err != nil { return err }
    if err := putInt32(output, 8, h.RootDispersion); err != nil { return err }
    if err := putInt32(output, 12, h.ReferenceID); err != nil { return err }
    if err := putInt64(output, 16, h.ReferenceTimestamp); err != nil { return err }
    if err := putInt64(output, 24, h.OriginateTimestamp); err != nil { return err }
    if err := putInt64(output, 32, h.ReceiveTimestamp); err != nil { return err }
    if err := putInt64(output, 40, h.TransmitTimestamp); err != nil { return err }

    return nil
}

func ParseRaw(input *[]byte) (*RawHeader, error) {

    inputData := *input

    // check incoming data length
    if len(inputData) != 48 {
        return nil, errors.New("Incoming packet must be 48 bytes")
    }

    // build output structure
    output := RawHeader{}

    // first byte
    b1 := inputData[0]
    output.Leap = (b1 >> 6) & 0x2
    output.Version = (b1 >> 3) & 0x7
    output.Mode = b1 & 0x7

    // next 3 bytes
    output.Stratum = inputData[1]
    output.Poll = inputData[2]
    output.Precision = inputData[3]

    // remaining components
    var err error
    output.RootDelay, err = getInt32(input, 4)
    if err != nil { return nil, err }
    output.RootDispersion, err = getInt32(input, 8)
    if err != nil { return nil, err }
    output.ReferenceID, err = getInt32(input, 12)
    if err != nil { return nil, err }
    output.ReferenceTimestamp, err = getInt64(input, 16)
    if err != nil { return nil, err }
    output.OriginateTimestamp, err = getInt64(input, 24)
    if err != nil { return nil, err }
    output.ReceiveTimestamp, err = getInt64(input, 32)
    if err != nil { return nil, err }
    output.TransmitTimestamp, err = getInt64(input, 40)
    if err != nil { return nil, err }

    return &output, nil
}

func getInt32(input *[]byte, position int) (int32, error) {
    inputData := *input

    if len(inputData) < position + 4 {
        return 0, fmt.Errorf("Buffer overflow cannot get 4 bytes @ %d", position)
    }

    var a int32
    for i := 0; i < 4; i++ {
        a = (a << 8) | int32(inputData[position + i])
    }
    return a, nil
}

func putInt32(output *[]byte, position int, data int32) error {
    outputData := *output

    if len(outputData) < position + 4 {
        return fmt.Errorf("Buffer overflow cannot put 8 bytes @ %d", position)
    }

    outputData[position] = byte((data >> 24) & 0xFF)
    outputData[position + 1] = byte((data >> 16) & 0xFF)
    outputData[position + 2] = byte((data >> 8) & 0xFF)
    outputData[position + 3] = byte(data & 0xFF)
    return nil
}

func getInt64(input *[]byte, position int) (int64, error) {
    inputData := *input

    if len(inputData) < position + 8 {
        return 0, fmt.Errorf("Buffer overflow cannot get 4 bytes @ %d", position)
    }

    var a int64
    for i := 0; i < 8; i++ {
        a = (a << 8) | int64(inputData[position + i])
    }
    return a, nil
}

func putInt64(output *[]byte, position int, data int64) error {
    outputData := *output

    if len(outputData) < position + 4 {
        return fmt.Errorf("Buffer overflow cannot put 8 bytes @ %d", position)
    }

    outputData[position] = byte((data >> 56) & 0xFF)
    outputData[position + 1] = byte((data >> 48) & 0xFF)
    outputData[position + 2] = byte((data >> 40) & 0xFF)
    outputData[position + 3] = byte((data >> 32) & 0xFF)
    outputData[position + 4] = byte((data >> 24) & 0xFF)
    outputData[position + 5] = byte((data >> 16) & 0xFF)
    outputData[position + 6] = byte((data >> 8) & 0xFF)
    outputData[position + 7] = byte(data & 0xFF)
    return nil
}
