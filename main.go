package main

import (
	"fmt"
	"log"
	"net/http"
	"net"
	"strings"
	"flag"
	"archive/zip"
	"os"
	"os/signal"
	"syscall"
	"io"
	"path/filepath"
)

func main() {
	

	//flag to specify whether we will be uploading folder or a single file
	folder := flag.Bool("f",false,"Use for serving folders on the server")
	zipped := flag.Bool("z",false,"Use for zipping folders and serving them as a single file on the server.(Deletes the zipped file once the server closes.)")
	save := flag.Bool("s",false,"Use with -z for saving the zipped files locally even after the server closes.")
	flag.Parse()
	

	if len(flag.Args())>0{
		if *folder{
			http.Handle("/",http.StripPrefix("/",http.FileServer(http.Dir(flag.Args()[0]))))
			fmt.Printf("Sharing folder on %s:8080\n",GetOutboundIP())
		}else{
			if *zipped{
				fmt.Println("zipping...")
				flag.Args()[0]=ZipFile()
				sigs := make(chan os.Signal, 3)
				signal.Notify(sigs, os.Interrupt,syscall.SIGTERM,syscall.SIGINT)

				go func() {
					<-sigs
            	    if !(*save){
						os.Remove(flag.Args()[0])
						fmt.Println("Closing server. File Deleted.")
						os.Exit(0)
					}else{
						fmt.Println("File saved in ",flag.Args()[0])
						os.Exit(0)
					}
    			}()
			}
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

func ZipFile() string {
	curDir,_ := os.Getwd()
	curDir = curDir+"/tmp.zip"
	source := flag.Args()[0]

	zipFile,err := os.Create(curDir)
	HandleError("os.Create: ",err)

	defer zipFile.Close()

	archive := zip.NewWriter(zipFile)
	defer archive.Close()

	info,err := os.Stat(source)
	HandleError("os.Stat: ",err)

	var baseDir string
	if info.IsDir(){
		baseDir= filepath.Base(source)
	}
	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {

		header, err := zip.FileInfoHeader(info)
		HandleError("zip.FileInfoHeader: ",err)
		if baseDir != "" {
			header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		HandleError("archive.CreateHeader: ",err)

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		HandleError("os.Open: ",err)
		defer file.Close()
		_, err = io.Copy(writer, file)

		return err
	})
	HandleError("filepath.Walk: ",err)
	return curDir
}

//function to get the ip address for other devices to communicate through
func GetOutboundIP() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    HandleError("net.Dial: ",err)
    defer conn.Close()
    localAddr := conn.LocalAddr().String()
    idx := strings.LastIndex(localAddr, ":")
    return localAddr[0:idx]
}

func HandleError(funcName string,err error){
	if err!=nil{
		log.Fatal(funcName,err)
	}
}