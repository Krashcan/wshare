package main

import (
	"fmt"
	"log"
	"os"
	"net/http"
	"net"
	"strings"
)

var fileName string
var tag string
func main() {
	tag = os.Args[1]
	//dont look for file name when user is looking up help
	if tag!= "-help"{
		fileName = os.Args[2]
	}

	//logic for various tags
	if tag=="-f"{//for sharing folders
		http.Handle("/",http.StripPrefix("/",http.FileServer(http.Dir(fileName))))
		fmt.Printf("Sharing on %s:8080\n",GetOutboundIP())
	}else if tag=="-s"{//for sharing files
		http.HandleFunc("/",ShareFile)
		fmt.Printf("Sharing on %s:8080\n",GetOutboundIP())	
	}else if tag=="-help"{
		fmt.Println("\n\twshare -f <absolute folder path> for folders\n\twshare -s <absolute file path> for files")
	}else{
		fmt.Println("\nCommand not found\n\twshare -help,for complete list of commands")
	}
	
	//only start the server when sharing is required
	if tag=="-f" || tag=="-s"{
		log.Fatal(http.ListenAndServe(":8080",nil))	
	}
}

func ShareFile(w http.ResponseWriter,r *http.Request){
	http.ServeFile(w,r,fileName)
}

func GetOutboundIP() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()
    localAddr := conn.LocalAddr().String()
    idx := strings.LastIndex(localAddr, ":")
    return localAddr[0:idx]
}