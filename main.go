package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"
)

type Duration struct {
	Bootup struct {
		Kernel          float64 `json:"kernel"`
		Initrd          float64 `json:"initrd"`
		Userspace       float64 `json:"userspace"`
		GraphicalTarget float64 `json:"graphical.target"`
	} `json:"bootup"`
	TimeUnit string `json:"time-unit"`
}

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

	//retBuf := bytes.Trim(buf, "\x00")
	retBuf := parseSystemdAnalyze(buf)

	return retBuf, err
}

func parseSystemdAnalyze(cmdStdout []byte) []byte {

	// Parses the systemd-analyze output and returns the duration in json format

	var duration Duration
	var jsonDuration []byte

	str1 := strings.TrimPrefix(string(cmdStdout[:]), "Startup finished in ")
	str2 := strings.Split(str1, " = ")
	str3 := strings.Split(str2[0], " + ")

	//strTotal = str2[1]

	for _, str := range str3 {

		str4 := strings.Split(str, " (")

		//fmt.Printf("%s\n", str4[1])
		timeDuration, _ := time.ParseDuration(strings.ReplaceAll(str4[0], "in ", ""))
		//fmt.Printf("%f\n", timeDuration.Seconds())

		if strings.Contains(str4[1], "kernel") {
			duration.Bootup.Kernel = timeDuration.Seconds()
		} else if strings.Contains(str4[1], "initrd") {
			duration.Bootup.Initrd = timeDuration.Seconds()
		} else if strings.Contains(str4[1], "userspace") {
			duration.Bootup.Userspace = timeDuration.Seconds()
		} else if strings.Contains(str4[1], "graphical.target") {
			duration.Bootup.GraphicalTarget = timeDuration.Seconds()
		}
	}

	//fmt.Printf("%q -- %d\n", str3, len(str3))

	duration.TimeUnit = "seconds" //default time unit

	jsonDuration, _ = json.Marshal(duration)

	return jsonDuration
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
