#!/usr/bin/env bash
curl -L https://github.com/romanesko/awesomeProject/archive/refs/heads/master.zip -o m.zip
unzip m.zip -d .
rm m.zip
mv awesomeProject-master awesome_project
cd awesome_project
touch server.env
echo 'AWP_DB_HOST=172.17.0.1' >> server.env
echo 'AWP_DB_NAME=postgres' >> server.env
echo 'AWP_DB_USER=postgres' >> server.env
echo 'AWP_DB_PASSWORD=secret' >> server.env
