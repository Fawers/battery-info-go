package battery

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

var ErrDefaultDeviceNotFound = errors.New("default BAT device not found")

type WattHour float32
type Watt float32
type Voltage float32
type Percentage float32

func (wh WattHour) String() string {
	return fmt.Sprintf("%.3f Wh", wh)
}

func (w Watt) String() string {
	return fmt.Sprintf("%.3f W", w)
}

func (v Voltage) String() string {
	return fmt.Sprintf("%.3f V", v)
}

func (p Percentage) String() string {
	return fmt.Sprintf("%.2f%%", p)
}

type History struct {
	Type   string
	Status string
	Time   time.Time
}

type Energy struct {
	Current, Empty, Full, FullDesign WattHour
	Percentage, Capacity             Percentage
}

type Battery struct {
	Present      bool
	Rechargeable bool
	State        string
	WarningLevel string
	Energy       Energy
	EnergyRate   Watt
	Voltage      Voltage
	ChargeCycles string
	TimeToEmpty  time.Duration
	TimeToFull   time.Duration
	Technology   string
	IconName     string
}

type BatteryInfo struct {
	SrcDevice   string
	NativePath  string
	Vendor      string
	Model       string
	Serial      uint64
	PowerSupply bool
	Updated     time.Time
	HasHistory  bool
	Statistics  bool
	Battery     Battery
	History     []History
}

func NewFromDevice(device string) (*BatteryInfo, error) {
	info, err := fetchDeviceInfo(device)

	if err != nil {
		return nil, err
	}

	bi := parseDeviceInfo(info)
	bi.SrcDevice = device
	return bi, nil
}

func NewFromDefaultDevice() (*BatteryInfo, error) {
	devs, err := listPowerDevices()

	if err != nil {
		return nil, err
	}

	for _, dev := range devs {
		if strings.LastIndex(dev, "BAT") != -1 {
			return NewFromDevice(dev)
		}
	}

	return nil, ErrDefaultDeviceNotFound
}

func (bi *BatteryInfo) Update() (*BatteryInfo, error) {
	return NewFromDevice(bi.SrcDevice)
}

func (bi *BatteryInfo) setStringField(field, value string) {
	switch field {
	case "native-path":
		bi.NativePath = value

	case "vendor":
		bi.Vendor = value

	case "model":
		bi.Model = value

	case "state":
		bi.Battery.State = value

	case "warning-level":
		bi.Battery.WarningLevel = value

	case "charge-cycles":
		bi.Battery.ChargeCycles = value

	case "technology":
		bi.Battery.Technology = value

	case "icon-name":
		bi.Battery.IconName = strings.TrimFunc(value, func(r rune) bool {
			return r == '\''
		})
	}
}

func (bi *BatteryInfo) setBoolField(field string, value bool) {
	switch field {
	case "power supply":
		bi.PowerSupply = value

	case "has history":
		bi.HasHistory = value

	case "has statistics":
		bi.Statistics = value

	case "present":
		bi.Battery.Present = value

	case "rechargeable":
		bi.Battery.Rechargeable = value
	}
}

func (bi *BatteryInfo) setFloatValue(field string, value float32) {
	switch field {
	case "energy":
		bi.Battery.Energy.Current = WattHour(value)

	case "energy-empty":
		bi.Battery.Energy.Empty = WattHour(value)

	case "energy-full":
		bi.Battery.Energy.Full = WattHour(value)

	case "energy-full-design":
		bi.Battery.Energy.FullDesign = WattHour(value)

	case "energy-rate":
		bi.Battery.EnergyRate = Watt(value)

	case "voltage":
		bi.Battery.Voltage = Voltage(value)

	case "percentage":
		bi.Battery.Energy.Percentage = Percentage(value)

	case "capacity":
		bi.Battery.Energy.Capacity = Percentage(value)
	}
}
