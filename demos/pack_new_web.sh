#!/bin/bash

set -e

cp -f ../app/app.toml new_web/conf

tar zcf new_web.tar.gz -C new_web/ .

echo `base64 -i new_web.tar.gz`>"new_web_data.txt"

rm -rf new_web.tar.gz