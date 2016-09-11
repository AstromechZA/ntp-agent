package translation

import (
    "time"
)

const intMask = 0xFFFFFFFF
const era1900 = -2208988800
const era2036 = 2085978496
const nanosecondsPerSecond = 1e9

func ConvertNTPToNanoSeconds(input uint64) int64 {
    // get integral seconds
    seconds := int64((input >> 32) & intMask)

    // pull the MSB of the seconds
    msb := seconds & 0x80000000
    if msb == 0 {
        seconds += era2036
    } else {
        seconds += era1900
    }
    // multiply out to nano seconds range
    seconds *= nanosecondsPerSecond

    // isolate fractions
    fraction := int64(input & intMask)
    // convert into seconds
    extra := fraction * nanosecondsPerSecond / 0x100000000

    return seconds + extra
}

func ConvertNanoSecondsToNTP(input int64) uint64 {
    seconds := input / nanosecondsPerSecond
    nanos := input % nanosecondsPerSecond

    // multiply back into int range
    lowHalf := uint64(nanos * 0x100000000 / nanosecondsPerSecond) & intMask

    var highHalf uint64
    if seconds >= era2036 {
        highHalf = ((uint64(seconds - era2036) & intMask) | 0x80000000) << 32
    } else {
        highHalf = ((uint64(seconds - era1900) & intMask)) << 32
    }

    return highHalf | lowHalf
}

func ConvertNanoSecondsToTime(input int64) time.Time {
    return time.Unix(input / nanosecondsPerSecond, input % nanosecondsPerSecond)
}
func ConvertTimeToNanoSeconds(input time.Time) int64 {
    return input.UnixNano()
}

func ConvertNTPToTime(input uint64) time.Time {
    return ConvertNanoSecondsToTime(ConvertNTPToNanoSeconds(input))
}

func ConvertTimeToNTP(input time.Time) uint64 {
    return ConvertNanoSecondsToNTP(ConvertTimeToNanoSeconds(input))
}
