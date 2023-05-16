package writer

import (
	config "NeSyCa/Config"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

type TCPWriterType struct {
	Target     config.TCPAddressType
	WriteChan  *chan []byte
	ReportChan *chan ReportType
}

type ReportType struct {
	Target   string
	Wrote    config.SizeType
	timeFrom int
}

func (f *TCPWriterType) Connect(conn *net.Conn) {
	logger := logrus.WithField("target", f.Target.Plain)
	var err error
	try := 1
	for {
		*conn, err = net.Dial("tcp", f.Target.Plain)
		if err == nil {
			logger.WithField("try", try).Debug("connect established")
			break
		}
		logger.WithField("try", try).WithError(err).Debug("connect error ...")
		try += 1
		time.Sleep(1 * time.Second)
		continue
	}
}
func (f *TCPWriterType) Start() {
	logger := logrus.WithField("target", f.Target.Plain)
	tickReport := time.NewTicker(1 * time.Second)
	rep := NewReport(f.Target.Plain)
	var conn net.Conn
	f.Connect(&conn)
	for {
		select {
		case data := <-*f.WriteChan:
			n, err := Write(&conn, &data)
			rep.Add(n)
			if err != nil {
				logger.Info("initiating new connection")
				conn.Close()
				conn = nil
				f.Connect(&conn)
			}
		case <-tickReport.C:
			*f.ReportChan <- *rep
			rep.Reset()
		}
	}
}
func (f *ReportType) Duration() int {
	return time.Now().Second() - f.timeFrom
}
func (f *ReportType) Reset() {
	f.timeFrom = time.Now().Second()
	f.Wrote = 0
}
func (f *ReportType) Add(v int) {
	f.Wrote += config.SizeType(v)
}
func (f *ReportType) BandWidth() config.SizeType {
	if f.Duration() == 0 {
		return config.SizeType(0)
	}
	return config.SizeType(int(f.Wrote) / f.Duration())
}
func (f *ReportType) WithFields() *logrus.Entry {
	bw := f.BandWidth()
	return logrus.
		WithField("target", f.Target).
		WithField("size", f.Wrote.Hr()).
		WithField("bandwidth", bw.Hr()).
		WithField("duration", f.Duration())
}
