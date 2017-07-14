#!/bin/sh
#FTPHOST=
#FTPUSER=
#FTPPASS=
FTPDIR=falcon
FILE=falcon.tar.gz

tar czvf ${FILE} *

## Send updates to server
ftp -nv <<EOF
	open $FTPHOST
	user $FTPUSER $FTPPASS
	binary
	cd $FTPDIR
	put $FILE
	quit
EOF
