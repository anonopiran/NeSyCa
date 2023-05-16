package config

import (
	"fmt"
	"math/rand"
	net "net"
	"time"

	"strconv"
)

type LogLevelType string
type TCPAddressType struct {
	Plain   string
	Address net.TCPAddr
}
type SizeType int
type SecondDuration time.Duration
type MiliSecondDuration time.Duration
type RateType struct {
	Min SizeType `env:"MIN,required"`
	Max SizeType `env:"MAX,required"`
}
type SettingsType struct {
	TCPTargets    []TCPAddressType   `env:"TCP_TARGETS,required"`
	Rate          RateType           `envPrefix:"RATE_"`
	ReportTick    SecondDuration     `env:"REPORT_TICK" envDefault:"10"`
	UpdateTick    SecondDuration     `env:"UPDATE_TICK" envDefault:"60"`
	DrainDelay    MiliSecondDuration `env:"DRAIN_DELAY" envDefault:"5"`
	ConnPerTarget int                `env:"CONN_PER_TARGET" envDefault:"5"`
	LogLevel      LogLevelType       `env:"LOG_LEVEL" envDefault:"warning"`
}

func (f *SizeType) Hr() string {
	const unit = 1024
	if *f < unit {
		return fmt.Sprintf("%d B", *f)
	}
	div, exp := int64(unit), 0
	for n := *f / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(*f)/float64(div), "KMGTPE"[exp])
}
func parseDur(s []byte) (time.Duration, error) {
	val := string(s)
	valInt, err := strconv.Atoi(val)
	if err != nil {
		return 0, err
	}
	return time.Duration(valInt), nil
}
func (f *RateType) GetRandom() SizeType {
	return SizeType(rand.Intn(int(f.Max-f.Min)) + int(f.Min))
}
