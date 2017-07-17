package apctl

/*
	Local application controller unit.
	main func:
		1. display meta of application.
		2. terminate/execute applications.
		3. application state maintain.
		4. APC to DHT communication.

*/

import(
	"fmt"
	"os/exec"
	"os"
	"net"
	"io/ioutil"
	"github.com/rainer37/Rnet/dht"
)

const AC_PREFIX string = "[APC]\t"
var appdir string = os.Getenv("PWD")+"/app/"

var apps map[string]App = make(map[string]App) // name to App mapping.
var friends map[string]string = make(map[string]string) // fname:ip.
var SOCKET_ADDR string = os.Getenv("PWD")+"/go.sock"
var global_ln net.Listener 

type App struct {
	name string
	sock string
	desp string // description of application
}

// ini func to refresh the local application states and metadata.
// TODO:
func AC_boot() {
	fmt.Println(AC_PREFIX+"Application Controller Booting...")
	/*
		build/compile the new applications.
	*/
	apps = make(map[string]App)
	//build_app("chat/chat")

	get_local_apps()

	// listen for global UDS for comm from all applications
	Listen_UDS(SOCKET_ADDR, DHT_UDS_handler)
}

// get all locally installed application from app_list.json
// TODO: format of json
func Get_app_list() {

	for i,v := range apps {
		fmt.Printf("[ %s ] : %s\n", i, v.desp)
	}

}

// return the list of peers as a string.
func Get_friends() string {
	peers := dht.Peer_list()
	var ps string = ""

	for _,v := range peers {
		ps = ps + " " + v
	}

	return ps
}

/*
	check all applications in the appdir, and build them all.
*/
func get_local_apps() {
	app_files, _ := ioutil.ReadDir(appdir)
	size := 0
	for _, f := range app_files {
		if f.Name()[0] != '.' {
			fmt.Println(f.Name())
			apps[f.Name()] = App{
				f.Name(),
			 	appdir+"/"+f.Name()+"/"+f.Name()+".sock",
			 	f.Name()+" is amazing",
			}
			build_app(f.Name())
			size++
		}
	}

	fmt.Printf(AC_PREFIX+"%d appplications found.\n", size)
}

/*
	build go application from source code with same name.
*/
func build_app(source string) {
	_, err := exec.Command("go","build", "-o", appdir+source+"/"+source, appdir+source+"/"+source+".go").CombinedOutput()
	if err != nil {
		os.Stderr.WriteString("Building "+err.Error())
	}
	//fmt.Println(string(output))
}

/*
	execute an executable, it does not have to be a go app.
*/
func Exec_app(source string, OS string) error {

	var cmd string
	var args []string

	if OS == "darwin" {
		cmd = "osascript"
		args = []string{"-e", "tell app \"Terminal\" to do script \"cd "+os.Getenv("PWD")+";"+appdir+source+"/"+source+"\""}
	} else if OS == "linux" {
		cmd = "gnome-terminal"
		args = []string{"-e", appdir+source+"/"+source}
	} else if OS == "windows" {
		fmt.Println("Windows cmd coming up....")
		return nil
	}

	_, err := exec.Command(cmd, args...).CombinedOutput()
	if err != nil {
		os.Stderr.WriteString(err.Error())
		return err
	}
	return nil
	//fmt.Println(string(output))
}	

func get_appsock_by_name(appname string) string {
	return apps[appname].sock
}
