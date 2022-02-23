package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/go-routeros/routeros"
)

type hostConfigList struct {
	Hosts []hostConfig `json:"hosts"`
}

type hostConfig struct {
	Name                  string `json:"name"`
	IP                    string `json:"ip"`
	Port                  int    `json:"port"`
	Username              string `json:"username"`
	Password              string `json:"password"`
	TLSInsecureSkipVerify bool   `json:"tlsInsecureSkipVerify,omitempty"`
}

func loadMappingFile(filename string) (*hostConfigList, error) {
	mappingFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer mappingFile.Close()

	bytes, err := ioutil.ReadAll(mappingFile)
	if err != nil {
		return nil, err
	}

	mapping := hostConfigList{}
	err = json.Unmarshal(bytes, &mapping)
	if err != nil {
		return nil, err
	}

	return &mapping, nil
}

func main() {
	host := flag.String("h", "", "host to set time for")
	mappingFile := flag.String("m", "mapping.json", "host-to-config-file mapping")
	ip := flag.String("i", "", "override IP address used to connect to host")
	flag.Parse()

	if *host == "" {
		log.Fatal("missing host parameter")
	}

	if *mappingFile == "" {
		log.Fatal("missing mapping file parameter")
	}
	hostMappingList, err := loadMappingFile(*mappingFile)
	if err != nil {
		log.Fatalf("failed to load mapping file: %v", err)
	}

	var hostConf *hostConfig
	for _, h := range hostMappingList.Hosts {
		if h.Name == *host {
			hostConf = &h
			break
		}
	}
	if hostConf == nil {
		log.Fatalf("no host config mapping found for host '%s'", *host)
	}

	if *ip != "" {
		hostConf.IP = *ip
	}

	client, err := routeros.DialTLS(fmt.Sprintf("%s:%d", hostConf.IP, hostConf.Port), hostConf.Username, hostConf.Password, &tls.Config{InsecureSkipVerify: hostConf.TLSInsecureSkipVerify})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to DialTLS: %v", err))
	}
	defer client.Close()

	tz := showTime(client)
	if tz == "" {
		tz = "Europe/Zurich"
	}

	loc, _ := time.LoadLocation(tz)
	t := time.Now().In(loc)

	_, err = client.RunArgs([]string{
		"/system/clock/set",
		fmt.Sprintf("=date=%s=", t.Format("Jan/02/2006")),
		fmt.Sprintf("=time=%s=", t.Format("15:04:05")),
	})
	if err != nil {
		log.Fatal(fmt.Errorf("failed to /system/clock/set: %v", err))
	}

	showTime(client)
}

func showTime(client *routeros.Client) string {
	res, err := client.Run("/system/clock/getall")
	if err != nil {
		log.Fatal(fmt.Errorf("failed to /system/clock/getall: %v", err))
	}

	if len(res.Re) > 0 {
		m := res.Re[0].Map
		fmt.Printf("Device time: %s %s %s\n", m["date"], m["time"], m["time-zone-name"])
		return res.Re[0].Map["time-zone-name"]
	}
	return ""
}
