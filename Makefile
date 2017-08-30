all: src/kotori/main.go
	cd src/kotori; go build -o ../../kotori; cd -