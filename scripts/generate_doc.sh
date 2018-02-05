#!/bin/bash

export GOPATH=${PWD}/gopath
src=${GOPATH}/src/github.com/yubo/falcon/ctrl/api/swagger/swagger.json
dst=${GOPATH}/src/github.com/yubo/falcon/ctrl/ctrl.swagger.json

cd ${GOPATH}/src/github.com/yubo/falcon/ctrl/api && \
mkdir -p $(dirname $dst) && \
bee generate docs && \
cp -f $src $dst
cd -

## api_reference.md
tempdir=$(mktemp -d -p .)
mkdir -p ${tempdir}/{tsdb,agent,alarm,service,transfer}
cd  ${tempdir}
cp -f ../lib/tsdb/tsdb.proto tsdb/tsdb.proto
cp -f ../agent/agent.proto agent/agent.proto
cp -f ../alarm/alarm.proto alarm/alarm.proto
cp -f ../service/service.proto service/service.proto
cp -f ../transfer/transfer.proto transfer/transfer.proto
protodoc --directories="tsdb=service_message,agent=service_message,alarm=service_message,service=service_message,transfer=service_message" \
	--title="falcon API Reference" \
	--output="api_reference.md" \
	--disclaimer="This is a generated documentation. Please read the proto files for more."
#sed -i 's/falcon\/falcon\.proto/falcon.proto/g' api_reference.md
cd -
cp -f ${tempdir}/api_reference.md ${GOPATH}/src/github.com/yubo/falcon/docs/api_reference.md
rm -rf ${tempdir}


## md -> html
render='github-markdown-render'
if ! type $render > /dev/null; then
	sudo npm install --global $render
fi

docs=$(find ./docs -name "*.md")
#docs='./docs/api_reference.md'

for doc in $docs
do
	doc=${doc#"./"}
	#doc=${doc%".md"}

	src="${doc}"
	dst="var/html/${doc}"
	echo "creating $dst"
	mkdir -p $(dirname $dst)
	echo '<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/github-markdown-css/2.10.0/github-markdown.min.css">
<style>
	.markdown-body {
		box-sizing: border-box;
		min-width: 200px;
		max-width: 980px;
		margin: 0 auto;
		padding: 45px;
	}

	@media (max-width: 767px) {
		.markdown-body {
			padding: 15px;
		}
	}
</style> </head>
<body> <article class="markdown-body">
' > $dst
	cat $src | $render >> $dst
	echo '</article></body>' >> $dst
done

