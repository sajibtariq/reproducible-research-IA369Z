#! /bin/bash

#	This program is free software; you can redistribute it and/or
#	modify it under the terms of the GNU General Public License
#	as published by the Free Software Foundation; either version 2
#	of the License, or (at your option) any later version.
#
#	This program is distributed in the hope that it will be useful,
#	but WITHOUT ANY WARRANTY; without even the implied warranty of
#	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#	GNU General Public License for more details.
#
#	You should have received a copy of the GNU General Public License
#	along with this program; if not, write to the Free Software
#	Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA
#	02110-1301, USA.
#

# this script has been tested and validated on the UHD x264 dataset across all 5 segment durations (2sec to 10sec)
#http://dev.cs.ucc.ie/misl/4K_non_copyright_dataset/10_sec/x264/sintel/DASH_Files/full/sintel_enc_x264_dash.mpd

# example call: ./true_byte_size_calc.sh ../files/347985/2sec_sintel_320x180_24fps_235kbps_segment1.m4s h264 2 24 320x180 0 1 0

# ./true_byte_size_calc.sh
# file location - ../files/347985/2sec_sintel_320x180_24fps_235kbps_segment1.m4s
# encoder - h264
# segment duration - 2
# frame rate in seconds - 24
# resolution (width x height) - 320x180
# start time - 0
# segment number - 1
# stall time - 0

# example call for 3 segments
# ./true_byte_size_calc.sh ../files/347985/2sec_sintel_320x180_24fps_235kbps_segment1.m4s h264 2 24 320x180 0 1 0
# P1203 output: 1.8775293295391329
# ./true_byte_size_calc.sh ../files/347985/2sec_sintel_320x180_24fps_235kbps_segment2.m4s h264 2 24 320x180 2 2 0
# P1203 output: 1.8775293295391324
# ./true_byte_size_calc.sh ../files/347985/2sec_sintel_640x360_24fps_1050kbps_segment3.m4s h264 2 24 640x360 4 3 0
# P1203 output: 2.3486846809025312

# make sure only 8 parameters are passed
if [ "$#" -ne 8 ]; then
    echo "Illegal number of parameters"
    echo "example input:"
fi



# file name
file_structure=$1
# codec
codec=$2
# segment duration
duration=$3
# frames per second
fps=$4
# resolution - WidthxHeight
resolution=$5
# start time of the
start=$6
# segment number - starting at XXXXX
in_number=$7
# stall time in seconds
stalling_time=$8
#
basecase="0."
#
dr=${file_structure%/*}

#touch $dr/head.json
touch $dr/tail.json
touch $dr/tailend.json
touch $dr/audio.json


# get the byte size of the segment
unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     filesize=$(stat -c%s $file_structure);; #size in bytes
    Darwin*)    filesize=$(stat -f%z $file_structure);; #size in bytes
esac

# 4 files are created before a final json file can be formed
# audio.json will hold the audio information (fake audio added as our DASH dataset has no audio)
# head.json will store the header and the segment data
# tail.json will store the stall data
# tailend.json will store the footer
file="$dr/head.json"

# get the byte size of the segment
unameOut="$(uname -s)"
case "${unameOut}" in
    Linux*)     mdatValue=`LANG=C grep -obUaP "\x00\x00\x00\x04\x68\xEF\xBC\x80" $file_structure | awk 'BEGIN{FS=":"}{print $1}'`;; #size in bytes
    Darwin*)    mdatValue=`LANG=C ggrep -obUaP "\x00\x00\x00\x04\x68\xEF\xBC\x80" $file_structure | awk 'BEGIN{FS=":"}{print $1}'`;; #size in bytes
esac

# print the first value - bit location this hex value occurs at
mdatValue=`echo $mdatValue | awk '{print $1}'`
# echo $mdatValue

# add 8 bits for header
header_val=$(($mdatValue + 8))
# echo $header_val

# get the file byte size less the header
without_header_val=$(($filesize-$header_val))
# echo $without_header_val

# determine the bitrate based on segment duration - multiply by 8 and divide by segment duration
printf -v bitrate '%.2f' "$(echo "scale=2; $without_header_val*8/$duration" | bc)"
# echo $bitrate

#printf -v will store the output in variable kbps
printf -v kbps '%.2f' "$(echo "scale=2; $bitrate/1024" | bc)"
#echo $kbps

#init is the time the stream started stalling
#start is the time the stream resumed after stalling
init=`echo "$start - $stalling_time" | bc`

#check if the head.json file exists
#this is to ensure that segment data are correctly separated by printing a comma before a segment if it is not the first segment to be written to the file
if [ -e "$file" ]
then
    # if the file exists:
    # remove the last 4 lines of the existing head.json file and
    # add a comma to the last line before we
    # write the next segment info
    for i in {1..4}
    do
        #sed -e '$ d' -i ''  $file
	      sed -i -e '$ d' $file
    done

    # convert duration and fps to floats
    printf -v durationFloat '%.1f' "$(echo "scale=2; $duration" | bc)"
    printf -v fpsFloat '%.1f' "$(echo "scale=2; $fps" | bc)"
    printf -v startFloat '%.1f' "$(echo "scale=2; $start" | bc)"

    printf -v summedFloat '%d' "$(echo "scale=2; ${start}+${duration}" | bc)"

    # add the audio duration info to the audio header.json file
    echo -e "{\n    \"I11\": {\n        \"segments\": [\n            { \"bitrate\": 192, \"codec\": \"aac\", \"duration\": $summedFloat, \"start\": 0 }\n        ],\n        \"streamId\": 42\n    },\n" > "$dr/audio.json"

    #append the relavent data to the head.json file
    echo -e "            },\n            {\n                \"bitrate\": $kbps,\n                \"codec\": \"$codec\",\n                \"duration\": ${durationFloat},\n                \"fps\": $fpsFloat,\n                \"resolution\": \"$resolution\",\n                \"start\": ${startFloat}\n            }\n        ],\n        \"streamId\": 42\n    }," >> "$file"
    #check if this segment stalled

    #	echo $stalling_time

    if [ ! -z "$stalling_time" ]
    then
        #if (( $stalling_time > 0. ))
	if (( $(echo "$stalling_time > 0" |bc -l) ))
	then

            # remove the last two lines from tail.json
            # for i in {1..2}
            # do
            #     sed -i '' -e '$ d $dr/tail.json'
            # done

            stallVals=`cat $dr/tail.json`

            STR=${stallVals#*[}
            STR=${STR%]*}

            # add the stall information to the tail.json file
            stall="    \"I23\": {\n        \"stalling\": ["
            # if there is a stall value, add this to the output
            if [ ! -z "$stalling_time" ]
            then
                if (( $(echo "$stalling_time > 0" |bc -l) )) ; then
                    if [ ! "$STR" == "" ]
                    then
                        # add stall
                        stall+="$STR,[${start},${stalling_time}]"
                    else
                        stall+="[0,0],[${start},${stalling_time}]"
                    fi
                fi
            fi

            # add tail to stall string
            stall+="],\n        \"streamId\": 42\n    }," > "$dr/tail.json"

            # add the stall information to the tail.json file
            echo -e "$stall " > "$dr/tail.json"
        fi
    fi



else
    #create the appropiate headers in the head and tail files and the appropiate footer in the tailend file

    # convert duration and fps to floats
    printf -v durationFloat '%.1f' "$(echo "scale=2; $duration" | bc)"
    printf -v fpsFloat '%.1f' "$(echo "scale=2; $fps" | bc)"
    # set the start time as a float if greater than zero
    if (( $(echo "$start > 0" |bc -l) )) ; then
        printf -v startFloat '%.1f' "$(echo "scale=2; $start" | bc)"
    else
        startFloat=0
    fi

    printf -v summedFloat '%d' "$(echo "scale=2; ${start}+${duration}" | bc)"

    # add the audio duration info to the audio header.json file
    echo -e "{\n    \"I11\": {\n        \"segments\": [\n            { \"bitrate\": 192, \"codec\": \"aac\", \"duration\": $summedFloat, \"start\": 0 }\n        ],\n        \"streamId\": 42\n    },\n" > "$dr/audio.json"
    # add the segment info to the main head.json file
    echo -e "    \"I13\": {\n        \"segments\": [\n            {\n                \"bitrate\": $kbps,\n                \"codec\": \"$codec\",\n                \"duration\": ${durationFloat},\n                \"fps\": $fpsFloat,\n                \"resolution\": \"$resolution\",\n                \"start\": ${startFloat}\n            }\n        ],\n        \"streamId\": 42\n    }," > "$file"

    # add the stall information to the tail.json file
    stall="    \"I23\": {\n        \"stalling\": ["
    # if there is a stall value, add this to the output
    if [ ! -z "$stalling_time" ]
    then
        if (( $(echo "$stalling_time > 0" |bc -l) )) ; then
            # add stall
            stall+="[${start},${stalling_time}]"
        else
            # we add this, so we can add other stalls later and not get a warning
            stall+="[0.0,0.0]"
        fi
    fi
    # add tail to stall string
    stall+="],\n        \"streamId\": 42\n    },"

    # add the stall information to the tail.json file
    echo -e "$stall" > "$dr/tail.json"

    # add the device details to the tailend.json file
    echo -e "    \"IGen\": {\n        \"device\": \"pc\",\n        \"displaySize\": \"1920x1080\",\n        \"viewingDistance\": \"150cm\"\n    }\n}" >  "$dr/tailend.json"
fi

#combine the parts into a final mode 0 compatible json file
cat "$dr/audio.json" "$file" "$dr/tail.json" "$dr/tailend.json" > "$dr/qoe-$in_number-json.json"

#run itu p.1203 (sending warnings/errors to /dev/null, effectively ignoring them) and extract the qoe only from the output
qoe=$(python3 -m itu_p1203 --print-intermediate $dr/qoe-$in_number-json.json 2>> errors.log | tail -n 6 | head -n1 | cut -f 1 -d ',' | cut -f 4 -d ' ')

# print out the qoe value
echo "$qoe"
