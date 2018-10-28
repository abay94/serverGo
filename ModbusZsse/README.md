## There is folder of package of package zsse as a module and a main.go as an example. ##
### As an input for this package we need to define: ###
 - - - -
1. Register type ( uint )
   1. DISCRETE_OUT_COIL = 0
   2. DISCRETE_IN_COIL = 1
   3. ANALOG_IN_REGISTER = 3
   4. ANALOG_HOLDING_REGISTER = 4
   - - - -
2. Starting address  ( uint )
* It is a starting byte of reading or writing on the choosen register
 - - - -
3. Unit address  ( uint )
 - - - -
4. Read or Write operation ( string )
* "R"  or "W"  --->  Read or Write
 - - - -
5. Value to Write for write requests   ( uint ) 
* Value for write operations or length for read operations
 - - - -
6. Destination address ( string )
* IP address with port e.g. "127.0.0.1:502"
