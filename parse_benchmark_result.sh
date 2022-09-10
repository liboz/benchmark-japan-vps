#!/usr/bin/env bash

sudo apt install -y ripgrep

# From https://stackoverflow.com/questions/17998978/removing-colors-from-output
sed -i 's/\x1B\[[0-9;]\{1,\}[A-Za-z]//g' result.txt
rg '(I/O Speed)\(average\) : (\d*\.\d* \w*)' -or '"$1","$2"' result.txt > parsed_result.csv
rg '^ *([\w ]*), (\w*) *(\d*\.\d* \w*) *(\d*\.\d* \w*)' -or '"$1","$2","$3","$4"' result.txt >> parsed_result.csv
rg '^(Single Core)[ |]*(\d*)' -or '"$1","$2"' result.txt >> parsed_result.csv
rg '^(Multi Core)[ |]*(\d*)' -or '"$1","$2"' result.txt >> parsed_result.csv
rg -U '^[- ]*(.*) ping statistics ---\n.*received, (\d*)% packet loss.*\n.* = (\d*\.\d*)/(\d*\.\d*)/(\d*\.\d*)/(\d*\.\d*)' -or '"$1","$2","$3","$4","$5","$6"' result.txt >> parsed_result.csv