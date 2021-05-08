build:
	rm ./cf.exe
	go build -ldflags "-s -w" cf.go
