#!/bin/bash

home=/Users/yukinomiu/go/src/hikari-go
version=`cat ${home}/script/version.txt`
client_home=${home}/command/hikari-client
server_home=${home}/command/hikari-server
target=${home}/target
config=${home}/config

# clean
echo 'clean target directory'
rm ${target}/*
echo 'clean target directory finished'

function build() {
    CGO_ENABLED=0 GOOS=$1 GOARCH=$2 go build -ldflags "-s -w" -o ${target}/$3
}

# darwin
echo 'build darwin files'
cd ${client_home}
build darwin amd64 hikari-client-darwin-x64-${version}
build darwin 386 hikari-client-darwin-x86-${version}
cd ${server_home}
build darwin amd64 hikari-server-darwin-x64-${version}
build darwin 386 hikari-server-darwin-x86-${version}
echo 'build darwin files finished'

# linux
echo 'build linux files'
cd ${client_home}
build linux amd64 hikari-client-linux-x64-${version}
build linux 386 hikari-client-linux-x86-${version}
cd ${server_home}
build linux amd64 hikari-server-linux-x64-${version}
build linux 386 hikari-server-linux-x86-${version}
echo 'build linux files finished'

# windows
echo 'build windows files'
cd ${client_home}
build windows amd64 hikari-client-windows-x64-${version}.exe
build windows 386 hikari-client-windows-x86-${version}.exe
cd ${server_home}
build windows amd64 hikari-server-windows-x64-${version}.exe
build windows 386 hikari-server-windows-x86-${version}.exe
echo 'build windows files finished'

# copy config
echo 'copy config files'
cp ${config}/* ${target}/
echo 'copy config files finished'

echo 'all finished'
