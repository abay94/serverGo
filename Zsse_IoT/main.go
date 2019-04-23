package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/objenious/senml"
	"gopkg.in/yaml.v2"
)

// Note: struct fields must be public in order for unmarshal to

type Config_http struct {
	IPAddress    string   `yaml:"ip_address"`
	DeviceIds    []string `yaml:"device_ids"`
	DeviceTokens []string `yaml:"device_tokens"`
	ChannelID    string   `yaml:"channel_id"`
	NamePrefix   string   `yaml:"name_prefix"`
}

func LoadConfiguration(file string) (Config_http, error) {
	conf := Config_http{}
	configFile, err := os.Open(file)
	defer configFile.Close()
	byteValue, _ := ioutil.ReadAll(configFile)
	errOnUnm := yaml.Unmarshal(byteValue, &conf)
	if errOnUnm != nil {
		log.Fatalf("error on unmarshal: %v", err)
	}
	return conf, err
}

func main() {

	config, _ := LoadConfiguration("config.yml")

	ipAddress := config.IPAddress
	channelID := config.ChannelID
	deviceIds := config.DeviceIds
	deviceTokens := config.DeviceTokens
	namePrefix := config.NamePrefix

	for true {
		for i := 0; i < len(deviceIds); i++ {
			val := rand.Float64() * 10
			str_val := fmt.Sprintf("%f", val)
			w := new(bytes.Buffer)
			payload := senml.Pack{{BaseName: namePrefix, Name: "temperature", Unit: "C", Value: senml.Float(val)}}
			err := json.NewEncoder(w).Encode(payload)
			fmt.Println(w)
			if err != nil {
				fmt.Println("There is problem with encoding to senml")
			}
			//payload := strings.NewReader("[{\"bn\":\"some-base-name\",\"bt\":48.56,\"n\":\"voltage\", \"u\":\"V\", \"v\" : 5.6}]")
			req, _ := http.NewRequest("POST", "http://"+ipAddress+"/http/channels/"+channelID+"/messages/", w)
			fmt.Println("Here is deviceTokens: ", deviceTokens[i])
			req.Header.Add("Authorization", deviceTokens[i])

			res, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println("There is problem on http connection")
			} else {
				defer res.Body.Close()
				body, err := ioutil.ReadAll(res.Body)
				if err != nil {
					fmt.Println("There is problem reading res body")
				} else {

					fmt.Println(string(body))
					fmt.Println(string(str_val))
				}

			}
		}
		//payload := strings.NewReader("[{\"bn\":\"some-base-name\",\"n\":\"voltage\"}, {\"n\":\"current\", \"temos\":1.1211212}]")
		// req, _ := http.NewRequest("POST", "http://"+ipAddress+"/http/channels/"+channelID+"/messages/", payload)
		// fmt.Println("Here is deviceTokens: ", deviceTokens[0])
		// req.Header.Add("Authorization", deviceTokens[0])

		// res, _ := http.DefaultClient.Do(req)
		// defer res.Body.Close()
		// body, _ := ioutil.ReadAll(res.Body)

		// fmt.Println(string(body))

		// fmt.Println("Here is the taken data: ", ipAddress)
		// fmt.Println("Here is the channel id: ", channelID)
		// fmt.Println("Here is deviceIds: ", deviceIds)
		// fmt.Println("Here is deviceTokens: ", deviceTokens)
		time.Sleep(7 * time.Second)
	}
}
