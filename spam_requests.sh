#!/bin/bash
x=1
while [ $x -le 25 ]
do
  curl 172.17.0.2:3333/?first=lil;
  x=$(( $x + 1 ))
done

#y=1
echo "\n\n This is the TCP test!\n"
#while [ $y -le 5 ]
#do
netcat 172.17.0.2 8080 < something.txt 
#  y=$(( $y + 1 ))
#done

echo "\n\n This is the UDP test!\n"

#dig google.com @172.17.0.2 -p 8080
dig topeos @172.17.0.2 -p 1058

#while true; do curl 172.17.0.2:3333/welcome/?first=lil; done
