package dht

/*
	Unstructed Plain Node Overlays For Testing
	(Conceptually not a Distributed Hash Table design)
	The design does not provide ANY fault-tolerance at all...

	Routing Table format:
		[ID : IP]

	ID space: string (127.0.0.1:1338 => 1270000000011338)
	IP space: normal (assuming unique and persistent IPv4)

	Joining Mechanism:
		1. Send IP to any existing node in system with msg starting with 'J' ID and ip:port.
		2. Responding peer replies with 'A' as acknowlegement, and the NLIST from peer to
		   Update local NLIST.
		3. Joining routine done. 
*/

import (
	"fmt"
	"net"
	"strings"
	"sync"
)

var mutex = &sync.Mutex{}
const PLAIN_OVERLAY string = "PLAIN_OVERLAY"

type Plain_node struct {
	Node
	NList map[string]string // list of peers with their ip's
}

func (p *Plain_node) Init(ip string, port string, extra ...string) (uint8, error) {
		
	p.IP = ip
	p.ID = id_generate_4(p.IP, p.Port_string) 
	p.Port_string = port
	p.State = 0
	p.NList = make(map[string]string)
	p.NList["1234567891231111"] = "123.456.789.123:1111"
	p.NList["2552552552551234"] = "255.255.255.255:1234"

	fmt.Printf("%sMy IP is [%s:%s] using [%s] as base arch\n", DHT_PREFIX, p.IP, p.Port_string, PLAIN_OVERLAY)

	l, err := net.Listen("tcp", p.IP+":"+p.Port_string)

	if err != nil {
		eprint(err)
		return 1, err
	}

	defer l.Close()

	fmt.Println(DHT_PREFIX+"Start Listening on ["+p.IP+":"+p.Port_string+"]")

	// add self information to the map
	p.NList[p.IP] = p.ID

	// server start and ready for receiving msgs.
	for {
		conn, err := l.Accept()

		fmt.Println(DHT_PREFIX+"Incoming connection from "+conn.RemoteAddr().String())

		if err != nil {
			eprint(err)
			conn.Close()
			return 1, err
		}
		go p.handleRequest(conn)
	}

	return 0,nil
}


func (p *Plain_node) Join(ip string, port string) (uint8, error) {

	fmt.Println(p.IP, p.Port_string)
	id := id_generate_4(p.IP, p.Port_string) // generate id for new node

	fmt.Printf("%sMy ID: [%s]\n", DHT_PREFIX, id)

	p.ID = id

	if id == "0" {
		fmt.Println(DHT_PREFIX+"Error on generating id")
	}

	// setup address of local and remote
	remoteAddr,_ := net.ResolveTCPAddr("tcp", ip+":"+port)
	localAddr,_ := net.ResolveTCPAddr("tcp", p.IP+":"+p.Port_string)

	conn, err := net.DialTCP("tcp", localAddr, remoteAddr)

	if err != nil {
		fmt.Println(DHT_PREFIX+"Cannot start connection to target address")
		return 1, err
	}

	// send J msg to the existing node for ack.
	conn.Write([]byte("J "+p.ID+" "+p.IP+":"+p.Port_string))

	// wait for the response.
	p.handleRequest(conn)

	// starting listening

	// l, err := net.Listen("tcp", p.IP+":"+p.Port_string)

	// if err != nil {
	// 	eprint(err)
	// 	return 1, err
	// }

	// defer l.Close()

	// fmt.Println(DHT_PREFIX+"Start Listening on ["+p.IP+":"+p.Port_string+"]")

	// // add self information to the map
	// p.NList[p.IP] = p.ID

	// // server start and ready for receiving msgs.
	// for {
	// 	conn, err := l.Accept()

	// 	fmt.Println(DHT_PREFIX+"Incoming connection from "+conn.RemoteAddr().String())

	// 	if err != nil {
	// 		eprint(err)
	// 		conn.Close()
	// 		return 1, err
	// 	}
	// 	go p.handleRequest(conn)
	// }	

	return 0, nil
} 

/************************************/
/*          INDIVIDUAL FNs          */
/************************************/

// error message printing routine
func eprint(err error)  {
	fmt.Println(DHT_PREFIX, err.Error())
}

// general server dispatcher
// handler multiplexer.
func (p *Plain_node) handleRequest(conn net.Conn) {
	// Make a buffer to hold incoming data.
	defer conn.Close()

	buf := make([]byte, 1024)

	// Read the incoming connection into the buffer.
	_ , err := conn.Read(buf)

	if err != nil {
		eprint(err)
		return 
	}

	msg := strings.Trim(string(buf), "\x00")
	msg = strings.Trim(msg, "\n")

	switch msg[0] {
		case 'J':
			fmt.Println(DHT_PREFIX+"Joinging Request Received.")
			p.handle_join(msg, conn)
		case 'A':
			fmt.Println(DHT_PREFIX+"Join Ack Received.")
			p.handle_join_ack(msg)
		default:
			fmt.Println(DHT_PREFIX+"Unknown msg format")
			conn.Write([]byte("Don't Know What You Mean by"+msg))
	}
}

// generate the id by the ipv4 address and port number
// ex. 192.168.0.1:1338 => 1921680000011338
// return string id or 0 for failure
func id_generate_4(ip string, port string) string {
	ips := strings.Split(ip, ".")

	if len(ips) != 4 {
		fmt.Println(DHT_PREFIX+"Malformed IPv4")
		return "0"
	}

	id := ""

	// fill 0's if not an three digits number in each part of IPs
	for _,v := range ips {
		id = id + strings.Repeat("0", 3-len(v)) + v
	}

	return id+port
}

// convert map[id]ip to single string in format:
// id1@ip1&id2@ip2&id3@ip3
func (p *Plain_node) generate_nlist() string {
	nlist := ""
	for id, ip := range p.NList {
		nlist = nlist + "&"+ id +"@"+ip
	}

	return nlist[1:]
}

/*	peer joining handler
 	msg : ['J' ID IP:PORT]
 	reply msg: ['A' NLIST]
*/
func (p *Plain_node) handle_join(msg string, conn net.Conn) {
	str := strings.Split(msg, " ")

	// check mulformat joining request
	if len(str) < 3 {
		fmt.Println("Too few inforamtion received, aborting Join")
		return 
	}

	id := str[1]
	ip := str[2]

	mutex.Lock()
	p.NList[id] = ip
	mutex.Unlock()

	fmt.Printf("%sWelcome [%s]\n", DHT_PREFIX, id)
	fmt.Println(p.NList)
	
	nlist := p.generate_nlist()

	conn.Write([]byte("A "+nlist))
}

/*
	Receive ack + nlist for update current nlist map
	msg : ['A' NLIST]
	Then update local state to 1(RUNNING)
*/
func (p *Plain_node) handle_join_ack(msg string) {

	str := strings.Split(msg, " ")

	if len(str) < 2 {
		fmt.Println("Too few inforamtion received, aborting Join")
		return 
	}

	peers := strings.Split(str[1], "&")

	fmt.Println("Current map", p.NList)

	// parse nlist received and update local nlist map
	for _, idip := range peers {
		id_ip := strings.Split(idip, "@")
		p.NList[id_ip[0]] = id_ip[1]
	}

	fmt.Println("After map", p.NList)

	p.State = 1 // 1 for normally RUNNING

	fmt.Println(DHT_PREFIX+"Join Succeed.")
}



