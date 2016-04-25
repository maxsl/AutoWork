#!/bin/bash

# ./mkcrt.sh czxichen@163.com
mkdir certs
rm certs/*
echo "make server cert"
openssl req -new -nodes -x509 -out certs/server.pem -keyout certs/server.key -days 3650 -subj "/C=CN/ST=SuZhou/L=YuanQu/O=Snail Company/OU=IT/CN=www.snail.com/emailAddress=$1"
echo "make client cert"
openssl req -new -nodes -x509 -out certs/client.pem -keyout certs/client.key -days 3650 -subj "/C=CN/ST=SuZhou/L=YuanQu/O=Snail Company/OU=IT/CN=www.snail.com/emailAddress=$1"