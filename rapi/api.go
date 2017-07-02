package rapi

/*
	api for applications to register, communicate, and more
*/

import(
	"net"
	"os"
	"fmt"
	"strings"
	"github.com/rainer37/Rnet/apctl"
)

/*
	register application with matadata.
*/
func Register() {

}

/*
	starting the uds server for the application.
*/
func Serve(handler func(net.Conn)) {
	apctl.Listen_UDS(os.Args[0]+".sock", handler)
}

/*
	send to central uds server for further forwarding.
*/
func Send(msg string) {
	// get direct name of app
	sock_addr := os.Args[0][strings.LastIndex(os.Args[0], "/")+1:]
	//sock_addr := os.Args[0]+".sock"
	apctl.Send_to_UDS("", msg+":"+sock_addr)
}

func Send_rip(remote_ip string, msg string) {
	sock_addr := os.Args[0][strings.LastIndex(os.Args[0], "/")+1:]
	//sock_addr := os.Args[0]+".sock"
	apctl.Send_to_UDS("", msg+":"+sock_addr+" "+remote_ip)
}

/*
	retrieve all the peers met.
*/
func Peers() {
	fmt.Println("Friends:")
	
	peers := strings.Split(apctl.Send_and_receive("", "P"), " ")

	for i,v := range peers {
		fmt.Printf("%d : %s\n", i, v)
	}
}