package battery

import "fmt"

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
