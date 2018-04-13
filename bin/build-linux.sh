#!/bin/bash
# usage ./build-linux.sh 1.0.13
curDir=`pwd`
GREEN="\e[1;32m"
RESET="\e[0m"
echo -e "${GREEN} current dir is $curDir ${RESET}"

d=`date "+%Y-%m-%d-%H-%M-%S"`
echo -e "${GREEN} mkpkg_time is $d ${RESET}"

pkg_version=$1
echo -e "${GREEN} pkg_version is $pkg_version ${RESET}" 

cd   $curDir/..
STRING_GAME=`git log | head -n 1 | awk '{print $2}'`
echo -e "${GREEN} game commit is $STRING_GAME ${RESET}" 

cd   $curDir/..
echo `pwd`

cd   $curDir/..

export GOPATH=`pwd`
export GOARCH=amd64
export GOOS=linux
cd bin


go build  -o game -ldflags "-X main.Commit=$STRING_GAME -X 'main.BUILD_TIME=`date`' -X main.VERSION=$1 -s -w"  ../src/main.go


read -p "Press any key to continue." var