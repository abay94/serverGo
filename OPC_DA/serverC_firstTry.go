package main


import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	//"github.com/360EntSecGroup-Skylar/excelize"
	"net"
	"github.com/influxdata/influxdb/client/v2"
	"log"
	"time"
	"strconv"
	"strings"
)

////////////////////////////// Struct for json files ///////////////////////////////////
type Main struct {
	Db_IP string `json:"db_IP"`
	Db_Port string `json:"db_Port"`
	Db_Name string `json:"db_Name"`
	Db_Username string `json:"db_UserName"`
	Db_UserPassword string `json:"db_UserPassword"`
	OPCda_server_IP string `json:"opcda_server_IP"`
	OPCda_server_port string `json:"opcda_server_Port"`
	OPCda_server_ConnectionItem string `json:"opcda_server_ConnectionItem"`
	Number_of_measurements string `json:"number_of_measurements"`
}

type Tag_group struct {
	Measurement string `json:"measurement"`
	Tags []Tags `json:"tags"`
}



type Tags struct {
	TagName   string `json:"tagName"`
}


////////////////////////////////////////////


////////////////////////////// Get data from OPC DA  ///////////////////////////////////

func getDataFromOPCda(tags []string, itemOPC string, ipOPC string, portOPC string, db_name string, db_ip string, db_port string, measurement string) []float64{

	resultData := make([]float64, len(tags))

	strEcho := `{"id":1, "operation":"connect", "item":"`+ itemOPC + `"}`
    servAddr := ipOPC + ":" + portOPC
    tcpAddr, err := net.ResolveTCPAddr("tcp", servAddr)
    if err != nil {
        println("ResolveTCPAddr failed:", err.Error())
        os.Exit(1)
    }
	//Make connection
    conn, err := net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        println("Dial failed:", err.Error())
        os.Exit(1)
	}
	println("connection was made")
	// Send data to server to make connection to Matrikon.OPC.Simulation.1
    _, err = conn.Write([]byte(strEcho))
    if err != nil {
        println("Write to server failed:", err.Error())
        os.Exit(1)
    }

    println("write to server = ", strEcho)

    reply := make([]byte, 2048)
	//Read the reply from the server
    _, err = conn.Read(reply)
    if err != nil {
        println("Write to server failed:", err.Error())
        os.Exit(1)
    }

    //println("reply from server=", string(reply))

	//Then we will send tags to get data
	strTags := "" 
	for i:=0; i < len(tags); i++ {
		if strTags == "" {
			strTags = `"` + tags[i] +`":"",`
		} else{
			strTags = strTags +`"`+ tags[i] +`":"",`
		}
		
		
	}
		
	
	strSend := `{"id":2,"operation":"read","items":{` + strTags +`}}`
	//println(strSend)
		_, err = conn.Write([]byte(strSend))
		
			if err != nil {
				println("Write to server failed:", err.Error())
				os.Exit(1)
			}

		replyTags := make([]byte, 4096)
		_, err = conn.Read(replyTags)
			if err != nil {
				println("Write to server failed:", err.Error())
				os.Exit(1)
			}
		println("Reply tags::", string(replyTags))
		//println("tags length:: ", len(tags))
		for i:=0; i < len(tags); i++ {
		data:= match(string(replyTags), string(tags[i]))
		
		//println("matched value ::: ", data)
		value, _ := strconv.ParseFloat(string(data), 64)
		resultData[i] = value
		//resultData = append(resultData, value)
		//println("result data ::: ", value)
		}
		fmt.Println("result data ::: ", resultData)
		
		write_influxdb(measurement, tags, resultData, db_name, db_ip, db_port)
		
	// resultData = resultData[:cap(resultData)]
	
	conn.Close()
	
	return resultData
}


func match(reply string, tagName string) string {
	i := strings.Index(reply, tagName)
	//println("index :  " , i)
	if i >= 0 {
		j := strings.Index(reply[i+len(tagName)+3:], `"`)
		//println("jjj: ", j)
		if j >= 0 {
			//println("on match function::", tagName, ": ", reply[i+3+len(tagName) : j+i+3+len(tagName)])
			return reply[i+3+len(tagName) : j+i+3+len(tagName)-1]
		}
	}
	return ""
}



////////////////////////////////////////////


//////////////////////////////  Write data to influx  ///////////////////////////////////



func write_influxdb(measurement string, tag_name []string, data []float64, db_name string, db_ip string, db_port string) string{
	
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://" + db_ip + ":" + db_port ,
		//Addr:     "http://localhost:8086",
		Username: "username",
		Password: "password",
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()


	if response, err := c.Query(client.Query{
		Command:  fmt.Sprintf("CREATE DATABASE %s", db_name),
	}); err == nil {
		if response.Error() != nil {
			response.Error()
		}
	} else {
		println(err)
	}
	
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  db_name,
		Precision: "s",
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create a point and add to batch
	tags := map[string]string{"cpu": "cpu-total"}

	
	//var fields map[string]interface{}
	fields := make(map[string]interface{})
	for i:=0; i<len(tag_name); i++{
		fields[tag_name[i]] = data[i]
	}
	fmt.Println(fields)
	pt, err := client.NewPoint(measurement, tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}
	
	// Close client resources
	if err := c.Close(); err != nil {
    		log.Fatal(err)
	}
	return "ok"

}




////////////////////////////////////////////



////////////////////////////// Function to read json ///////////////////////////////////


func readJson() {


	// Open our jsonFile
	jsonFile, err := os.Open("C_config_main.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened C_config_main.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()


	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)


	// we initialize our Users array
	var main Main


	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above
	json.Unmarshal(byteValue, &main)
	

 


	number_of_measurements, _ := strconv.Atoi(main.Number_of_measurements)
	itemOPC := main.OPCda_server_ConnectionItem
	ipOPc := main.OPCda_server_IP
	portOPC := main.OPCda_server_port
	db_name := main.Db_Name
	db_ip := main.Db_IP
	db_port := main.Db_Port





for true {

	for i := 1; i < number_of_measurements + 1; i++ {
		
		jsonFile, err := os.Open("C_config_tags_" + strconv.Itoa(i) + ".json")
		if err != nil {
			fmt.Println(err)
		}
		defer jsonFile.Close()
		byteValue, _ := ioutil.ReadAll(jsonFile)

		var tag_group Tag_group
		json.Unmarshal(byteValue, &tag_group)

		tag_names := make([]string, len(tag_group.Tags))
		for i:=0; i < len(tag_group.Tags); i++{
			tag_names[i] = tag_group.Tags[i].TagName
		}
		
		measurement := tag_group.Measurement
		getDataFromOPCda(tag_names, itemOPC, ipOPc, portOPC, db_name, db_ip, db_port, measurement)
		

	}



	time.Sleep(10 * time.Second)
}



}


////////////////////////////////////////////

func main(){
	readJson()
}
