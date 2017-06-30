package apctl

/*
	Local application controller unit.
	main func:
		1. display meta of application.
		2. terminate/execute applications.
		3. application state maintain.
		4. APP to DHT communication.

*/

import(
	"fmt"
	"os/exec"
	"os"
	"net"
	"strings"
	"github.com/rainer37/Rnet/dht"
)

const AC_PREFIX string = "[APC]\t"
const appdir string = "./app/"

var apps map[string]App = make(map[string]App)
var SOCKET_ADDR string = os.Getenv("PWD")+"/go.sock"
var global_ln net.Listener 

type App struct {
	name string
	sock string
}

// ini func to refresh the local application states and metadata.
// TODO:
func AC_boot() {
	fmt.Println(AC_PREFIX+"Application Controller Booting...")
	/*
		build/compile the new applications.
	*/
	apps = make(map[string]App)
	build_app("chatchat/chat")

	// listen for global UDS for comm from all applications
	Listen_UDS(SOCKET_ADDR, DHT_UDS_handler)
}

func build_app(source string) {
	_, err := exec.Command("go","build", "-o", appdir+source, appdir+source+".go").CombinedOutput()
	if err != nil {
		os.Stderr.WriteString("1"+err.Error())
	}
	//fmt.Println(string(output))
}

func Exec_app(source string) {
	output, err := exec.Command("gnome-terminal","-e", appdir+source).CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
	}
	fmt.Println(string(output))
}	

// get all locally installed application from app_list.json
// TODO: format of json
func Get_app_list() {

	for i,v := range apps {
		fmt.Printf("%s : %s", i, v.name)
		// TODO: read from apps map.
	}

	fmt.Println("Sample APP 1: Chat with your friend")
	fmt.Println("Sample APP 2: Restaurant Nearby")
	fmt.Println("Sample APP 3: Where am i")
}


/*
	Forward the data from outside to the local application
	through UDS. 
*/
func Dispatch_app_data(msg string) {
	fmt.Println(msg)
	m := strings.Split(msg, " ")
	amsg := m[1][:strings.Index(m[1],":")-1]
	sock := m[1][strings.Index(m[1],":"):]
	fmt.Println(sock, amsg)
}


/*
	Forward app data to DHT through UDS.
*/
func Release_app_data(msg string) {
		// send data over internet DHT nodes.
	dht.Send_to_ext_DHT("192.168.31.205:1338",msg)
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

	fmt.Println(AC_PREFIX+"Server got:", data)
	//return

	if string(data)[0] == 'U' {
		Dispatch_app_data(data)
	} else {
		Release_app_data(data)
	}
	
}

