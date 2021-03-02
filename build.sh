export CGO_ENABLED=0
go get github.com/mitchellh/gox

mkdir bin
rm bin/* -fr
cd bin

#gox -output="fuck 1984_{{.OS}}_{{.Arch}}" -osarch="linux/amd64" -osarch="windows/386" ..
gox -output="fuckwsl2" -osarch="windows/386" ..
cd ..
