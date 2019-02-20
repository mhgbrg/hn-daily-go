build:
	go build

digest: build
	./hndaily digest 2019-02-01
