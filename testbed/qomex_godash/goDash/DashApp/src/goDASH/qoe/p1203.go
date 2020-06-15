/*
 *	goDASH, golang client emulator for DASH video streaming

 *	This program is free software; you can redistribute it and/or
 *	modify it under the terms of the GNU General Public License
 *	as published by the Free Software Foundation; either version 2
 *	of the License, or (at your option) any later version.
 *
 *	This program is distributed in the hope that it will be useful,
 *	but WITHOUT ANY WARRANTY; without even the implied warranty of
 *	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *	GNU General Public License for more details.
 *
 *	You should have received a copy of the GNU General Public License
 *	along with this program; if not, write to the Free Software
 *	Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA
 *	02110-1301, USA.
 */

package qoe

import (
	"bytes"
	"fmt"
	glob "goDASH/global"
	"goDASH/logging"
	"goDASH/utils"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

const test = "sda"
const linux = "linux"
const darwin = "darwin"

// output strings for the main body of the P.1203 json file
const bodyPrintStringHeader = "{\n    \"I11\": {\n        \"segments\": [\n            { \"bitrate\": 192, \"codec\": \"aac\", \"duration\": %s, \"start\": 0 }\n        ],\n        \"streamId\": 42\n    },\n    \"I13\": {\n        \"segments\": [\n"
const bodyPrintString = "            {\n                \"bitrate\": %s,\n                \"codec\": \"%s\",\n                \"duration\": %s,\n                \"fps\": %s,\n                \"resolution\": \"%s\",\n                \"start\": %s\n            }"
const bodyPrintStringTail = "\n        ],\n        \"streamId\": 42\n    },\n"

// head and tail parts of the stall section of the P.1203 json file
const stallHead = "    \"I23\": {\n        \"stalling\": ["
const stallTail = "],\n        \"streamId\": 42\n    },\n"

// device details of the P.1203 json file
const deviceString = "    \"IGen\": {\n        \"device\": \"pc\",\n        \"displaySize\": \"1920x1080\",\n        \"viewingDistance\": \"150cm\"\n    }\n}"

// GetOS : return a string equating to the current runtime operating system
func GetOS() string {

	return runtime.GOOS
}

// createP1203 : create the P1203 value
func createP1203(log *map[int]logging.SegPrintLogInformation, c chan float64) {

	// for each of the logs, lets create a P.1203 compliant Json file
	// get the body
	bodyString := createP1203body(*log)
	// get the stalls
	stallString := createP1203stalls(*log)
	// add all the output together
	jsonString := strings.Join([]string{bodyString, stallString, deviceString}, "")

	// write the output to a json file (file for the last map in the log)
	createP1203file(*log, jsonString)

	// calculate the P1203 value and return to the channel
	c <- getP1203Val(*log)

}

// getP1203Val : return the P1203 value for this segment
func getP1203Val(logMap map[int]logging.SegPrintLogInformation) (p1203Val float64) {

	// needed for the read from json file
	logSize := len(logMap)
	fileInput := logMap[logSize].SegmentFileName + ".json"

	// check if file exists
	if _, err := os.Stat(fileInput); err != nil {
		// input segment file does not exist, stop the app
		fmt.Println("*** The file " + fileInput + " does not exist or cannot be found.  please check if correct path is used ***")
		// stop the app
		utils.StopApp()
	}

	// calculate the P1203 value - P1203 knows from looking at the generated json file what mode to use
	out, err := exec.Command("bash", "-c", "python3 -m itu_p1203 --print-intermediate "+fileInput+" 2> /dev/null | tail -n 6 | head -n1 | cut -f 1 -d ',' | cut -f 4 -d ' '").Output()
	if err != nil {
		log.Fatal(err)
	}

	// get the P1203 value and remove any return carrige
	stringVal := strings.TrimSuffix(string(out), "\n")

	// save this string as a float64
	p1203Val, err = strconv.ParseFloat(stringVal, 64)
	if err != nil {
		log.Fatal(err)
	}

	return
}

// createP1203file : create a P1203 json file for the last downloaded segment
func createP1203file(log map[int]logging.SegPrintLogInformation, jsonString string) {

	// write the output to a json file (file for the last map in the log)
	// needed for the write to file
	logSize := len(log)
	fileLocation := log[logSize].SegmentFileName + ".json"

	// create the file to the provided file location
	out, err := os.Create(fileLocation)
	if err != nil {
		fmt.Println("*** " + fileLocation + " cannot be created ***")
		// stop the app
		utils.StopApp()
	}
	defer out.Close()

	// Write the jsonString to file
	_, err = out.Write([]byte(jsonString))
	if err != nil {
		fmt.Println("*** " + fileLocation + " cannot be saved ***")
		// stop the app
		utils.StopApp()
	}
}

// createP1203body : create the body of the json file
func createP1203body(log map[int]logging.SegPrintLogInformation) (bodyValues string) {

	var bodyVal string

	// for each of the logs, lets create a P.1203 compliant Json file
	for a := 1; a <= len(log); a++ {

		// needed for main body
		kbps := fmt.Sprintf("%.2f", log[a].P1203Kbps)
		codec := log[a].RepCodec
		segmentDuration := fmt.Sprintf("%.1f", float64(log[a].SegmentDuration))
		fps := fmt.Sprintf("%.1f", float64(log[a].RepFps))
		resolution := strconv.Itoa(log[a].RepWidth) + "x" + strconv.Itoa(log[a].RepHeight)
		start := fmt.Sprintf("%.1f", float64(log[a].PlayStartPosition/glob.Conversion1000)-float64(log[a].SegmentDuration))

		// local val
		var bodyLoop string

		// if we have been in this loop before, we need to add a comma to the end of the string
		if len(log) > 1 && len(log) != a {
			bodyLoop = bodyPrintString + "%s\n"
			// get the body values
			bodyVal = fmt.Sprintf(bodyLoop, kbps, codec, segmentDuration, fps, resolution, start, ",")
		} else {
			bodyVal = fmt.Sprintf(bodyPrintString, kbps, codec, segmentDuration, fps, resolution, start)
		}

		// save them to our string
		bodyValues = strings.Join([]string{bodyValues, bodyVal}, "")
	}

	// needed for audio header
	clipDuration := strconv.Itoa(log[len(log)].PlayStartPosition / glob.Conversion1000)

	//audioPrintHeader
	AudioHeaderVal := fmt.Sprintf(bodyPrintStringHeader, clipDuration)

	// add all the stall values to the json head and tail
	return strings.Join([]string{AudioHeaderVal, bodyValues, bodyPrintStringTail}, "")
}

// createP1203stalls : create the stall string
func createP1203stalls(log map[int]logging.SegPrintLogInformation) (stallValues string) {

	// for each of the logs, lets create a P.1203 compliant Json file
	for a := 1; a <= len(log); a++ {

		// needed for stalls
		stallTime := fmt.Sprintf("%.3f", float64(log[a].PlaybackTime/glob.Conversion1000))
		// output is an int
		// stallDuration := fmt.Sprintf("%.3f", float64(utils.Abs(log[a].StallTime)/glob.Conversion1000))
		// output is a float
		stallDuration := fmt.Sprintf("%.3f", (float64(utils.Abs(log[a].StallTime)) / float64(glob.Conversion1000)))

		// local val
		var stallLoop string

		// different choices for stall outputs
		switch a {
		// a == 1
		case 1:
			stallLoop = "[%s,%s]"
			// get the stall values
			stallVal := fmt.Sprintf(stallLoop, stallTime, stallDuration)
			// save them to our string
			stallValues = strings.Join([]string{stallValues, stallVal}, "")
		// default for all other cases
		default:
			// we only want stalls if there is a stall time
			if utils.Abs(log[a].StallTime) > 0 {
				stallLoop = ",[%s,%s]"
				stallVal := fmt.Sprintf(stallLoop, stallTime, stallDuration)
				stallValues = strings.Join([]string{stallValues, stallVal}, "")
			}
		}
	}
	// add all the stall values to the json head and tail
	return strings.Join([]string{stallHead, stallValues, stallTail}, "")
}

// GetKBPS : return the kbps value for this segment
func GetKBPS(fileInput string, segDuration int64, debugLog bool, isByteRangeMPD bool, segSize int) (kbpsFloat float64) {

	// if this is a byte-range semgent, save the segment duration to withoutHeaderVal
	withoutHeaderVal := int64(segSize)

	// if this is not a byte-range semgent, calcualte the withoutHeaderVal
	if !isByteRangeMPD {

		// get the correct version of grep based on O/S
		var grep string
		var fi os.FileInfo
		var err error

		// check if file exists
		if fi, err = os.Stat(fileInput); err != nil {
			// input segment file does not exist, stop the app
			fmt.Println("*** The" + fileInput + " segment does not exist or cannot be found.  please check if correct path is used ***")
			// stop the app
			utils.StopApp()
		}

		// set the version of grep we will use
		switch GetOS() {
		// linux use grep
		case linux:
			// size in bytes using grep
			grep = "grep"
			logging.DebugPrint(glob.DebugFile, debugLog, "\nDEBUG: ", "grep being used on Linux")
		// mac use grep
		case darwin:
			// size in bytes using ggrep (from brew or port)
			grep = "ggrep"
			logging.DebugPrint(glob.DebugFile, debugLog, "\nDEBUG: ", "ggrep being used on OSX")
		}

		// get the location of this hex value in the input file - return 2 positions
		out, err := exec.Command("bash", "-c", "LANG=C "+grep+" -obUaP \"\\x00\\x00\\x00\\x04\\x68\\xEF\\xBC\\x80\" "+fileInput+" | awk 'BEGIN{FS=\":\"}{print $1}'").Output()
		if err != nil {
			log.Fatal(err)
		}

		// sometimes we can't read the hex value from the input segment
		// in this instance we just use the entire segment size as input to P.1203
		if !(len(out) == 0) {
			logging.DebugPrint(glob.DebugFile, debugLog, "\nDEBUG: ", "P1203 has the correct hex value ")

			// get the index of the first return carrige
			returnIndex := bytes.IndexByte(out, 10)
			// save this value as a string
			mdatValue := string(out[:returnIndex])
			// convert this value to an int64
			mdatValueInt, err := strconv.ParseInt(mdatValue, 10, 64)
			if err != nil {
				log.Fatal(err)
			}
			// add 8 bits for header
			mdatValueInt += 8
			// get the file byte size less the header
			withoutHeaderVal = fi.Size() - mdatValueInt
		}
	}

	// determine the bitrate based on segment duration - multiply by 8 and divide by segment duration
	kbpsInt := ((withoutHeaderVal * 8) / segDuration)
	// convert kbps to a float
	kbpsFloat = float64(kbpsInt) / glob.Conversion1024

	kbpsFloatStringVal := fmt.Sprintf("%3f", kbpsFloat)

	logging.DebugPrint(glob.DebugFile, debugLog, "\nDEBUG: ", "P1203 bitrate is "+kbpsFloatStringVal)

	return
}
