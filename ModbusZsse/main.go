package main

import (
	"fmt"
	"ModbusZsse"
)

func main() {
	
	ans := ModbusZsse.Init(4, 200, 1, "R", 10, "192.168.1.94:502") // as input (regType_in uint, raddr_in uint, uaddr_in uint, operation_in string, value_in uint, dest_in string)
	fmt.Println("goooood: ", ans)
}