package writer

import (
	"net"

	"github.com/sirupsen/logrus"
)

func Write(conn *net.Conn, d *[]byte) (int, error) {
	n, err := (*conn).Write(*d)
	logger := logrus.WithField("conn", *conn).WithField("size", n)
	if err != nil {
		logger.WithError(err).Debug("error writing data")
	}
	return n, err
}
func NewReport(target string) *ReportType {
	rpt := ReportType{Target: target}
	rpt.Reset()
	return &rpt
}
