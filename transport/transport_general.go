package transport

import (
	"fmt"
)

const TRANS_PREFIX string = "[TRANS] "

func Trans_Boot() {
	fmt.Println(TRANS_PREFIX+"Transport Communication Unit Booting...")
} 

