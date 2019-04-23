package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
)

type ResponseJson struct {
	MainfluxID       string     `json:"mainflux_id,omitempty"`
	MainfluxKey      string     `json:"mainflux_key,omitempty"`
	MainfluxChannels []Channels `json:"mainflux_channels,omitempty"`
	Content          string     `json:"content,omitempty"`
}

type Channels struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Metadata string `json:"metadata,omitempty"`
}

func getMacAddr() ([]string, error) {
	// Here we will find mac-addresses of the device those are owned as external key and id
	ifas, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var as []string
	for _, ifa := range ifas {
		a := ifa.HardwareAddr.String()
		if a != "" {
			as = append(as, a)
		}
	}
	return as, nil
}

func updateConfig(nameFile string, responseFromServer []byte) (err error) {
	// Here we will create the config.yml which needs for main.go script
	r := bytes.NewReader(responseFromServer)
	decoder := json.NewDecoder(r)
	val := &ResponseJson{}
	er := decoder.Decode(val)
	if er != nil {
		log.Fatal(err)
	}

	f, err := os.Create(nameFile)
	if err != nil {
		fmt.Println(err)
		return
	}
	deviceID := val.MainfluxID
	deviceToken := val.MainfluxKey
	channelID := val.MainfluxChannels[0].ID
	name := val.MainfluxChannels[0].Name
	//written string is the body of config.yml which are fulfilled by response of GET request to the Mainflux
	writtenString := "ip_address : \"192.168.1.54\" \ndevice_ids : [ \"" + deviceID + "\"] \n" + "device_tokens :  [ \"" + deviceToken + "\"]\n" + "channel_id : \"" + channelID + "\" \n" + "name_prefix :   \"" + name + "\" "
	l, err := f.WriteString(writtenString)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return err
	}
	fmt.Println(l, "bytes written successfully")
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func main() {

	as, err := getMacAddr()
	if err != nil {
		log.Fatal(err)
	}
	externalId := as[0]
	externalKey := as[1]
	fmt.Println(externalId, externalKey)
	//payload := strings.NewReader("[{\"bn\":\"some-base-name\",\"bt\":48.56,\"n\":\"voltage\", \"u\":\"V\", \"v\" : 5.6}]")
	req, _ := http.NewRequest("GET", "http://192.168.1.54:8200/things/bootstrap/"+externalId, nil)
	req.Header.Add("Authorization", externalKey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("There is a problem on http connection")
	} else {
		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("There is problem reading res body")
		} else {

			fmt.Println(string(body))
			err := updateConfig("config.yml", body)

			if err != nil {
				fmt.Println(err)
			}
			cmd1 := exec.Command("/bin/bash", "-c", "go run main.go")
			fmt.Println("Starting command")
			// Runs the main.go script which reads config.yml and sends data with token and id that is written in config.yml
			err1 := cmd1.Run()
			fmt.Println("DONE permission: ", err1)
		}
	}

}
