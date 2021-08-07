#!/usr/bin/env bash

rm -rf output
mkdir -p output/templates
mkdir -p output/html
mkdir -p output/books
cp script/* output/
cp templates/* output/templates/
cp books/* output/books/
chmod +x output/bootstrap.sh

go build -o output/
