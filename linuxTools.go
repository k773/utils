package utils

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type LinuxHwTools struct {
	Logger LoggerO
}

// GetDiskSpace /dev/sda1      162420480 38199960 124204136  24% /  <- use "/" as mnt
func (*LinuxHwTools) GetDiskSpace(mnt ...string) (usage float64, free, used, total int64, e error) {
	var sb []byte
	if sb, e = exec.Command("bash", "-c", "df").CombinedOutput(); e == nil {
		for _, line := range strings.Split(string(sb), "\n") {
			if s := strings.Fields(line); len(s) == 6 {
				if Contains(mnt, s[5]) {
					a, _ := strconv.ParseInt(s[1], 10, 64) // total
					b, _ := strconv.ParseInt(s[2], 10, 64) // used
					c, _ := strconv.ParseInt(s[3], 10, 64) // available
					total += a * 1024
					used += b * 1024
					free += c * 1024
				}
			}

			if used != 0 {
				usage = float64(used) / float64(total)
			}
		}
	}
	return
}

// GetIOStats Parses info from /proc/pid/io:
// rchar: 4086576068
// wchar: 5667644758219
// syscr: 18960
// syscw: 619175514
// read_bytes: 5739536384
// write_bytes: 10860712153088
// cancelled_write_bytes: 41168896
func (*LinuxHwTools) GetIOStats(pid int) (m map[string]int64) {
	m = map[string]int64{}

	if err, lines := ReadFileByLines(fmt.Sprintf("/proc/%v/io", pid)); err == nil {
		for _, line := range lines {
			if s := strings.Split(line, " "); len(s) == 2 {
				m[s[0][:len(s[0])-1]], _ = strconv.ParseInt(s[1], 10, 64)
			}
		}
	}
	return
}

// CollectIOUsageStats Simple collects IO statistics for t time period
func (l *LinuxHwTools) CollectIOUsageStats(pid int, t time.Duration) (m map[string]int64) {
	a := l.GetIOStats(pid)
	time.Sleep(t)
	b := l.GetIOStats(pid)
	return subtractMaps(b, a)
}

func (*LinuxHwTools) GetProcPid(a string) (pid int, e error) {
	var sb []byte
	if sb, e = exec.Command("bash", "-c", "pidof "+a).CombinedOutput(); e == nil {
		pid, e = strconv.Atoi(strings.Trim(string(sb), "\n\r "))
	}
	return
}

// resp: m1-m2
func subtractMaps(m1, m2 map[string]int64) map[string]int64 {
	for k, v := range m1 {
		if v2, h := m2[k]; h {
			m1[k] = v - v2
		}
	}
	return m1
}

func (l *LinuxHwTools) GetCpuLoad(t time.Duration) (load float64, e error) {
	var i0, t0, i1, t1 uint64
	if i0, t0, e = l.GetCPUSample(); e == nil {
		time.Sleep(t)
		if i1, t1, e = l.GetCPUSample(); e == nil {
			t := t1 - t0
			load = float64(t-(i1-i0)) / float64(t)
		}
	}
	return
}

func (*LinuxHwTools) GetCPUSample() (idle, total uint64, e error) {
	var lines []string
	if e, lines = ReadFileByLines("/proc/stat"); e == nil {
		for _, line := range lines {
			if fields := strings.Fields(line); fields[0] == "cpu" {
				for i := 1; i < len(fields); i++ {
					val, _ := strconv.ParseUint(fields[i], 10, 64)
					total += val // tally up all the numbers to get total ticks
					if i == 4 {  // idle is the 5th field in the cpu line
						idle = val
					}
				}
				return
			}
		}
	}
	return
}

func (*LinuxHwTools) GetRamUsage() (usage float64, available, used, total int64, e error) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	var lines []string
	if e, lines = ReadFileByLines("/proc/meminfo"); e == nil {
	a:
		for _, line := range lines {
			if s := strings.Fields(line); len(s) == 3 {
				v, _ := strconv.ParseInt(s[1], 10, 64)
				switch s[0] {
				case "MemTotal:":
					total = v * 1024
					if available != 0 {
						break a
					}
				case "MemAvailable:":
					available = v * 1024
					if total != 0 {
						break a
					}
				default:
					//fmt.Println(line, s)
				}
			}
		}
		available = available + int64(m.HeapIdle)
		used = total - available
		usage = float64(used) / float64(total)
	} //

	return
}

func (*LinuxHwTools) GetFdCount() (usage float64, available, used, total int, e error) {
	var out []byte
	if out, e = exec.Command("/bin/sh", "-c", "ulimit -n").Output(); e == nil {
		total, e = strconv.Atoi(strings.Trim(string(out), "\n\r "))
		if out, e = exec.Command("/bin/sh", "-c", fmt.Sprintf("ls /proc/%v/fd/ | wc -l", os.Getpid())).Output(); e == nil {
			used, e = strconv.Atoi(strings.Trim(string(out), "\n\r "))
			available = total - used
			usage = float64(used) / float64(total)
		}
	}
	return
}

func (l *LinuxHwTools) CalculateTraffic(interface_ string, d time.Duration) (rx, tx float64, e error) {
	rx2, tx2, e := l.GetTraffic(interface_, d)
	return float64(rx2) / float64(d/time.Second), float64(tx2) / float64(d/time.Second), e
}

func (l *LinuxHwTools) GetTraffic(interface_ string, d time.Duration) (rx, tx int64, e error) {
	if rx, tx, e = l.GetTrafficSample(interface_); e == nil {
		time.Sleep(d)
		var rx2, tx2 int64
		if rx2, tx2, e = l.GetTrafficSample(interface_); e == nil {
			rx, tx = rx2-rx, tx2-tx
		}
	}
	return
}

func (*LinuxHwTools) GetTrafficSample(interface_ string) (rx, tx int64, e error) {
	var out []byte
	if out, e = exec.Command("/bin/sh", "-c", `ip -s -c link show `+interface_+` | sed -n -e 4p -e 6p |  cut -d " " -f5`).Output(); e == nil {
		if s := strings.Split(string(out), "\n"); len(s) >= 2 {
			var e1, e2 error
			rx, e1 = strconv.ParseInt(s[0], 10, 64)
			tx, e2 = strconv.ParseInt(s[1], 10, 64)
			e = JoinErrors(e1, e2)
		} else {
			e = errors.New("wrong output split")
		}
	}
	return
}

// GetPrimaryNetworkInterface returns the name of an interface through which traffic goes to 8.8.8.8
func (l *LinuxHwTools) GetPrimaryNetworkInterface() (iface string, e error) {
	var out []byte
	if out, e = exec.Command("/bin/sh", "-c", `ip route get 8.8.8.8 | sed -n 's/.*dev \([^\ ]*\).*/\1/p'`).Output(); e == nil {
		iface = string(out)
	}
	return
}
