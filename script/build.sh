#!/bin/bash

DIR=`dirname $0`

BINDIR=$DIR/../bin
CMDDIR=$DIR/../cmd

FILE=$1
if [[ $FILE == "" ]]; then
  FILE="simple/main"
fi

OS=$2
if [[ $OS == "" ]]; then
  OS="linux"
fi

ARCH=$3
if [[ $ARCH == "" ]]; then
  ARCH="amd64"
fi

OUTPUT=$FILE"_"$OS"_"$ARCH
if [[ $OS == "windows" ]]; then
  OUTPUT=$OUTPUT.exe
fi


echo build $CMDDIR/$FILE.go -o $BINDIR/$OUTPUT
GOOS=$OS  GOARCH=$ARCH go build -o $BINDIR/$OUTPUT $CMDDIR/$FILE.go

echo finished