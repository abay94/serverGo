There are two type of config file 
- - - - 
1. Main
2. Tags

First one responds for attributes : 
- IP address of OPC matrikon app
- PORT of OPC matrikon app
- Connection Item
- Number of measurements ( which indicates to main.go how many tag configs need to read in order to classify all tags by measuremnts ) 
- Destination database ( in our case Influxdb ) IP address, db_name, db_port, username, password

Second one responds for attributes :
- measurement ( in which all mentioned tags in config file will be put )
- tagName  ( all tag names given in OPC matricon app )
