v = v1
ALL := linux32 linux64 linuxarm linuxarm64 darwin32 darwin64 move
all: $(ALL)

linux32:
	GOOS=linux GOARCH=386 go build -o dnsblast
	tar -cvzf dnsblast-$(v)-linux-386.tar.gz dnsblast
	rm dnsblast
linux64:
	GOOS=linux GOARCH=amd64 go build -o dnsblast
	tar -cvzf dnsblast-$(v)-linux-amd64.tar.gz dnsblast
	rm dnsblast
linuxarm:
	GOOS=linux GOARCH=arm go build -o dnsblast
	tar -cvzf dnsblast-$(v)-linux-arm.tar.gz dnsblast
	rm dnsblast
linuxarm64:
	GOOS=linux GOARCH=arm64 go build -o dnsblast
	tar -cvzf dnsblast-$(v)-linux-arm64.tar.gz dnsblast
	rm dnsblast
#win32:
#	GOOS=windows GOARCH=386 go build -o dnsblast.exe
#	tar -cvzf dnsblast-$(v)-windows-386.tar.gz dnsblast.exe
#	rm dnsblast.exe
#win64:
#	GOOS=windows GOARCH=amd64 go build -o dnsblast.exe
#	tar -cvzf dnsblast-$(v)-windows-amd64.tar.gz dnsblast.exe
#	rm dnsblast.exe
darwin32:
	GOOS=darwin GOARCH=386 go build -o dnsblast
	tar -cvzf dnsblast-$(v)-darwin-386.tar.gz dnsblast
	rm dnsblast
darwin64:
	GOOS=darwin GOARCH=amd64 go build -o dnsblast
	tar -cvzf dnsblast-$(v)-darwin-amd64.tar.gz dnsblast
	rm dnsblast
move:
	if [ -d builds ]; then rm -fr builds; fi
	mkdir builds
	mv *tar.gz builds
clean:
	if [ -f dnsblast ]; then rm dnsblast; fi
	if [ -f  *tar.gz ]; then rm *tar.gz; fi
	if [ -d builds ]; then rm -fr builds; fi
