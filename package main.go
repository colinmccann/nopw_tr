package main

import (
	"log"
	"net"
	"os"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {
	// open up the listening addres for returning ICMP packets. Or is this two way somehow?
	c, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}
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
	if _, err := c.WriteTo(wb, &net.UDPAddr{IP: net.ParseIP("142.1.217.155"), Zone: "en0"}); err != nil {
		log.Fatal(err)
	}

	rb := make([]byte, 1500)
	n, peer, err := c.ReadFrom(rb)
	if err != nil {
		log.Fatal(err)
	}

	rm, err := icmp.ParseMessage(58, rb[:n])
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Got %+v from %+v", rm, peer)

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

/*

package main

import (
	"fmt"
	"net"
)

func main() {
	output, err := destAddr("ixmaps.ca")
	if err != nil {
		fmt.Println("ERROR")
	}
	fmt.Println(output)
}

func destAddr(dest string) ([4]byte, error) {
	destAddr := [4]byte{0, 0, 0, 0}
	addrs, err := net.LookupHost(dest)
	if err != nil {
		return destAddr, err
	}
	addr := addrs[0]

	ipAddr, err := net.ResolveIPAddr("ip", addr)
	if err != nil {
		return destAddr, err
	}
	fmt.Println(ipAddr.IP.To4())
	copy(destAddr[:], ipAddr.IP.To4())
	return destAddr, nil
}

*/
