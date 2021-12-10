package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	// TODO - create an init/setup func, put this in
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// open up the listening address for returning ICMP packets. Or is this two way somehow?
	// c should be renamed socket?
	// c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	icmpSocket, err := net.ListenPacket("ip4:icmp", "0.0.0.0")

	if err != nil {
		log.Fatal(err)
	}
	defer icmpSocket.Close()

	c := ipv4.NewPacketConn(icmpSocket)
	defer c.Close()

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}

	wb, err := wm.Marshal(nil)
	if err != nil {
		log.Fatal(err)

	}
	// if _, err := c.WriteTo(wb, nil, &dst); err != nil {
	// 	fmt.Fprintf(os.Stderr, "Could not send the ICMP packet: %s\n", err)
	// 	os.Exit(1)
	// }
	if _, err := c.WriteTo(wb, nil, &net.UDPAddr{IP: net.ParseIP("142.1.217.155"), Zone: "en0"}); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Got here")

	readBuffer := make([]byte, 1500)
	// n, peer, err := c.ReadFrom(readBuffer)
	readBytes, _, hopNode, err := c.ReadFrom(readBuffer)
	if err != nil {
		log.Fatal(err)
	}

	readMessage, err := icmp.ParseMessage(58, readBuffer[:readBytes])
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Received readMessage - %+v - from hopNode - %+v", readMessage, hopNode)

	// client := &http.Client{
	// 	Transport: &http.Transport{
	// 		DialContext: dialer.DialContext,
	// 	},
	// }

	// resp, err := client.Get("http://ixmaps.ca")
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()

	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(string(body))
}
