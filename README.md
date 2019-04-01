# pika
NY Cab trip data wrapped in a API to make it more useful

## Objective
- Provide an API wrapper for NY Cab trip data
- Provide a terminal user interface client to query the API

## Steps
- The project is built with **Go version 1.11** and is completely dockerized. There are docker containers for Go with the API code and MySQL with the db scripts. The project uses **go.mod** for dependencies.
- I used the commands to build and start the application
``` docker-compose build```
```docker-compose up```
- The project also has a client package which is built with only standard Go libraries. To run the client cd to the client package and use the command ```go run client.go```.
- I have developed this using VSCode in Ubuntu 16.04.6 LTS

## Client
- After running the client the terminal would request the user to choose from options "**1 Query**" or "**2 Clear Cache**". Entering 1 would start the Query, while entering 2 would perform the Clear Cache
- If 1 is entered and hence Query is chosen to be executed, the user would then be presented with the next option "**Enter Trip Date in dd/mm/yyyy format**"
- Once the date is entered, the user would then be prompted to enter medallion(s) using "**Enter 1 or more Medallions in a comma separated format, without spaces**"
- After medallions have been entered the next option "**Ignore cached data? Y to ignore or N to use cache**" would show up. Depending on the input the caching mechanism will either be used or ignored. 
- If the query is successful the server would respond with a json shown on the terminal in the format given below
```{"trips":[{"medallion":"D7D598CD99978BD012A87A76A7C891B7","count":3},{"medallion":"5455D5FF2BD94D10B304A15D4B7F2735","count":2},{"medallion":"801C69A08B51470871A8110F8B0505EE","count":1}]}```
- In case it responds with ```{"error":"dial tcp 152.28.1.1:3306: connect: connection refused"}```, please try to restart the server.
![Alt text](/Client_Terminal_Screenshot.png?raw=true "Client Terminal")
