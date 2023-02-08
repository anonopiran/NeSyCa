package main

import (
	config "NeSyCa/Config"
	utils "NeSyCa/Utils"
	writer "NeSyCa/Writer"
	"math/rand"
	"time"

	"net/http"
	_ "net/http/pprof"

	log "github.com/sirupsen/logrus"
)

var cfg = config.Config()
var startAt int64
var totalSize = 0
var sizeChannel = make(chan int)

func stats() {
	go func() {
		ticker := time.NewTicker(time.Duration(cfg.StatsInterval) * time.Second)
		for {
			<-ticker.C
			duration := time.Now().Unix() - startAt
			log.WithField("duration", utils.SecondsToHuman(int(duration))).WithField("volume", utils.ByteCountIEC(int64(totalSize))).WithField("rate", utils.ByteCountIEC(int64(totalSize/int(duration)))).Warn()
		}
	}()
	go func() {
		for {
			v := <-sizeChannel
			totalSize += v
		}
	}()
}
func main() {
	go func() {
		http.ListenAndServe(":1234", nil)
	}()
	stats()
	startAt = time.Now().Unix()
	for {
		val := int(rand.Float32()*(cfg.SizeMax-cfg.SizeMin) + cfg.SizeMin)
		sleep := rand.ExpFloat64() * cfg.RateLambda
		target := cfg.UDPTargets[rand.Intn(len(cfg.UDPTargets))].AsString()
		log.WithField("sleep", sleep).WithField("size", val).WithField("target", target).Info()
		go writer.Write(&target, val, sizeChannel)
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}
