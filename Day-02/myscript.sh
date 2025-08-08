#!/usr/bin/bash

#signal interrupt used to get CTRL+C signal to set a trap on it like exiting the code.
trap bashtrap SIGINT

clear;

bashtrap() {
  echo "CTRL+C Detected !...executing bash trap!"

  exit 1;
}

for a in `seq 1 10`; do
     echo "$a/10 to Exit."
     sleep 0.5;
done 

echo "Exit Bash trap Example!!"

