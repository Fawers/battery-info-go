package battery

import (
	"strings"
	"time"
)

// History as listed by upower -i.
type History struct {
	Type   string
	Status string
	Time   time.Time
}

// Energy holds the information of the energy
// of the battery, as well as its current percentage
// and capacity.
type Energy struct {
	Current, Empty, Full, FullDesign WattHour
	Percentage, Capacity             Percentage
}

// Battery holds the specific information of the battery,
// such as its current state, if it's rechargeable, its
// voltage, technology used, and time to empty or full.
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

// BatteryInfo holds the general information of the battery,
// such as its model, vendor, serial, among others.
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

// NewForDevice constructs a new *BatteryInfo for the
// given device. The device can be acquired with ListPowerDevices.
// Alternatively, if something wrong happens when fetching or
// parsing the battery information, an error may be returned.
//
// The following error types may be returned:
// *InvalidDeviceError,
// *CommandError.
func NewForDevice(device string) (*BatteryInfo, error) {
	info, err := fetchDeviceInfo(device)

	if err != nil {
		return nil, err
	}

	bi := parseDeviceInfo(info)
	bi.SrcDevice = device
	return bi, nil
}

// NewForDefaultDevice returns a *BatteryInfo for the default
// device, which is the one that ends with BAT#, where # is a
// number.
//
// The following error types may be returned:
// *DefaultDeviceNotFoundError,
// *InvalidDeviceError,
// *CommandError.
func NewForDefaultDevice() (*BatteryInfo, error) {
	devs := ListPowerDevices()

	for _, dev := range devs {
		if strings.LastIndex(dev, "BAT") != -1 {
			return NewForDevice(dev)
		}
	}

	return nil, &DefaultDeviceNotFoundError{}
}

// Update returns a new *BatteryInfo by re-parsing the battery
// information as done by NewForDevice.
//
// Possible returnable errors are the same as for NewForDevice().
func (bi *BatteryInfo) Update() (*BatteryInfo, error) {
	return NewForDevice(bi.SrcDevice)
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
