package main

import (
	"math/rand"
	"net"
	"time"
)

func randomNum(min, max int) int {
	return rand.Intn(max-min) + min
}

func getRandomTimeInterval(avg int) int {
	return rand.Intn(avg) * 2
}

func connectCollector(collector string) (*net.UDPConn, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", collector)
	if err != nil {
		log.Fatal(err)
	}
	
	return net.DialUDP("udp", nil, udpAddr)
}

func sendPackets(conn *net.UDPConn) error {
	rand.Seed(time.Now().Unix())
	n := randomNum(50, 1000)
	// add spike data
	if opts.SpikeProto != "" {
		GenerateSpike()
	}
	if n > 900 {
		data := GenerateNetflow(8)
		buffer := BuildNFlowPayload(data)
		_, err := conn.Write(buffer.Bytes())
		if err != nil {
			return err
		}
	} else {
		data := GenerateNetflow(16)
		buffer := BuildNFlowPayload(data)
		_, err := conn.Write(buffer.Bytes())
		if err != nil {
			return err
		}
	}
	
	return nil
	
	// add some periodic spike data
	//if n < 150 {
	//	sleepInt := time.Duration(3000)
	//	time.Sleep(sleepInt * time.Millisecond)
	//}
	//sleepInt := time.Duration(n)
	//time.Sleep(sleepInt * time.Millisecond)
}