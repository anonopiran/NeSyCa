package main

import (
	config "NeSyCa/Config"
	utils "NeSyCa/Utils"
	writer "NeSyCa/Writer"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
)

var cfg config.SettingsType
var startAt int64
var totalUDPSize = 0
var totalHTTPSize = 0
var udpSizeChannel = make(chan int)
var httpSizeChannel = make(chan int)

func stats() {
	go func() {
		ticker := time.NewTicker(time.Duration(cfg.StatsInterval) * time.Second)
		for {
			<-ticker.C
			duration := time.Now().Unix() - startAt
			log.WithField("duration", utils.SecondsToHuman(int(duration))).WithField("udp volume", utils.ByteCountIEC(int64(totalUDPSize))).WithField("udp rate", utils.ByteCountIEC(int64(totalUDPSize/int(duration)))).WithField("http volume", utils.ByteCountIEC(int64(totalHTTPSize))).WithField("http rate", utils.ByteCountIEC(int64(totalHTTPSize/int(duration)))).Warn()
		}
	}()
	go func() {
		for {
			select {
			case vu := <-udpSizeChannel:
				totalUDPSize += vu
			case vh := <-httpSizeChannel:
				totalHTTPSize += vh
			}
		}
	}()
}
func main() {
	cfg = *config.Config()
	stats()
	startAt = time.Now().Unix()
	udp_enabled := len(cfg.UDPTargets) > 0
	http_enabled := len(cfg.HttpTargets) > 0
	for {
		val := int(rand.Float32()*(cfg.SizeMax-cfg.SizeMin) + cfg.SizeMin)
		sleep := rand.ExpFloat64() * cfg.RateLambda
		if udp_enabled {
			target := cfg.UDPTargets[rand.Intn(len(cfg.UDPTargets))].AsString()
			go writer.WriteUDP(&target, val, udpSizeChannel)
			log.WithField("type", "udp").WithField("sleep", sleep).WithField("size", val).WithField("target", target).Info()
		}
		if http_enabled {
			target := cfg.HttpTargets[rand.Intn(len(cfg.HttpTargets))].AsString()
			go writer.WriteHTTP(&target, val, httpSizeChannel)
			log.WithField("type", "http").WithField("sleep", sleep).WithField("size", val).WithField("target", target).Info()
		}
		time.Sleep(time.Duration(sleep) * time.Second)
	}
}
