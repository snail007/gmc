#!/bin/bash

set -e

cp -f ../app/web_new.toml new_web/conf/app.toml
cp -f ../app/api_new.toml new_api/conf/app.toml

tar zcf new_web.tar.gz -C new_web/ .
tar zcf new_api.tar.gz -C new_api/ .

webdata=`base64 -i new_web.tar.gz`
apidata=`base64 -i new_api.tar.gz`

cat <<EOF >../../gmct/module/new/data.go
package newx

var (
	webData = "$webdata"
	apiData = "$apidata"
)
EOF

rm -rf new_web.tar.gz
rm -rf new_api.tar.gz