package utils

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type LinuxHwTools struct {
}

func (*LinuxHwTools) GetIOStats(pid int) (m map[string]int64) {
	m = map[string]int64{}

	if err, lines := ReadFileByLines(fmt.Sprintf("/proc/%v/io", pid)); err == nil {
		for _, line := range lines {
			if s := strings.Split(line, " "); len(s) == 2 {
				m[s[0]], _ = strconv.ParseInt(s[1], 10, 64)
			}
		}
	}
	return
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
					fmt.Println(line, s)
				}
			}
		}
		used = total - available
		usage = float64(used) / float64(total)
	}
	return
}
