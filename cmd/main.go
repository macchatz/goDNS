package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/paulbellamy/ratecounter"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var k chan string

func getRoot(w http.ResponseWriter, r *http.Request) {
	hasFirst := r.URL.Query().Has("first")
	if hasFirst {
		first := r.URL.Query().Get("first")
		k <- first
	}
	io.WriteString(w, "This is my website!\n")

}

func redirect_Client(net_type string, buf []byte) {

	conn, err := net.Dial(net_type, "8.8.8.8")
	if err != nil {
		fmt.Println("This connection is not ok!")
	}
	fmt.Println("Google responded with this shit:", conn.LocalAddr().String())

}

func getTCPreq(con net.Conn) {
	io.WriteString(con, "This is the real netcat!\n")
	fmt.Println("Received TCP request on Local Address:", con.LocalAddr().String())
}

func getUDPreq(udpServer net.PacketConn, addr net.Addr, buf []byte) {

	fmt.Println("Received UDP request on Local Address:", addr.String())
	time := time.Now().Format(time.ANSIC)
	responseStr := fmt.Sprintf("time received: %v. Your message: %v!", time, string(buf))

	udpServer.WriteTo([]byte(responseStr), addr)
}

func getReq(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)

	fmt.Println("URL:", url)
	io.WriteString(w, "Thanks dig!\n")
}

//interface: Do you want to buy a DNS service

func main() {
	k = make(chan string, 50)

	fmt.Println("Welcome to the shitiest docker DNS server!")

	go func() {
		tcp_connection, TCP_err := net.Listen("tcp", ":8080")
		if TCP_err != nil {
			fmt.Println("TCP_TSAPOU")
		}

		go func() {
			for {
				conn, err := tcp_connection.Accept()
				if err != nil {
					// handle error
				}
				go getTCPreq(conn)
			}
		}()
	}()

	fmt.Println("Stopeakimas")

	go func() {

		udp_connection, UDP_err := net.ListenUDP("udp", &net.UDPAddr{
			Port: 1058,
			IP:   net.ParseIP("0.0.0.0"),
		})
		if UDP_err != nil {
			fmt.Println("UDP_TSAPOU")
		}

		for {
			buf := make([]byte, 1024)
			_, addr, err := udp_connection.ReadFromUDP(buf)
			if err != nil {
				log.Fatal(err)
				continue
			}
			go getUDPreq(udp_connection, addr, buf)
		}

	}()
	fmt.Println("Stopeakimas42")
	//upd_listener := net.ListenUDP("udp4",nil)

	//handle 3333 listener
	mux := http.NewServeMux()
	mux.HandleFunc("/", getRoot)
	mux.Handle("/metrics", promhttp.Handler())

	//handle 53
	zux := http.NewServeMux()
	zux.HandleFunc("/", getReq)
	//for prometheus metrics
	go func() {
		http.ListenAndServe(":6060", nil)
	}()

	go http.ListenAndServe(":3333", mux)
	fmt.Println("Server listening on port :3333")

	counter := ratecounter.NewRateCounter(1 * time.Second)
	for i := 0; i < 4; i++ {
		go func() {
			for {
				select {
				case <-k:
					counter.Incr(1)
				}
			}
		}()
	}

	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-ticker.C:
			fmt.Println("RPS", counter.Rate())
		}
	}

}
