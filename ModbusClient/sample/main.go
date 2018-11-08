package main

import (
	"fmt"
	"ModbusZsse"
	"encoding/json"
	"io/ioutil"
	"time"
	"os"
)

type JsonObject struct {
	Reg_type uint `json:"reg_type"`
	Start_addr uint `json:"start_addr"`
	Unit_addr uint `json:"unit_addr"`
	Rw_operation string `json:"rw_operation"`
	Value_or_length uint `json:"value_or_length"`
	Dst_address string `json:"dst_address"`
	Delay int `json:"delay"`
}

func start(config_file_name string){

	jsonFile, err := os.Open(config_file_name)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened: " + config_file_name )
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var jsonElement JsonObject
	json.Unmarshal(byteValue, &jsonElement)

	reg_type := jsonElement.Reg_type
	start_addr := jsonElement.Start_addr
	unit_addr := jsonElement.Unit_addr
	rw_operation := jsonElement.Rw_operation
	value_or_length := jsonElement.Value_or_length
	dst_address := jsonElement.Dst_address
	delay := jsonElement.Delay
	req := ModbusZsse.Init(reg_type, start_addr, unit_addr, rw_operation, value_or_length, dst_address)

	for true {

		
		ans := ModbusZsse.Run(req)
		fmt.Println("Here is the taken data: ", ans)
		time.Sleep(time.Duration(delay) * time.Second)
	}

} 

func main() {
	start("config.json")
}