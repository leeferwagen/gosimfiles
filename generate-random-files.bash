#!/bin/bash

rm -rf /tmp/simfile
mkdir -p /tmp/simfile

COUNT=1
LEN=2048
while [ "$COUNT" -lt 300 ]
do
	OFILE=/tmp/simfile/file-`printf "%04d" "$COUNT"`
	echo $OFILE
	openssl rand -base64 $LEN > $OFILE
	COUNT=$(( $COUNT + 1 ))
done
