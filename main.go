package main

import (
	"fmt"
	"log"
	"net/http"
	"net"
	"strings"
	"flag"
)

func main() {
	//flag to specify whether we will be uploading folder or a single file
	folder := flag.Bool("f",false,"Use for serving folders on the server")
	
	flag.Parse()
	

	if len(flag.Args())>0{
		if *folder{
			http.Handle("/",http.StripPrefix("/",http.FileServer(http.Dir(flag.Args()[0]))))
			fmt.Printf("Sharing folder on %s:8080\n",GetOutboundIP())
		}else{
			http.HandleFunc("/",ShareFile)
			fmt.Printf("Sharing file on %s:8080\n",GetOutboundIP())
		}
		
		log.Fatal(http.ListenAndServe(":8080",nil))
	}else{
		fmt.Println("Invalid usage. No file mentioned. Use wshare -h for help.")
	}
	
}


//function to share files
func ShareFile(w http.ResponseWriter,r *http.Request){
	http.ServeFile(w,r,flag.Args()[0])
}
//function to get the ip address for other devices to communicate through
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