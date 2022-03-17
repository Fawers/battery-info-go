package battery

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func runCmd(cmd *exec.Cmd) ([]string, error) {
	out, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	content := strings.Split(
		strings.TrimSpace(string(out)), "\n")

	return content, nil
}

func listPowerDevices() ([]string, error) {
	cmd := exec.Command("upower", "-e")
	return runCmd(cmd)
}

func fetchDeviceInfo(device string) ([]string, error) {
	cmd := exec.Command("upower", "-i", device)
	return runCmd(cmd)
}

func parseDeviceInfo(info []string) *BatteryInfo {
	binfo := new(BatteryInfo)
	binfo.History = make([]History, 0)
	var historyType string

	for _, line := range info {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "History") {
			j := strings.IndexRune(line, '(') + 1
			k := strings.IndexRune(line, ')')
			historyType = line[j:k]

		} else if i := strings.IndexRune(line, ':'); i != -1 {
			field := line[:i]

			switch field {
			case "native-path", "vendor", "model", "state", "warning-level":
				fallthrough

			case "charge-cycles", "technology", "icon-name":
				j := strings.LastIndex(line, " ")
				binfo.setStringField(field, line[j+1:])

			case "power supply", "has history", "has statistics", "present", "rechargeable":
				j := strings.LastIndex(line, " ")
				b := false

				if line[j+1:] == "yes" {
					b = true
				}

				binfo.setBoolField(field, b)

			case "serial":
				j := strings.LastIndex(line, " ")
				valueStr := line[j+1:]
				value, _ := strconv.ParseUint(valueStr, 10, 64)
				binfo.Serial = value

			case "energy", "energy-empty", "energy-full", "energy-full-design":
				fallthrough

			case "energy-rate", "voltage":
				k := strings.LastIndex(line, " ")
				line = line[:k]
				fallthrough

			case "percentage", "capacity":
				line = strings.TrimRight(line, "%")
				j := strings.LastIndex(line, " ")
				valueStr := line[j+1:]
				value, _ := strconv.ParseFloat(valueStr, 32)

				binfo.setFloatValue(field, float32(value))

			case "time to empty", "time to full":
				k := strings.LastIndex(line, " ")
				unit := line[k+1]
				line = line[:k]
				j := strings.LastIndex(line, " ") + 1
				time_ := line[j:]

				dur, _ := time.ParseDuration(fmt.Sprintf("%s%c", time_, unit))

				if strings.HasSuffix(field, "empty") {
					binfo.Battery.TimeToEmpty = dur
				} else {
					binfo.Battery.TimeToFull = dur
				}

			case "updated":
				k := strings.LastIndex(line, "(") - 1
				line = line[:k]
				j := strings.LastIndex(line, "  ") + 2
				dateStr := line[j:]

				date, _ := time.Parse("Mon 2 Jan 2006 03:04:05 PM MST", dateStr)
				binfo.Updated = date
			}
		} else if historyParts := strings.Split(line, "\t"); len(historyParts) == 3 {
			timestampStr, _ := strconv.Atoi(historyParts[0])
			timestamp := time.Unix(int64(timestampStr), 0)
			status := historyParts[2]
			binfo.History = append(binfo.History, History{
				Type:   historyType,
				Time:   timestamp,
				Status: status,
			})
		}
	}

	return binfo
}
