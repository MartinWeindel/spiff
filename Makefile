VERBOSE=-v

all: grammar test release

grammar:
	go get github.com/pointlander/peg
	(cd $(GOPATH)/src/github.com/pointlander/peg; git checkout 1d0268dfff9bca9748dc9105a214ace2f5c594a8; go install .)
	peg dynaml/dynaml.peg

release: spiff_linux_amd64.zip spiff_darwin_amd64.zip	

test: ensure
	go test $(VERBOSE) ./...

spiff_linux_amd64.zip: ensure
	GOOS=linux GOARCH=amd64 go build -o spiff++/spiff++ .
	(cd spiff++; zip spiff_linux_amd64.zip spiff++)
	rm spiff++/spiff++

ensure:
	dep ensure $(VERBOSE)
	# patch yaml parser to resolve conflict with spiff merge operator `<<:`
	sed -i.bak 's/n.value == "<<"/n.value == "<<<<"/' vendor/gopkg.in/yaml.v2/decode.go

spiff_darwin_amd64.zip: ensure
	GOOS=darwin GOARCH=amd64 go build -o spiff++/spiff++ .
	(cd spiff++; zip spiff_darwin_amd64.zip spiff++)
	rm spiff++/spiff++
