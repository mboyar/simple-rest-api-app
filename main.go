package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	
)

const strVersion string = "v0.1"
const bufSize uint16 = 500

func main() {

	strMsg := "Server ready, endpoints: /version and /duration \n"

	fmt.Println(strMsg)

	strVersion := getVersion()
	bufDuration, _ := getDuration()

	go responser(bufDuration, strVersion, strMsg)

	select {}
}

func getVersion() string {
	return strVersion
}

func getDuration() ([]byte, error) {

	//preinstalled command name in systemd based systems.
	//see http://manpages.ubuntu.com/manpages/bionic/man1/systemd-analyze.1.html for the details
	strCmd := "systemd-analyze"

	_, err := exec.LookPath(strCmd)
	if err != nil {
		log.Println(strCmd + " command cannot found in your system")
		return []byte("systemd-analyze fake output"), err
	}

	//calling with either "time" arg or w/o any arg results same bootup duration output
	cmd := exec.Command(strCmd, "time")
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		log.Fatal(err)
	}

	//Exec command as non-block
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	//reading from io.ReadCloser stream to dynamic allocated byte array
	buf := make([]byte, bufSize)
	if _, err := io.ReadFull(stdout, buf); err != nil {
		log.Println(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

	retBuf := bytes.Trim(buf, "\x00")

	return retBuf, err
}

func responser(duration []byte, version string, msg string) {

	strDuration := string(duration)

	//h0,h1,h2 are callback functions triggered when any request got by http server

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

	//Listen port 8080 in blocking mode
	log.Fatal(http.ListenAndServe(":8080", nil))
}
