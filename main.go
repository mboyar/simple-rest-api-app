package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	//"strings"
	//"unicode/utf8"
)

const strVersion string = "v0.1"

func main() {

	strMsg := "Server ready, endpoints: /version and /duration \n"

	fmt.Println(strMsg)

	strVersion := getVersion()
	bufDuration, _ := getDuration()

	responser(bufDuration, strVersion, strMsg)
}

func getVersion() string {
	return strVersion
}

func getDuration() ([]byte, error) {

	strCmd := "systemd-analyze"

	_, err := exec.LookPath(strCmd)
	if err != nil {
		log.Fatal(strCmd + " command cannot found in your system")
	}

	cmd := exec.Command(strCmd, "time")
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	buf := make([]byte, 500)
	if _, err := io.ReadFull(stdout, buf); err != nil {
		log.Println(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	//fmt.Println("len(buf):", len(buf))
	//fmt.Println("cap(buf):", cap(buf))
	retBuf := bytes.Trim(buf, "\x00")

	return retBuf, err
}

func responser(duration []byte, version string, msg string) {

	strDuration := string(duration)

	h0 := func(w http.ResponseWriter, _ *http.Request) {

		io.WriteString(w, msg)
	}

	h1 := func(w http.ResponseWriter, _ *http.Request) {

		io.WriteString(w, "Version: "+version+"\n")
	}

	h2 := func(w http.ResponseWriter, _ *http.Request) {

		io.WriteString(w, "Startup duration of the system: "+strDuration)
	}

	http.HandleFunc("/", h0)
	http.HandleFunc("/version", h1)
	http.HandleFunc("/duration", h2)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
