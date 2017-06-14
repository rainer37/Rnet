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
const (
	PLAIN_OVERLAY string = "PLAIN_OVERLAY"
	APP_PREFIX byte = 'U'
)
type Plain_node struct {
	Node
	NList map[string]string // list of peers with their ip's
}

func (p *Plain_node) Init(ip string, port string, extra ...string) (uint8, error) {
		
	p.IP = ip
	p.Port_string = port
	p.ID = id_generate_4(p.IP, p.Port_string) 
	p.State = 1
	p.NList = make(map[string]string)
	p.NList[ip+":"+port] = p.ID

	fmt.Printf("%sMy IP is [%s:%s] using [%s] as base arch\n", DHT_PREFIX, p.IP, p.Port_string, PLAIN_OVERLAY)

	p.listen()

	return 0, nil
}


func (p *Plain_node) Join(ip string, port string) (uint8, error) {

	fmt.Printf("%sLocal IP: [%s:%s]\n",DHT_PREFIX, p.IP, p.Port_string)
	
	id := id_generate_4(p.IP, p.Port_string) // generate id for new node

	fmt.Printf("%sLocal ID: [%s]\n", DHT_PREFIX, id)

	p.ID = id

	if id == "0" {
		fmt.Println(DHT_PREFIX+"Error on generating id")
	}

	//conn, err := p.connect(ip+":"+ port)
	conn, err := net.Dial("tcp", ip+":"+ port)

	if err != nil {
		fmt.Println(DHT_PREFIX+"Cannot start connection to target address")
		return 1, err
	}
	// send J msg to the existing node for ack.
	p.Send(conn, "J "+p.ID+" "+p.IP+":"+p.Port_string)

	// wait for the response.
	err = p.handleRequest(conn)

	if err != nil {
		fmt.Println(DHT_PREFIX+"Did not receive ACK, cannot proceed to listen")
		return 1, err
	}

	//set up done, enter the listening states
	p.listen()

	return 0, nil
} 

func (p *Plain_node) Send(conn net.Conn, msg string) (uint8, error) {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		return 1, err
	}
	return 0,nil
}

/************************************/
/*          INDIVIDUAL FNs          */
/************************************/

// listen on local ip and port in the Node
func (p *Plain_node) listen() {
	l, err := net.Listen("tcp", p.IP+":"+p.Port_string)

	if err != nil {
		eprint(err)
		return
	}

	defer l.Close()

	fmt.Println(DHT_PREFIX+"Start Listening on ["+p.IP+":"+p.Port_string+"]")

	// add self information to the map
	p.NList[p.IP+":"+p.Port_string] = p.ID

	// server start and ready for receiving msgs.
	for {
		conn, err := l.Accept()

		fmt.Println(DHT_PREFIX+"Incoming connection from "+conn.RemoteAddr().String())

		if err != nil {
			eprint(err)
			conn.Close()
		}
		go p.handleRequest(conn)
	}
}

// error message printing routine
func eprint(err error)  {
	fmt.Println(DHT_PREFIX, err.Error())
}

// general server dispatcher
// handler multiplexer.
func (p *Plain_node) handleRequest(conn net.Conn) error {
	// Make a buffer to hold incoming data.
	defer conn.Close()

	buf := make([]byte, 1024)

	// Read the incoming connection into the buffer.
	_ , err := conn.Read(buf)

	if err != nil {
		eprint(err)
		return err
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
		case 'B':
			fmt.Println(DHT_PREFIX+"Newbie joined.")
			p.add_newbie(msg)
		case APP_PREFIX:
			fmt.Println(DHT_PREFIX+"Application Data Received.")
		default:
			fmt.Println(DHT_PREFIX+"Unknown msg format")
			conn.Write([]byte("Don't Know What You Mean by"+msg))
	}

	return nil
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
	p.NList[ip] = id
	mutex.Unlock()

	fmt.Printf("%sWelcome [%s]\n", DHT_PREFIX, id)
	p.print_peers()

	nlist := p.generate_nlist()

	// TODO: make it concurrent

	for i,_ := range p.NList {
		if i != ip && i != p.IP+":"+p.Port_string {
			fmt.Println(i,"doesn't know",ip)
			c,err := net.Dial("tcp", i)
			if c == nil {
				eprint(err)
				// forget it if not connectable
				continue
			}
			p.Send(c, "B "+id+" "+ip)
			c.Close()
		}
	}

	p.Send(conn, "A "+nlist)
	//conn.Write([]byte("A "+nlist))
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

	fmt.Println("Current map")
	p.print_peers()

	// parse nlist received and update local nlist map
	for _, idip := range peers {
		id_ip := strings.Split(idip, "@")
		p.NList[id_ip[0]] = id_ip[1]
	}

	fmt.Println("After map")
	p.print_peers()

	p.State = 1 // 1 for normally RUNNING

	fmt.Println(DHT_PREFIX+"Join Succeed.")
}

func (p *Plain_node) add_newbie(msg string) {
	fmt.Println(msg)
	str := strings.Split(msg, " ")

	if len(str) < 3 {
		fmt.Println("Too few inforamtion received, aborting updating NList")
		return 
	}

	id := str[1]
	ip := str[2]

	mutex.Lock()
	p.NList[ip] = id
	mutex.Unlock()

	p.print_peers()
}

func (p *Plain_node) print_peers() {
	c := 0
	for i,v := range p.NList {
		fmt.Printf("[%d]: [%s:%s]\n", c, i, v)
		c++
	}
}

