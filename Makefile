bin:
	-mkdir bin

build: bin/myframework



bin/myframework: bin
	go build -o cmd/sche_app.go bin/myframework

run:
	go run cmd/sche_app.go --master="10.8.12.174:25050" --address="10.8.12.174" --docker-image="hello-world"
