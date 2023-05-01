package main

import (
	config "NeSyCa/Config"
	writer "NeSyCa/Writer"
	"math/rand"
	"net"
	"time"

	"github.com/sirupsen/logrus"
)

// var cfg config.SettingsType
// var startAt int64
// var totalUDPSize = 0
// var totalHTTPSize = 0
// var udpSizeChannel = make(chan int)
// var httpSizeChannel = make(chan int)

//	func stats() {
//		go func() {
//			ticker := time.NewTicker(time.Duration(cfg.StatsInterval) * time.Second)
//			for {
//				<-ticker.C
//				duration := time.Now().Unix() - startAt
//				log.WithField("duration", utils.SecondsToHuman(int(duration))).WithField("udp volume", utils.ByteCountIEC(int64(totalUDPSize))).WithField("udp rate", utils.ByteCountIEC(int64(totalUDPSize/int(duration)))).WithField("http volume", utils.ByteCountIEC(int64(totalHTTPSize))).WithField("http rate", utils.ByteCountIEC(int64(totalHTTPSize/int(duration)))).Warn()
//			}
//		}()
//		go func() {
//			for {
//				select {
//				case vu := <-udpSizeChannel:
//					totalUDPSize += vu
//				case vh := <-httpSizeChannel:
//					totalHTTPSize += vh
//				}
//			}
//		}()
//	}
//
//	func main() {
//		cfg = *config.Config()
//		stats()
//		startAt = time.Now().Unix()
//		udp_enabled := len(cfg.UDPTargets) > 0
//		http_enabled := len(cfg.HttpTargets) > 0
//		for {
//			val := int(rand.Float32()*(cfg.SizeMax-cfg.SizeMin) + cfg.SizeMin)
//			sleep := rand.ExpFloat64() * cfg.RateLambda
//			if udp_enabled {
//				target := cfg.UDPTargets[rand.Intn(len(cfg.UDPTargets))].AsString()
//				go writer.WriteUDP(&target, val, udpSizeChannel)
//				log.WithField("type", "udp").WithField("sleep", sleep).WithField("size", val).WithField("target", target).Info()
//			}
//			if http_enabled {
//				target := cfg.HttpTargets[rand.Intn(len(cfg.HttpTargets))].AsString()
//				go writer.WriteHTTP(&target, val, httpSizeChannel)
//				log.WithField("type", "http").WithField("sleep", sleep).WithField("size", val).WithField("target", target).Info()
//			}
//			time.Sleep(time.Duration(sleep) * time.Second)
//		}
//	}
type request struct {
	target    net.TCPAddr
	size      int
	rate      int
	batchSize int
}

var rnd = rand.New(rand.NewSource(time.Now().Unix()))

func main() {
	for {
		req := makeRequest()
		logrus.Debugf("request: %+v", req)
		writer.WriteTCP(net.TCPAddr(req.target), req.size, req.rate, req.batchSize)
	}
}
func makeRequest() request {
	cfg := config.Config()
	return request{
		target:    net.TCPAddr(cfg.TCPTargets[rnd.Intn(len(cfg.TCPTargets))]),
		size:      rnd.Intn(int(cfg.FileSizeMax)-int(cfg.FileSizeMin)) + int(cfg.FileSizeMin),
		batchSize: rnd.Intn(int(cfg.BatchSizeMax)-int(cfg.BatchSizeMin)) + int(cfg.BatchSizeMin),
		rate:      rnd.Intn(int(cfg.RateLimitMax)-int(cfg.RateLimitMin)) + int(cfg.RateLimitMin),
	}
}
