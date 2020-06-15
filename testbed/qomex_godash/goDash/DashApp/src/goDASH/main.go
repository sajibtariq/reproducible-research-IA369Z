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

package main

import (
	//to read inputs
	"encoding/json"
	"flag"
	"fmt" // to read arguments to application
	glob "goDASH/global"
	"goDASH/http"
	"goDASH/logging"
	"goDASH/player"
	"goDASH/utils"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

// variable to determine if debug log string will print
var debugLog = false

// variable to determine if hls is on
var hlsBool = false

// variable to determine if hls is on
var adaptSegDurBool = false

// variable to determine if we will print to terminal
var printLog = false

// should we extend the print over additional columns
var extendPrintLog = false

// variable to determine if quic is on
var quicBool = false

// variable to determine if we are using the goDASHbed testbed
var useTestbedBool = false

// variable to determine if getHeader is on
var getHeaderBool = false
var getHeaderReadFromFile string

// variable to determine if we are using the config file
var configSet = false

// variable to determine if QoE is on
var getQoEBool = false

// where to save the downloaded files
var fileDownloadLocation = "../files/"

// slices for our encoders, algorithms and HLS
var codecSlice = []string{glob.RepRateCodecAVC, glob.RepRateCodecHEVC, glob.RepRateCodecVP9, glob.RepRateCodecAV1}
var algorithmSlice = []string{glob.ConventionalAlg, glob.ElasticAlg, glob.LogisticAlg, glob.TestAlg, glob.ProgressiveAlg, glob.MeanAverageAlg, glob.GeomAverageAlg, glob.EMWAAverageAlg, glob.ArbiterAlg, glob.BBAAlg}
var hlsSlice = []string{glob.HlsOff, glob.HlsOn}

// default value for the exponential ratio
var exponentialRatio = 0.0

// byte range boolean
var isByteRangeMPD = false

// dictionary for printHeaders
var printHeadersData map[string]string

// set up the debug log file
func init() {
	// create this new folder location and create the log file
	os.MkdirAll(glob.DebugFolder, os.ModePerm)
	// utils.WriteFile(glob.DebugFile)
}

// main function
func main() {

	//to print a message in case there would be an error and stop the application
	defer utils.RecoverPanic()

	os.Setenv("VERSION", "0.6.2")

	var structList []http.MPD

	// creating the flag structure of the help output
	// this sets each flag
	urlPtr := flag.String(glob.URLName, "", "a list of urls specifying the location of the video clip MPD files - \"[<url>,<url>]\"")
	configPtr := flag.String(glob.ConfigName, "", "config file for this video stream - \"[path/to/config/file]\" - values in the config file have precedence over all parameters passed via command line")
	debugPtr := flag.String(glob.DebugName, glob.DebugOff, "set debug information for this video stream - \"["+glob.DebugOn+"|"+glob.DebugOff+"]\"")
	codecPtr := flag.String(glob.CodecName, glob.RepRateCodecAVC, "video codec to use - used when accessing multi-codec MPD files - \"["+glob.RepRateCodecAVC+"|"+glob.RepRateCodecHEVC+"|"+glob.RepRateCodecVP9+"|"+glob.RepRateCodecAV1+"]\"")
	maxHeightPtr := flag.Int(glob.MaxHeightName, 2160, "maximum height resolution to stream - defaults to maximum resolution height in MPD file")
	streamDurationPtr := flag.Int(glob.StreamDurationName, 0, "number of seconds to stream - defaults to maximum stream duration in MPD file")
	maxBufferPtr := flag.Int(glob.MaxBufferName, 30, "maximum stream buffer in seconds")
	initBufferPtr := flag.Int(glob.InitBufferName, 2, "initial number of segments to download before stream starts")
	adaptPtr := flag.String(glob.AdaptName, glob.ConventionalAlg, "DASH algorithms - \""+glob.ConventionalAlg+"|"+glob.ElasticAlg+"|"+glob.ProgressiveAlg+"|"+glob.LogisticAlg+"|"+glob.MeanAverageAlg+"|"+glob.GeomAverageAlg+"|"+glob.EMWAAverageAlg+"|"+glob.ArbiterAlg+"|"+glob.BBAAlg+"\"")
	fileStorePtr := flag.String(glob.FileStoreName, "", "folder location within "+fileDownloadLocation+" to store the streamed DASH files - if no folder is passed, output defaults to \"../files\" folder")
	terminalPrintPtr := flag.String(glob.TerminalPrintName, glob.TerminalPrintOff, "extend the output logs to provide additional information - \"["+glob.TerminalPrintOn+"|"+glob.TerminalPrintOff+"]\"")
	hlsPtr := flag.String(glob.HlsName, glob.HlsOff, "HLS setting - used for redownloading chunks at a higher quality rep_rate - \""+glob.HlsOff+"|"+glob.HlsOn+"\"")
	quicPtr := flag.String(glob.QuicName, glob.QuicOff, "download the stream using the QUIC transport protocol - \"["+glob.QuicOn+"|"+glob.QuicOff+"]\"")
	expRatioPtr := flag.Float64(glob.ExpRatioName, 0, "download the stream with exponential parameter : ratio - this only works with these algorithms (XXXXXXXXX)")
	getHeaderPtr := flag.String(glob.GetHeaderName, glob.GetHeaderOff, "get the header information for all segments across all of the MPD urls - based on:  \"["+glob.GetHeaderOff+"|"+glob.GetHeaderOn+"|"+glob.GetHeaderOnline+"|"+glob.GetHeaderOffline+"]\" "+glob.GetHeaderOff+": do not get headers, "+glob.GetHeaderOn+": get all headers defined by MPD, "+glob.GetHeaderOnline+": get headers from webserver based on algorithm input and "+glob.GetHeaderOffline+": get headers from header file based on algorithm input (file created by "+glob.GetHeaderOn+")")
	printHeaderPtr := flag.String(glob.PrintHeaderName, "", "print columns based on selected print headers:")
	useTestbedPtr := flag.String(glob.UseTestBedName, glob.UseTestBedOff, "setup https certs and use goDASHbed testbed - \"["+glob.UseTestBedOn+"|"+glob.UseTestBedOff+"]\"")
	QoEPtr := flag.String(glob.QoEName, glob.QoEOff, "print per segment QoE values (P1203 mode 0 and Claye) - \"["+glob.QoEOn+"|"+glob.QoEOff+"]\"")
	LogFilePtr := flag.String(glob.DebugFileName, glob.DebugFile, "Location to store the debug logs")

	// nicer print out for flags details
	flag.Usage = func() {
		fmt.Println("")
		fmt.Println("Flags for " + glob.AppName + ":")
		flag.PrintDefaults()
		fmt.Println("  - help or -h\n" + "\tPrint help screen")
	}

	// parse the arguments to the application
	flag.Parse()

	// check if no arguments are passed to the application
	if len(os.Args) == 1 {
		// print error message
		fmt.Println("*** Arguments are needed ***")
		// stop the app
		utils.StopApp()
	}

	// check config is first - check the config arguement
	if utils.IsFlagSet(glob.ConfigName) {

		// we can't print anything to debug here, as we have not set debug boolean so far :)

		// check if the config file exists
		if fi, err := os.Stat(*configPtr); err == nil {
			// if the config file exists
			if !strings.HasPrefix(*configPtr, "-") && !strings.HasPrefix(*configPtr, "[") {
				// if the file is empty stop the application
				if fi.Size() < 1 {
					// print error message
					fmt.Println("*** The" + glob.ConfigName + " file is empty, add content to " + glob.ConfigName + " file or remove from app arguements ***")
					// stop the app
					utils.StopApp()
				}

				// get some new values from the config file
				configURLPtr, configAdaptPtr, configCodecPtr, configMaxHeightPtr, configStreamDurationPtr, configMaxBufferPtr, configInitBufferPtr, configHlsPtr, configFileStorePtr, configGetHeaderPtr, configDebugPtr, configTerminalPrintPtr, configQuicPtr, configExpRatioPtr, configPrintHeaderPtr, configUseTestbedPtr, configQoEPtr, configLogFilePtr := logging.Configure(*configPtr, glob.DebugFile, debugLog)

				// check for variables with no value assigned in the config file
				utils.CheckStringVal(&configURLPtr, urlPtr)
				utils.CheckStringVal(&configAdaptPtr, adaptPtr)
				utils.CheckStringVal(&configCodecPtr, codecPtr)
				utils.CheckIntVal(&configMaxHeightPtr, maxHeightPtr)
				utils.CheckIntVal(&configStreamDurationPtr, streamDurationPtr)
				utils.CheckIntVal(&configMaxBufferPtr, maxBufferPtr)
				utils.CheckIntVal(&configInitBufferPtr, initBufferPtr)
				utils.CheckStringVal(&configHlsPtr, hlsPtr)
				utils.CheckStringVal(&configFileStorePtr, fileStorePtr)
				utils.CheckStringVal(&configGetHeaderPtr, getHeaderPtr)
				utils.CheckStringVal(&configDebugPtr, debugPtr)
				utils.CheckStringVal(&configTerminalPrintPtr, terminalPrintPtr)
				utils.CheckStringVal(&configQuicPtr, quicPtr)
				utils.CheckFloatVal(&configExpRatioPtr, expRatioPtr)
				utils.CheckStringVal(&configPrintHeaderPtr, printHeaderPtr)
				utils.CheckStringVal(&configUseTestbedPtr, useTestbedPtr)
				utils.CheckStringVal(&configQoEPtr, QoEPtr)
				utils.CheckStringVal(&configLogFilePtr, LogFilePtr)

				// set our config boolean to true
				configSet = true

			} else {
				// I don't think this is going to be called
				fmt.Println("Path to the file needed after '-" + glob.ConfigName + "' : for example : ./goDASH -" + glob.ConfigName + " ../config/config")
				// stop the app
				utils.StopApp()
			}
		} else if os.IsNotExist(err) {
			// the config file does not exists
			fmt.Println("*** " + glob.ConfigName + " file does not exist, please check file location ***")
			// stop the app
			utils.StopApp()

		} else {
			// some times we can get file errors, so just stop the program
			fmt.Println("*** " + glob.ConfigName + " file can not be read properly, please check file location ***")
			// stop the app
			utils.StopApp()
		}
	}

	// set debug is the second check - check the debug argument
	if utils.IsFlagSet(glob.DebugName) || configSet {

		// we can't print anything to debug here, as we have not set debug boolean so far :)

		if *debugPtr == glob.DebugOn {
			// set the debug logger boolean to true
			debugLog = true

			// set the debug logging location
			if utils.IsFlagSet(glob.DebugFileName) || configSet {
				// if the log file is not the same as the default log file setting
				if glob.DebugFile != *LogFilePtr {
					// reset the global location for this log file
					glob.DebugFile = glob.DebugFolder + *LogFilePtr + glob.FileFormat
					// fmt.Println(glob.DebugFile)
				}
			}
			// create the log file
			utils.WriteFile(glob.DebugFile)

			// print the first debug log string to the debug log
			logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.DebugName+" set to true ")

			// if the config file was set, then only now can we print those logs to debug
			if configSet {
				// print value to debug log
				logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.ConfigName+" set to "+*configPtr)
			}

		} else if *debugPtr == glob.DebugOff {
			// set the debug logger boolean to false
			debugLog = false

		} else {
			// print error message
			fmt.Println("*** -" + glob.DebugName + " must be set to either " + glob.DebugOn + " or " + glob.DebugOff + " (" + glob.DebugOff + " by default). ***")
			// stop the app
			utils.StopApp()
		}
	}

	// set testbed is the third check - check the useTestbed argument
	if utils.IsFlagSet(glob.UseTestBedName) || configSet {

		// print the first debug log string to the debug log
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.UseTestBedName+" set to "+*useTestbedPtr)

		if *useTestbedPtr == glob.UseTestBedOn {
			// set the extend logger boolean to true
			useTestbedBool = true

		} else if *useTestbedPtr == glob.UseTestBedOff {
			// set the extend logger boolean to false
			useTestbedBool = false

		} else {
			// print error message
			fmt.Println("*** -" + glob.UseTestBedName + " must be set to either " + glob.UseTestBedOn + " or " + glob.UseTestBedOff + " (" + glob.UseTestBedOff + " by default). ***")
			// stop the app
			utils.StopApp()
		}
	}

	// set quic is the fourth check - check the quic argument
	if utils.IsFlagSet(glob.QuicName) || configSet {

		// print the first debug log string to the debug log
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.QuicName+" set to "+*quicPtr)

		if *quicPtr == glob.QuicOn {
			// set the extend logger boolean to true
			quicBool = true

		} else if *quicPtr == glob.QuicOff {
			// set the extend logger boolean to false
			quicBool = false

		} else {
			// print error message
			fmt.Println("*** -" + glob.QuicName + " must be set to either " + glob.QuicOn + " or " + glob.QuicOff + " (" + glob.QuicOff + " by default). ***")
			// stop the app
			utils.StopApp()
		}
	}

	// set url is the fifth check - check the url arguement
	if utils.IsFlagSet(glob.URLName) || configSet {

		// print value to debug log
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.URLName+" set to "+*urlPtr)

		if !strings.HasPrefix(*urlPtr, "-") {
			structList = http.ReadURLArray(*urlPtr, debugLog, useTestbedBool, quicBool)

			// save the current MPD Rep_rate Adaptation Set
			// check if the codec is in the MPD urls passed in
			codecList, codecIndexList := http.GetCodec(structList, *codecPtr, debugLog)
			// determine if the passed in codec is one of the codecs we use (checking the first MPD only)
			usedCodec, codecIndex := utils.FindInStringArray(codecList[0], *codecPtr)
			// check the codec and print error is false
			if !usedCodec {
				// print error message
				fmt.Printf("*** -" + glob.CodecName + " " + *codecPtr + " is not in the provided MPD, please check " + *urlPtr + " ***\n")
				// stop the app
				utils.StopApp()
			}

			// get the current adaptation set, number of representations and min and max index based on max resolution height
			currentMPDRepAdaptSet := codecIndexList[0][codecIndex]
			mpdLength := len(structList[0].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation)
			mpdIndex0 := structList[0].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[0].BandWidth
			mpdIndexMax := structList[0].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[mpdLength-1].BandWidth

			// if the MPD is reversed (index 0 for represenstion is the lowest rate)
			// then reverse the represenstions
			if mpdIndex0 < mpdIndexMax {

				// define a new structList
				var reversedStructList []http.MPD
				// create it with content
				reversedStructList = http.ReadURLArray(*urlPtr, debugLog, useTestbedBool, quicBool)

				// loop over the existing list and reverse the representations
				i := 0
				for j := mpdLength - 1; j >= 0; j-- {

					// save the lowest index of structList in the highest index of reversedStructList
					reversedStructList[0].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[j] = structList[0].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[i]
					// reset the ID number of reversedStructList
					reversedStructList[0].Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[j].ID = j + 1
					// increment i
					i = i + 1
				}
				//reset the structlist to the new rates
				structList[0].Periods[0].AdaptationSet[currentMPDRepAdaptSet] = reversedStructList[0].Periods[0].AdaptationSet[currentMPDRepAdaptSet]
			}

		} else {
			fmt.Println("*** A URL(s) arguement is needed for the MPD(s) location ***")
			// stop the app
			utils.StopApp()
		}
	}

	// check the QoE argument
	if utils.IsFlagSet(glob.QoEName) || configSet {

		// print the first debug log string to the debug log
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.QoEName+" set to "+*QoEPtr)

		if *QoEPtr == glob.QoEOn {
			// set the extend logger boolean to true
			getQoEBool = true

			// check if P1203 is in the system PATH
			_, err := exec.LookPath(glob.P1203exec)
			if err != nil {
				log.Fatal(glob.P1203exec + " has not been found in $PATH, either turn \"QoE off\" or make sure P1203 has been installed and added to your $PATH")
				os.Exit(3)
			}
			logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", glob.P1203exec+" is installed")

		} else if *QoEPtr == glob.QoEOff {
			// set the extend logger boolean to false
			getQoEBool = false

		} else {
			// print error message
			fmt.Println("*** -" + glob.QoEName + " must be set to either " + glob.QoEOn + " or " + glob.QoEOff + " (" + glob.QoEOff + " by default). ***")
			// stop the app
			utils.StopApp()
		}
	}

	// check the printHeaders arguement
	if utils.IsFlagSet(glob.PrintHeaderName) || configSet {

		// only unmarhsall json if parameters were passed
		if *printHeaderPtr != "" {
			err := json.Unmarshal([]byte(*printHeaderPtr), &printHeadersData)
			if err != nil {
				panic(err)
			}
		}
		// if the printHeaders map is empty
		if len(printHeadersData) > 1 {
			extendPrintLog = true
			logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.PrintHeaderName+": print additional headers "+*printHeaderPtr)
		} else {
			logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.PrintHeaderName+": print no additional headers "+*printHeaderPtr)
		}
	}

	// check the extend argument
	if utils.IsFlagSet(glob.TerminalPrintName) || configSet {

		// print the first debug log string to the debug log
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.TerminalPrintName+" set to "+*terminalPrintPtr)

		if *terminalPrintPtr == glob.TerminalPrintOn {
			// set the extend logger boolean to true
			printLog = true

		} else if *terminalPrintPtr == glob.TerminalPrintOff {
			// set the extend logger boolean to false
			printLog = false

		} else {
			// print error message
			fmt.Println("*** -" + glob.TerminalPrintName + " must be set to " + glob.TerminalPrintOn + " or " + glob.TerminalPrintOff + " (" + glob.TerminalPrintOn + " by default). ***")
			// stop the app
			utils.StopApp()
		}
	}

	// check the getHeader argument
	if utils.IsFlagSet(glob.GetHeaderName) || configSet {

		// print the first debug log string to the debug log
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.GetHeaderName+" set to "+*getHeaderPtr)

		// get all segments headers for all MPD files
		if *getHeaderPtr == glob.GetHeaderOn {
			getHeaderBool = true

			// do not get segment headers
		} else if *getHeaderPtr == glob.GetHeaderOff {
			getHeaderBool = false

			// get segment headers from webserver based on algorithm input
		} else if *getHeaderPtr == glob.GetHeaderOnline {
			getHeaderBool = false

			// get segment headers from header file based on algorithm input
		} else if *getHeaderPtr == glob.GetHeaderOffline {
			getHeaderBool = false

			// loop over all MPD urls(s)
			for mpdListIndex := 0; mpdListIndex < len(structList); mpdListIndex++ {
				// variables
				isByteRangeMPD := false
				var segmentDurationArray []int

				// determine if this MPD is byte-range
				baseURL := http.GetRepresentationBaseURL(structList[mpdListIndex], 0)
				if baseURL != glob.RepRateBaseURL {
					isByteRangeMPD = true
				}

				// get the segment duration
				if isByteRangeMPD {
					// if this is a byte-range MPD, get byte range metrics
					_, segmentDurationArray = http.GetByteRangeSegmentDetails(structList, mpdListIndex)
				} else {
					// if not, get standard profile metrics
					_, segmentDurationArray = http.GetSegmentDetails(structList, mpdListIndex)
				}
				// current segment duration for the first MPD in the url list
				segmentDuration := segmentDurationArray[0]

				// get the MPD title
				headerURL := http.GetFullStreamHeader(structList[mpdListIndex], isByteRangeMPD)
				mpdTitle := (strings.Split(headerURL, "."))[0]

				// get the profile from the MPD file
				profiles := strings.Split(structList[mpdListIndex].Profiles, ":")
				numProfile := len(profiles) - 2
				profile := profiles[numProfile]

				// create the file name
				fileName := glob.DebugFolder + strconv.Itoa(segmentDuration) + "sec_" + mpdTitle
				// if byte-range add this
				if isByteRangeMPD {
					fileName += glob.ByteRangeString
				}
				// add the tail to the file
				fileName += "_" + profile + ".csv"

				// now check if the file already exists
				_, err := os.Stat(fileName)
				if err == nil {
					logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", fileName+" already exists")
				} else {
					logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", fileName+" does not exists")
					// WHAT DO WE DO NOW IF THE FILE DOES NOT EXIST ???
				}
			}

		} else {
			// print error message
			fmt.Println("*** -" + glob.GetHeaderName + " must be set to one of these string values \"[" + glob.GetHeaderOff + "|" + glob.GetHeaderOn + "|" + glob.GetHeaderOnline + "|" + glob.GetHeaderOffline + "]\" (" + glob.GetHeaderOff + " by default). Use -h for more info ***")
			// stop the app
			utils.StopApp()
		}
	}

	// check the adaptive algorithm argument
	if utils.IsFlagSet(glob.AdaptName) || configSet {
		// print value to debug log
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.AdaptName+" set to "+*adaptPtr)

		// determine if the passed in algorithm is one of the algorithms we use
		usedAlgorithm, _ := utils.FindInStringArray(algorithmSlice, *adaptPtr)

		// check the algorithm and print error is false
		if !usedAlgorithm {
			// print error message
			fmt.Printf("*** -"+glob.AdaptName+" must be either %v and not "+*adaptPtr+" ***\n", algorithmSlice)
			// stop the app
			utils.StopApp()
		}

		if *adaptPtr == "exponential" {

			if utils.IsFlagSet(glob.ExpRatioName) || configSet {
				// logging
				s := fmt.Sprintf("%.2f", *expRatioPtr)
				logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.ExpRatioName+" set to "+s)

				exponentialRatio = *expRatioPtr
			} else {
				//if there is -adapt exponential and nothing after
				//fmt.Printf("*** - " + expRatioName + " + value between 0 and 1 required with ' -" + adaptName + " ' ***\n")
				//utils.StopApp()

				logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.ExpRatioName+" set to default value of 0.0")
				//we use the default value of 0.0
				exponentialRatio = *expRatioPtr
			}

		}
	} else {
		//if there is no -adapt but -expRatio
		if utils.IsFlagSet(glob.ExpRatioName) {
			fmt.Printf("*** -adapt exponential required with : ' -" + glob.ExpRatioName + " ' ***\n")
			utils.StopApp()
		}
	}

	// check the fileStore argument
	if utils.IsFlagSet(glob.FileStoreName) || configSet {
		// print value to debug log
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.FileStoreName+" set to "+*fileStorePtr)

		// update the location to store the downloaded DASH files
		fileDownloadLocation = filepath.Join(fileDownloadLocation, *fileStorePtr)

		// create this new folder location
		os.MkdirAll(fileDownloadLocation, os.ModePerm)

	}

	// check the max resolution height argument
	if utils.IsFlagSet(glob.MaxHeightName) || configSet {
		// print value to debug log
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.MaxHeightName+" set to "+strconv.Itoa(*maxHeightPtr))

		// the input must be a positive number
		if *maxHeightPtr < 1 {
			// print error message
			fmt.Println("*** -" + glob.MaxHeightName + " must be a positive number and not " + strconv.Itoa(*maxHeightPtr) + " ***")
			// stop the app
			utils.StopApp()
		}
	}

	// check the stream duration argument
	if utils.IsFlagSet(glob.StreamDurationName) || configSet {
		// print value to debug log
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.StreamDurationName+" set to "+strconv.Itoa(*streamDurationPtr))

		// the input must be a positive number
		if *streamDurationPtr < 0 {
			// print error message
			fmt.Println("*** -" + glob.StreamDurationName + " must be a positive number and not " + strconv.Itoa(*streamDurationPtr) + " ***")
			// stop the app
			utils.StopApp()
		}

		// first work out if we are using a byte-range MPD
		baseURL := http.GetRepresentationBaseURL(structList[0], 0)
		if baseURL != glob.RepRateBaseURL {
			isByteRangeMPD = true
		}

		// variables
		var segmentDurationArray []int
		var maxSegments int

		// get max number segments and segment duration from the first URL MPD - index 0
		if isByteRangeMPD {
			// if this is a byte-range MPD, get byte range metrics
			maxSegments, segmentDurationArray = http.GetByteRangeSegmentDetails(structList, 0)
		} else {
			// if not, get standard profile metrics
			maxSegments, segmentDurationArray = http.GetSegmentDetails(structList, 0)
		}
		// get the segment duration of the last segment (typically larger than normal)
		lastSegmentDuration := http.SplitMPDSegmentDuration(structList[0].MaxSegmentDuration)
		// current segment duration for the first MPD in the url list
		segmentDuration := segmentDurationArray[0]
		// get MPD stream duration
		mpdStreamDuration := segmentDuration*(maxSegments-1) + lastSegmentDuration
		// determine if MPD stream time is larger than streamDurationPtr othewise error and stop

		if mpdStreamDuration < (*streamDurationPtr) {

			fmt.Println("*** -" + glob.StreamDurationName + ", " + strconv.Itoa(*streamDurationPtr) + " seconds, must not be larger than the maximum MPD stream duration of " + strconv.Itoa(mpdStreamDuration) + " second ***")
			// stop the app
			utils.StopApp()
		}

		// if no values passed in for segment duration, stream the entire clip
		if *streamDurationPtr == 0 {
			*streamDurationPtr = (mpdStreamDuration * glob.Conversion1000)
		} else {
			// otherwise use the passed in segment number
			// convert this segment number to seconds
			*streamDurationPtr *= glob.Conversion1000
		}
	}

	// check the max buffer argument
	if utils.IsFlagSet(glob.MaxBufferName) || configSet {
		// print value to debug log
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.MaxBufferName+" set to "+strconv.Itoa(*maxBufferPtr))

		// the input must be a positive number
		if *maxBufferPtr < 1 {
			// print error message
			fmt.Println("*** -" + glob.MaxBufferName + " must be a positive number (in seconds) and not " + strconv.Itoa(*maxBufferPtr) + " ***")
			// stop the app
			utils.StopApp()
		}
	}

	// check the initial number of buffer segments argument
	if utils.IsFlagSet(glob.InitBufferName) || configSet {
		// print value to debug log
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.InitBufferName+" set to "+strconv.Itoa(*initBufferPtr))

		// the input must be a positive number
		if *initBufferPtr < 0 {
			// print error message
			fmt.Println("*** -" + glob.InitBufferName + " must be a positive number and not " + strconv.Itoa(*initBufferPtr) + " ***")
			// stop the app
			utils.StopApp()
		}
	}

	// check the codec argument
	if utils.IsFlagSet(glob.CodecName) || configSet {
		// print value to debug log
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.CodecName+" set to "+*codecPtr)

		// determine if the passed in codec is one of the codecs we use
		usedCodec, _ := utils.FindInStringArray(codecSlice, *codecPtr)

		// check the codec and print error is false
		if !usedCodec {
			// print error message
			fmt.Printf("*** -"+glob.CodecName+" must be either %v and not "+*codecPtr+" ***\n", codecSlice)
			// stop the app
			utils.StopApp()
		}
	}

	// check the hls argument
	if utils.IsFlagSet(glob.HlsName) || configSet {
		// print value to debug log
		logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "-"+glob.HlsName+" set to "+*hlsPtr)

		// determine if the passed in hls is one of the hls we use
		usedHLS, _ := utils.FindInStringArray(hlsSlice, *hlsPtr)

		// check hls and print error is false
		if !usedHLS {
			// print error message
			fmt.Printf("*** -"+glob.HlsName+" must be either %v and not "+*hlsPtr+" ***\n", hlsSlice)
			// stop the app
			utils.StopApp()
		} else if *hlsPtr != "off" {
			hlsBool = true
		}
	}

	// its time to stream, call the algorithm file in player.go
	player.Stream(structList, glob.DebugFile, debugLog, *codecPtr, glob.CodecName, *maxHeightPtr,
		*streamDurationPtr, *maxBufferPtr, *initBufferPtr, *adaptPtr, *urlPtr, fileDownloadLocation, extendPrintLog, *hlsPtr, hlsBool, *quicPtr, quicBool, getHeaderBool, *getHeaderPtr, exponentialRatio, printHeadersData, printLog, useTestbedBool, getQoEBool)

}
