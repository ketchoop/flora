VERSION ?= $(git describe --tags)
DIST_DIRS := find * -type d -exec

init:
	go get -u github.com/laher/goxc
	dep ensure

build:
	go build -ldflags "-X main.version=${VERSION}" cmd/flora/flora.go

build-all:
	gox -verbose \
	-ldflags "-X main.version=${VERSION}" \
	-os="linux darwin windows" \
	-arch="amd64" \
	-output="dist/{{.OS}}-{{.Arch}}/{{.Dir}}" ./cmd/flora/

pre-dist:
	go get github.com/mitchellh/gox
	mkdir -p dist

dist: pre-dist build-all
	cd dist && \
	$(DIST_DIRS) cp ../LICENSE {} \; && \
	$(DIST_DIRS) cp ../README.md {} \; && \
	$(DIST_DIRS) tar -zcf flora-${VERSION}-{}.tar.gz {} \; && \
	$(DIST_DIRS) zip -r flora-${VERSION}-{}.zip {} \; && \
	cd ..

vet:
	go tool vet .

all: bootastrap vet dist 

.PHONY: init build build-darwin build-win build-all dist vet
