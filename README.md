# ntp-agent

I wanted to learn how NTP clients work, so I'm writing my own in Golang.

See the beginning of this story here: http://astromechza.github.io/2016/09/10/ntp-agent-part-1.html

## commit db13b50

```
$ ntp-agent --help
ntp-agent is a simple binary for pulling and setting a
more accurate time.

Although not as accurate as true NTP, it may be effective enough for some
use cases.

Given a number of remote NTP servers, this application will calculate an
average clock offset and if you approve, set the current date and time
accordingly.

See www.ntp.org for a list of useful ntp servers to pull from.

  -assume-yes
        Don't prompt for sync
  -version
        Print the version string

$ ntp-agent --version
Version: 0.1 (commit cf720b7 @ 2016-09-25)
Project: https://github.com/AstromechZA/ntp-agent
```

At this point, the binary can be used to set the date and time. It does it in
a fairly simple manner, none of the complex clock selection, clustering,
etc. But it works pretty effectively so far.
