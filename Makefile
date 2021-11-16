install:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger 
swagger:
	/home/direwolf/go/bin/swagger generate spec -o ./swagger.json --scan-models