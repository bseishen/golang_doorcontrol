module golang_doorcontrol

go 1.24.2

require (
	github.com/bseishen/golang_doorcontrol v0.0.0-20200210182034-236d2f334874
	github.com/eclipse/paho.mqtt.golang v1.5.0
	github.com/jacobsa/go-serial v0.0.0-20180131005756-15cf729a72d4
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/mattn/go-sqlite3 v1.14.27
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.22.0 // indirect
)

replace github.com/bseishen/golang_doorcontrol => .
