BUILDDIR=./build
BIN=lbdl
PACKAGE=github.com/jeffchannell/lbdl/main
VERSION=`git describe --abbrev=0 --tags`'-'`git rev-parse --short HEAD`

all: init linux #win

clean:
	rm -rf ${BUILDDIR}

init:
	mkdir -p ${BUILDDIR}/downloads
	touch ${BUILDDIR}/downloads/.keep
	mkdir -p ${BUILDDIR}/torrents
	touch ${BUILDDIR}/torrents/.keep
	touch ${BUILDDIR}/magnet.list

linux:
	GOOS=linux GOARCH=amd64 go build -ldflags "-linkmode external -extldflags -static" -o ${BUILDDIR}/${BIN}.x86_64.linux ${PACKAGE}/arch/linux
	#GOOS=linux GOARCH=arm64 go build -ldflags "-linkmode external -extldflags -static" -o ${BUILDDIR}/${BIN}.arm64.linux ${PACKAGE}/arch/linux
	#GOOS=linux GOARCH=amd64 go build -ldflags '-X main.version='${VERSION} -o ${BUILDDIR}/${BIN}.x86_64.linux ${PACKAGE}/arch/linux
	#GOOS=linux GOARCH=arm64 go build -ldflags '-X main.version='${VERSION} -o ${BUILDDIR}/${BIN}.arm64.linux ${PACKAGE}/arch/linux

#win:
	#GOOS=windows GOARCH=386 go build -ldflags "-linkmode external -extldflags -static" -o ${BUILDDIR}/${BIN}.exe ${PACKAGE}/arch/win
	#GOOS=windows GOARCH=386 go build -ldflags '-X main.version='${VERSION} -o ${BUILDDIR}/${BIN}.exe ${PACKAGE}/arch/win
