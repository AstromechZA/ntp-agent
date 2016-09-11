package translation

import (
    "testing"
    "time"
)

func absDint(v1 int64, v2 int64) int64 {
    if v1 < v2 { return v2 - v1 }
    return v1 - v2
}

func absDuint(v1 uint64, v2 uint64) uint64 {
    if v1 < v2 { return v2 - v1 }
    return v1 - v2
}

func TestTimeToNanoSecondsToTime(t *testing.T) {
    t1 := time.Now()
    s := ConvertTimeToNanoSeconds(t1)
    t2 := ConvertNanoSecondsToTime(s)
    if t1 != t2 {
        t.Error("Timestamps were not equal")
    }
}

const ntpTestVal1 = uint64(0xdb7f8bf3f9fb5658)

func TestConvertNTPToNanoSecondsToNTP(t *testing.T) {
    ns := ConvertNTPToNanoSeconds(ntpTestVal1)
    ti := ConvertNanoSecondsToTime(ns)
    ns2 := ConvertTimeToNanoSeconds(ti)
    ntp := ConvertNanoSecondsToNTP(ns2)

    d := absDuint(ntp, ntpTestVal1)
    if d > 10 { t.Error("NTP conversion was too far out") }
}

func TestConvertNanoSecondsToNTPToNanoSeconds(t *testing.T) {
    ns := time.Now().UnixNano()
    ntp := ConvertNanoSecondsToNTP(ns)
    ns2 := ConvertNTPToNanoSeconds(ntp)
    d := absDint(ns, ns2)
    if d > 10 { t.Error("NTP conversion was too far out") }
}
