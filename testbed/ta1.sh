#!/usr/bin/env bash
mod=$1_Band_data
net=$2
doc=$3
num=$4
host=$5
a=1
v=1

#end=$((SECONDS+650))

while [ $a -lt 2 ]; do
    array=()
    i=1
    j=1
    cut -d, -f 1,2 --output-delimiter=' ' ~/reproducible-research-IA369Z/testbed/Band_data/Band_data/$mod/$net/$doc$num.csv | while read col1 col2 ; do
    array[$i]=$col2
    array1[$j]=$col1

    #echo ${array[i]}
    #sleep 1
    #float=${array[i]}
    #ban=${float%.*}
    ban=$(bc<<<"scale=3; ${array[i]}*$host")
    echo $ban
    if [ $v -eq 1 ]; then

         sudo  tc qdisc add dev ap1-eth10 handle 1:0 root htb default 1 && tc class add dev ap1-eth10 parent 1:0 classid 1:1 htb rate  "$(bc<<<"scale=2; $ban/1")"kbit ceil "$(bc<<<"scale=2; $ban/1")"kbit && echo  "second" $SECONDS 'time' ${array1[j]} "band" "$(bc<<<"scale=2; $ban/1")"kbit

         sudo  tc qdisc add dev s3-eth1 handle 1:0 root htb default 1 && tc class add dev s3-eth1 parent 1:0 classid 1:1 htb rate "$(bc<<<"scale=2; $ban/1")"kbit ceil "$(bc<<<"scale=2; $ban/1")"kbit &&  echo  "second" $SECONDS 'time' ${array1[j]} "band" "$(bc<<<"scale=2; $ban/1")"kbit

         sleep ${array1[j]}

         v=0

    else

         sudo tc qdisc del dev ap1-eth10 root &&  sudo tc qdisc add dev ap1-eth10 handle 1:0 root htb default 1 && tc class add dev ap1-eth10 parent 1:0 classid 1:1 htb rate "$(bc<<<"scale=2; $ban/1")"kbit ceil "$(bc<<<"scale=2; $ban/1")"kbit &&  echo  "second" $SECONDS 'time' ${array1[j]} "band" "$(bc<<<"scale=2; $ban/1")"kbit

         sudo tc qdisc del dev s3-eth1 root  &&  sudo  tc qdisc add dev  s3-eth1 handle 1:0 root htb default 1 && tc class add dev s3-eth1 parent 1:0 classid 1:1 htb rate "$(bc<<<"scale=2; $ban/1")"kbit ceil "$(bc<<<"scale=2; $ban/1")"kbit &&  echo  "second" $SECONDS 'time' ${array1[j]} "band" "$(bc<<<"scale=2; $ban/1")"kbit

        sleep ${array1[j]}

    i=$((i + 1))
    j=$((j + 1))
    fi


    if [ $SECONDS -ge 1100 ]; then
         break
    fi

done

    a=$((a + 1))

done
