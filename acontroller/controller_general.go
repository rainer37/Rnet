package acontroller

/*
	Local application controller unit.
	main func:
		1. display meta of application.
		2. terminate/execute applications.
		3. application state maintain.

*/

import(
	"fmt"
)

var apps map[string]App = make(map[string]App)

// application struct storing basic info.
type App struct {
	name, path, version, author, state, description string
}

type Rapp interface {
	Boot()
	Send(msg string,address_or_id string)
	Self_report()
}

// ini func to refresh the local application states and metadata.
// TODO:
func AC_boot() {
	fmt.Println("Application Controller Booting...")
}

// get all locally installed application from app_list.json
// TODO: format of json
func Get_app_list() {

	for i,v := range apps {
		fmt.Printf("%s : %s", i, v.description)
		// TODO: read from apps map.
	}

	fmt.Println("Sample APP 1: Chat with your friend")
	fmt.Println("Sample APP 2: Restaurant Nearby")
	fmt.Println("Sample APP 3: Where am i")
}
