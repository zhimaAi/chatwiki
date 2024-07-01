#!/bin/bash

rm -rf dist
npm i
npm run build:win

if [ "$(command -v zip | wc -l)" -eq 0 ]
then
  apt update
  apt install -y zip
fi

cd ./dist/win-unpacked || exit
zip -r ../chatwiki.zip ./*
