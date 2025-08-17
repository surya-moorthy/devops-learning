#!/usr/bin/bash

choice=4

echo "1.Bash"
echo "2.Scripting"
echo "3.Tutorial"
echo "4.Exit"
echo -n "Enter your choice:"

while [ $choice -eq 4 ]; do

read choice
 
if [ $choice -eq 1 ]; then
  
 echo "You want to enter into bash"

else 
    if [ $choice -eq 2 ]; then 
       echo "You want to Scripting Huh!"

    else [ $choice -eq 3 ]; 
       echo "You want to watch tutorial!"
       if [ $choice -eq  4 ]; then
          exit
        fi
    fi
fi

done

