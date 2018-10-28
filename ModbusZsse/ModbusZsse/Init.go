package ModbusZsse

import (
	"fmt"
	"flag"
	"reflect"
)

func Init(regType_in uint, raddr_in uint, uaddr_in uint, operation_in string, value_in uint, dest_in string) ([]int){
	
	fmt.Println("Modbus Client Started")
	regType := flag.Uint("t", regType_in ,"Register type")
	raddr := flag.Uint("r", raddr_in ,"Starting address")
	uaddr := flag.Uint("u", uaddr_in ,"Unit address")
	operation := flag.String("o",operation_in, "Read or Write operation")
	value := flag.Uint("v", value_in ,"Value to Write for write requests")
	dest :=  flag.String("dst",dest_in, "Destination address")
	flag.Parse()
	requestHandler := ModbusRequest{
			*regType,
			*raddr,
			*uaddr,
			*operation,
			*value,
			*dest,
		}
		fmt.Println("Called correctly")  //<<<<<<<<  delete it 

		
		result, err := requestHandler.Handlerequest()
		checkError(err)
		var need []int
		if(result != nil) {
			//fmt.Println("dsafasf: ", reflect.TypeOf(result))
			//fmt.Println("dsafasf: ", result)
			
			right := result.([]uint16)

			need := make([]int, len(right))

			for i:= 0 ; i < len(right); i++ {
				if (right[i]> 32767){
					need[i] = int(right[i]) - 65536
				}else{
					need[i] = int(right[i])
				}
			}
			//fmt.Println("dsafasf: ", need)
			return need
		}
		return need
	
}


