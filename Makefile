default:
	mkdir -p ~/bin
	cc shells.c -o ~/bin/shells
	cc gitupdater.c -o ~/bin/gitupdater
	go build -o ~/bin/gocat gocat/main.go
	go build -o ~/bin/auto-hashcat auto-hashcat/main.go
	go build -o ~/bin/walker walker/main.go
	go build -o ~/bin/walker pwndb/main.go
