package controller

import (
	config "NeSyCa/Config"
	writer "NeSyCa/Writer"
	"crypto/rand"
	"time"

	"github.com/sirupsen/logrus"
)

func Run() {
	cfg := config.Config()
	tickWrite := time.NewTicker(time.Duration(cfg.DrainDelay))
	tickUpdate := time.NewTicker(time.Duration(cfg.UpdateTick))
	tickReport := time.NewTicker(time.Duration(cfg.ReportTick))
	chWrite, chRep := generateWriters()
	var d []byte
	rate := updateRandData(&d)
	repSet := *new([]writer.ReportType)
	for {
		select {
		case <-tickWrite.C:
			*chWrite <- d
		case <-tickUpdate.C:
			rate = updateRandData(&d)
		case <-tickReport.C:
			total := config.SizeType(0)
			perSite := make(map[string]config.SizeType)
			for _, r_ := range repSet {
				total += r_.Wrote
				perSite[r_.Target] += r_.Wrote
			}
			dur := config.SizeType(time.Duration(cfg.ReportTick).Seconds())
			realRate := total / dur
			logrus.WithField("set_rate", rate.Hr()).WithField("rate", realRate.Hr()).Warn()
			log := logrus.NewEntry(logrus.New())
			for t_, v_ := range perSite {
				realRate = v_ / dur
				log = log.WithField(t_, realRate.Hr())
			}
			log.Info()
			repSet = *new([]writer.ReportType)
		case rpt := <-*chRep:
			rpt.WithFields().Debug()
			repSet = append(repSet, rpt)
		}
	}
}
func generateWriters() (*chan []byte, *chan writer.ReportType) {
	cfg := config.Config()
	chWrite := make(chan []byte, 100)
	chRep := make(chan writer.ReportType)
	for _, targ := range cfg.TCPTargets {
		for c := 0; c < cfg.ConnPerTarget; c += 1 {
			writerObj := writer.TCPWriterType{
				Target:     targ,
				ReportChan: &chRep,
				WriteChan:  &chWrite,
			}
			go writerObj.Start()
		}
	}
	return &chWrite, &chRep
}
func calcBatchSize(bw config.SizeType, drain time.Duration) config.SizeType {
	return config.SizeType(float64(bw) * drain.Seconds())
}

func getRandom(d *[]byte, s config.SizeType) {
	new_d := make([]byte, s)
	rand.Read(new_d)
	*d = new_d
}
func updateRandData(d *[]byte) config.SizeType {
	cfg := config.Config()
	rt := cfg.Rate.GetRandom()
	bs := calcBatchSize(rt, time.Duration(cfg.DrainDelay))
	getRandom(d, bs)
	logrus.WithField("Rate", rt.Hr()).WithField("Batchsize", bs.Hr()).Info("rate updated")
	return rt
}
