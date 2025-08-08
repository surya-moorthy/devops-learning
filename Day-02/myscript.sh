#!/usr/bin/bash


array=("Hello world" Hello world "dandan dandan")

ELEMENTS=${#array[@]}

for ((i=0;i<$ELEMENTS;i++)); do 
    echo ${array[${i}]}
done 

declare -a ARRAY

#creating a file descriptor for reading

exec 10<&0

exec < $1

while read LINE; do 
     ARRAY[$count]=$LINE
     ((count++))
done

echo Number of elements : ${#ARRAY[@]}

echo ${ARRAY[@]}

exec 0<&10 10<&-
