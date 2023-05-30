package config

import (
	"fmt"
	net "net"
	"strconv"
	"time"

	"gonum.org/v1/gonum/stat/distuv"
)

type LogLevelType string
type TCPAddressType struct {
	Plain   string
	Address net.TCPAddr
}
type SizeType int
type SecondDuration time.Duration
type MiliSecondDuration time.Duration
type DistType string
type uvType interface {
	Rand() float64
}

const (
	UniformE   DistType = "uniform"
	LogNormalE DistType = "log-normal"
)

type RateType struct {
	Dist         DistType `env:"DISTRIBUTION,required"`
	UniformMin   float64  `env:"UNIFORM_MIN"`
	UniformMax   float64  `env:"UNIFORM_MAX"`
	LogNormSigma float64  `env:"LOGNORM_SIGMA"`
	LogNormMu    float64  `env:"LOGNORM_MU"`
	uv           uvType
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
func (f *RateType) Validate() error {
	switch f.Dist {
	case UniformE:
		if f.UniformMin == 0 || f.UniformMax == 0 {
			return fmt.Errorf("UNIFORM_MIN and UNIFORM_MAX should be positive values. current: %f-%f", f.UniformMin, f.UniformMax)
		}
		f.uv = distuv.Uniform{Min: f.UniformMin, Max: f.UniformMax}
	case LogNormalE:
		if f.LogNormMu == 0 || f.LogNormSigma == 0 {
			return fmt.Errorf("LOGNORM_MU and LOGNORM_SIGMA should be positive values. current: %f-%f", f.LogNormMu, f.LogNormSigma)
		}
		f.uv = distuv.LogNormal{Mu: f.LogNormMu, Sigma: f.LogNormSigma}
	default:
		return fmt.Errorf("RATE_DISTRIBUTION should be one of %s or %s. current: %s", UniformE, LogNormalE, f.Dist)
	}
	return nil
}
func (f *RateType) GetRandom() SizeType {
	return SizeType(f.uv.Rand())
}
