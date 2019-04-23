Here will be two scripts firstRunScript and main. FirstRunScript is obliged to get config from server(mainflux)
and creates config.yml with below shown format:

* ip_address : "192.168.1.54" 
* device_ids : [ "de7fe372-1e1c-40ea-b1b4-c2df96db2677"] 
* device_tokens :  [ "a5f7d893-3562-4552-bbe9-2b7f1dbf2813"] 
* channel_id : "e297a164-0277-4c47-903a-11913b162e85" 
* name_prefix :   "ch_1" 

Afterthat, it will run main.go which reads config.yml and sends data in exery 7 seconds with thingsID and thingsToken 
regarding to the config.yml
