package apctl

import(
	"fmt"
	"strings"
	"net"
	"github.com/rainer37/Rnet/dht"
)

/*
	Forward the data from outside to the local application
	through UDS. 
*/
func Dispatch_app_data(msg string) {
	fmt.Println(msg)

	m := strings.Split(msg, " ")

	c := strings.LastIndex(m[1], ":")

	sock, amsg := m[1][c+1:], m[1][:c]
	
	fmt.Println(sock, amsg)

	Send_to_UDS(sock, amsg)
}


/*
	Forward app data to DHT through UDS.
*/
func Release_app_data(msg string, peer_addr string) {
	// send data over internet DHT nodes.
	fmt.Printf(AC_PREFIX+"forward [%s] to DHT\n", msg)
	dht.Send_to_ext_DHT(peer_addr, msg)
}

/*
	handle message from applications/DHT

	TODO: application data encoding unit.
*/
func DHT_UDS_handler(c net.Conn) {
	defer c.Close()

	buf := make([]byte, 1024)

	nr, err := c.Read(buf)
	if err != nil {
		return
	}

	data := string(buf[0:nr])

	fmt.Printf(AC_PREFIX+"received [%s] from app\n", data)

	if string(data)[0] == 'U' {
		Dispatch_app_data(data)
	} else if string(data)[0] == 'P' { 
		c.Write([]byte(Get_friends()))
	} else {
		d := strings.Split(data, " ")
		Release_app_data(d[0], d[1])
	}
}

