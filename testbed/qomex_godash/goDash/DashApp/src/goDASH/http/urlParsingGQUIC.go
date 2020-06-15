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

package http

/*
import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"goDASH/logging"
	"goDASH/utils"

	glob "goDASH/global"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/cavaliercoder/grab"
	//"github.com/lucas-clemente-http3/quic-go/http3"

	quic "github.com/lucas-clemente/quic-go"
	//"github.com/lucas-clemente/quic-go/http3"

	"github.com/lucas-clemente/quic-go/h2quic"
	//"github.com/lucas-clemente/quic-go/internal/protocol"
)

// getURLBody :
// * get the response body of the url
// * calculate the rtt
// * return the response body and the rtt
func getURLBody(url string, isByteRangeMPD bool, startRange int, endRange int, quicBool bool, debugFile string, debugLog bool, useTestbedBool bool) (io.ReadCloser, time.Duration) {

	var client *http.Client
	var cert tls.Certificate
	var caCertPool *x509.CertPool
	var caCert []byte
	var err error
	var config *tls.Config
	var quicConfig *quic.Config
	var tr *http.Transport
	var trQuic *h2quic.RoundTripper

	if useTestbedBool {
		logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "Testbed in use")

		// Read the key pair to create certificate
		cert, err = tls.LoadX509KeyPair(glob.HTTPcertLocation, glob.HTTPkeyLocation)
		if err != nil {
			log.Println("Unable to load X509 key and cert")
			log.Fatal(err)
		}
		logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "loading X509 key and cert: "+glob.HTTPcertLocation+" "+glob.HTTPkeyLocation)

		// Create a CA certificate pool and add cert.pem to it
		caCert, err = ioutil.ReadFile(glob.HTTPcertLocation)
		if err != nil {
			log.Println("Unable to read X509 cert")
			log.Fatal(err)
		}
		logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "reading X509 cert")

		caCertPool := x509.NewCertPool()
		// add cert to pool
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "No certs appended, using system certs only")
		}
	}

	// if we want to use quic
	if quicBool {
		if !useTestbedBool {
			//roundTripper := &http3.RoundTripper{
			//	TLSClientConfig: &tls.Config{
			//		RootCAs: testdata.GetRootCA(),
			//	},
			//}
			//defer roundTripper.Close()
			//client = &http.Client{
			//	Transport: roundTripper,
			//}
		} else {
			//versions := protocol.SupportedVersions
			// set up the config
			logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "creating tls config for quic")
			quicConfig = &quic.Config{
				// use insecure SSL - if needed only use during internal tests
				// this is set statically in the globalVar.go file (set to true if needed)
				//InsecureSkipVerify: glob.InsecureSSL,
				//RootCAs:            caCertPool,
				//Certificates:       []tls.Certificate{cert},
				//Versions: versions,
			}
			// set up our http transport
			logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "creating our http transport using our tls config for quic")

			trQuic = &h2quic.RoundTripper{QuicConfig: quicConfig}
			// set up the client
			logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "creating our client using our http transport and our tls config for quic")
			client = &http.Client{Transport: trQuic}
		}
		// otherwise use a normal-ish HTTP client
	} else {
		// set up a http client
		if useTestbedBool {
			// set up the config
			logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "creating tls config")
			config = &tls.Config{
				// use insecure SSL - if needed only use during internal tests
				// this is set statically in the globalVar.go file (set to true if needed)
				InsecureSkipVerify: glob.InsecureSSL,
				RootCAs:            caCertPool,
				Certificates:       []tls.Certificate{cert},
			}
			// set up our http transport
			logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "creating our http transport using our tls config")
			tr = &http.Transport{TLSClientConfig: config}
			// set up the client
			logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "creating our client using our http transport and our tls config")
			client = &http.Client{Transport: tr}

		} else {
			client = &http.Client{}
		}
	}

	// request the url
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		fmt.Println(err)
		fmt.Println("the URL " + url + " doesn't match with anything")
		// stop the app
		utils.StopApp()
	}

	// determine the rtt for this segment
	start := time.Now()
	if useTestbedBool {
		if quicBool {
			// use our new transport for calculating the rtt
			if _, err := trQuic.RoundTrip(req); err != nil {
				log.Fatal(err)
			}
		} else {
			// use our new transport for calculating the rtt
			if _, err := tr.RoundTrip(req); err != nil {
				log.Fatal(err)
			}
		}
	} else {
		// use http default transport for calculating the rtt
		if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
			log.Fatal(err)
		}
	}
	// get rtt
	rtt := time.Since(start)

	// add the byte ranges, if byte-range
	if isByteRangeMPD {
		byteRange := "bytes=" + strconv.Itoa(startRange) + "-" + strconv.Itoa(endRange)
		req.Header.Add("Range", byteRange)
	}

	var resp *http.Response

	// if we want to use quic
	if quicBool {
		resp, err = client.Do(req)
	} else {
		//request the URL using the client
		resp, err = client.Do(req)
	}
	if err != nil {
		fmt.Println(err)
		fmt.Println("the URL " + url + " doesn't match with anything")
		// stop the app
		utils.StopApp()
	}

	// get protocol version
	protocol := resp.Request.Proto

	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "URL is : "+url)
	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "Protocol is : "+protocol)

	//Check if the GET method has sent a status code equal to 200
	if resp.StatusCode != http.StatusOK && !isByteRangeMPD {
		// add this to the debug log
		fmt.Println("The URL returned a non status okay error code: " + strconv.Itoa(resp.StatusCode))
		// stop the app
		utils.StopApp()
	}
	//fmt.Println("len : ", resp.ContentLength)

	// return the response body
	return resp.Body, rtt

}

// getURLProgressively :
// * get the response body of the url
// * calculate the rtt and throughtput for the download per second
// * return the rtt
func getURLProgressively(url string, isByteRangeMPD bool, startRange int, endRange int, fileLocation string) time.Duration {

	var thrPerSecond []int64

	// set up a http client
	client := grab.NewClient()
	// request the url and save to a file location
	req, err := grab.NewRequest(fileLocation, url)
	// if there is an error, stop the app
	if err != nil {
		fmt.Println(err)
		fmt.Println("the URL " + url + " doesn't match with anything")
		// stop the app
		utils.StopApp()
	}

	// determine the rtt for this segment
	start := time.Now()
	if _, err := http.DefaultTransport.RoundTrip(req.HTTPRequest); err != nil {
		log.Fatal(err)
	}
	// get rtt
	rtt := time.Since(start)
	//fmt.Printf("grab RTT in %dms for %s\n", rtt, url)

	// add the byte ranges, if byte-range
	if isByteRangeMPD {
		byteRange := "bytes=" + strconv.Itoa(startRange) + "-" + strconv.Itoa(endRange)
		req.HTTPRequest.Header.Add("Range", byteRange)
	}

	//request the URL using the client
	resp := client.Do(req)

	// start UI loop, (maybe we should put 1 instead of 1000 to have it in millisecond)
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()

	// Check if the download has finished or not
	//start = time.Now()
	for !resp.IsComplete() {
		select {
		case <-t.C:
			/*
				fmt.Printf("transferred %v / %v bytes (%.2f%%) in %dms\n",
					resp.BytesComplete(),
					resp.Size,
					100*resp.Progress(), time.Since(start)/1000000)
*/ /*
		thrPerSecond = append(thrPerSecond, resp.BytesComplete())

	case <-resp.Done:
		// download is complete
		/*
			fmt.Printf("transferred %v / %v bytes (%.2f%%) in %dms\n",
				resp.BytesComplete(),
				resp.Size,
				100*resp.Progress(), time.Since(start)/1000000)
*/ /*
			thrPerSecond = append(thrPerSecond, resp.BytesComplete())
			break
		}
	}
	// check for errors
	if err := resp.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Download failed: %v\n", err)
		// stop the app
		utils.StopApp()
	}

	/* We can't use this as progressive has a different status code
	//Check if the GET method has sent a status code equal to 200
	if resp.HTTPResponse.StatusCode != http.StatusOK && !isByteRangeMPD {
		// add this to the debug log
		fmt.Println("The URL returned a non status okay error code: " + strconv.Itoa(resp.HTTPResponse.StatusCode))
		// stop the app
		utils.StopApp()
	}
*/
//
// 	// return the rtt
// 	return rtt
//
// }
//
// // GetURLByteRangeBody :
// // * get the response body of the url and return an io.ReadCloser
// // * based on byte-ranges
// func GetURLByteRangeBody(url string, startRange int, endRange int) (io.ReadCloser, time.Duration) {
//
// 	// set up a http client
// 	client := &http.Client{}
// 	// request the url
// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		fmt.Println(err)
// 		fmt.Println("the URL " + url + " doesn't match with anything")
// 		// stop the app
// 		utils.StopApp()
// 	}
//
// 	//req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
// 	start := time.Now()
// 	if _, err := http.DefaultTransport.RoundTrip(req); err != nil {
// 		log.Fatal(err)
// 	}
// 	// get rtt
// 	rtt := time.Since(start)
//
// 	// add the byte ranges
// 	byteRange := "bytes=" + strconv.Itoa(startRange) + "-" + strconv.Itoa(endRange-1)
// 	req.Header.Add("Range", byteRange)
//
// 	//request the URL using the client
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		fmt.Println(err)
// 		fmt.Println("the URL " + url + " doesn't match with anything")
// 		// stop the app
// 		utils.StopApp()
// 	}
//
// 	//Check if the GET method has sent a status code equal to 200
// 	if resp.StatusCode != http.StatusOK {
// 		// add this to the debug log
// 		fmt.Println("The URL returned a non status okay error code: " + strconv.Itoa(resp.StatusCode))
// 		// stop the app
// 		utils.StopApp()
// 	}
// 	//fmt.Println("len : ", resp.ContentLength)
//
// 	// return the response body
// 	return resp.Body, rtt
//
// }
//
// // GetURL :
// // * return the content of the body of the url
// func GetURL(url string, isByteRangeMPD bool, startRange int, endRange int, quicBool bool, debugFile string, debugLog bool, useTestbedBool bool) ([]byte, time.Duration) {
//
// 	// get the response body and rtt for this url
// 	responseBody, rtt := getURLBody(url, isByteRangeMPD, startRange, endRange, quicBool, debugFile, debugLog, useTestbedBool)
//
// 	// Lets read from the http stream and not create a file to store the body
// 	body, err := ioutil.ReadAll(responseBody)
// 	//bodyString := string(body)
// 	if err != nil {
// 		fmt.Println("Unable to read from url")
// 		// stop the app
// 		utils.StopApp()
// 	}
//
// 	// close the responseBody
// 	responseBody.Close()
//
// 	// return the body of the responseBody
// 	return body, rtt
// }
//
// // GetRepresentationBaseURL :
// // * get BaseURL for byte-range MPD
// func GetRepresentationBaseURL(mpd MPD, currentMPDRepAdaptSet int) string {
// 	return mpd.Periods[0].AdaptationSet[currentMPDRepAdaptSet].Representation[0].BaseURL
// }
//
// // JoinURL :
// /*
//  * func joinURL(baseURL string, append string) string
//  *
//  * join components of urls together
//  * return the URL
//  */
// func JoinURL(baseURL string, append string, debugLog bool) string {
//
// 	// if "append" already contains "http", then do nothing
// 	if !(strings.Contains(append, "http")) {
// 		// get the base of the current url
// 		base := path.Base(baseURL)
// 		// replace this base url with the required file string
// 		urlHeaderString := strings.Replace(baseURL, base, append, -1)
// 		//logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "complete URL: "+urlHeaderString)
//
// 		// return the new url
// 		return urlHeaderString
// 	}
// 	// return the new url
// 	return append
// }
//
// // GetFile :
// /*
//  * Function getFile :
//  * get the provided file from the online HTTP server and save to folder
//  */
// func GetFile(currentURL string, fileBaseURL string, fileLocation string, isByteRangeMPD bool, startRange int, endRange int, segmentNumber int, segmentDuration int, addSegDuration bool, quicBool bool, debugFile string, debugLog bool, useTestbedBool bool) (time.Duration, int) {
//
// 	// create the string where we want to save this file
// 	var createFile string
//
// 	// join the new file location to the base url
// 	urlHeaderString := JoinURL(currentURL, fileBaseURL, debugLog)
//
// 	logging.DebugPrint(debugFile, debugLog, "DEBUG: ", "get file from URL: "+urlHeaderString+"\n")
//
// 	if urlHeaderString == "" {
// 		fmt.Println("null urlHeader")
// 	}
//
// 	//request the URL with GET
// 	body, rtt := getURLBody(urlHeaderString, isByteRangeMPD, startRange, endRange, quicBool, debugFile, debugLog, useTestbedBool)
//
// 	// we only want the base file of the url (sometimes the segment media url has multiple folders)
// 	base := path.Base(fileBaseURL)
//
// 	// we need to create a file to save for the byte-range content
// 	if isByteRangeMPD {
// 		s := strings.Split(base, ".")
// 		base = s[0] + "_segment" + strconv.Itoa(segmentNumber) + ".m4s"
// 	}
//
// 	// create the new file location, or not
// 	if addSegDuration {
// 		createFile = fileLocation + "/" + strconv.Itoa(segmentDuration) + "sec_" + base
// 	} else {
// 		createFile = fileLocation + "/" + base
// 	}
//
// 	// save the file to the provided file location
// 	out, err := os.Create(createFile)
// 	if err != nil {
// 		fmt.Println("*** " + createFile + " cannot be downloaded ***")
// 		// stop the app
// 		utils.StopApp()
// 	}
// 	defer out.Close()
//
// 	// Write the body to file
// 	_, err = io.Copy(out, body)
// 	if err != nil {
// 		fmt.Println("*** " + createFile + " cannot be saved ***")
// 		// stop the app
// 		utils.StopApp()
// 	}
//
// 	fi, err := os.Stat(createFile)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	size := strconv.FormatInt(fi.Size(), 10)
// 	segSize, err := strconv.Atoi(size)
// 	if err != nil {
// 		logging.DebugPrint(debugFile, debugLog, "Error : ", "Cannot convert the size to an int when getting a file")
// 		utils.StopApp()
// 	}
//
// 	// close the body connection
// 	body.Close()
//
// 	return rtt, segSize
// }
//
// // GetFileProgressively :
// /*
//  * get the provided file from the online HTTP server and save to folder
//  * get a 1-second piece of each file
//  */
// func GetFileProgressively(currentURL string, fileBaseURL string, fileLocation string, isByteRangeMPD bool, startRange int, endRange int, segmentNumber int, segmentDuration int, addSegDuration bool, debugLog bool) (time.Duration, int) {
//
// 	// create the string where we want to save this file
// 	var createFile string
//
// 	// join the new file location to the base url
// 	urlHeaderString := JoinURL(currentURL, fileBaseURL, debugLog)
// 	logging.DebugPrint(glob.DebugFile, debugLog, "DEBUG: ", "get file from URL: "+urlHeaderString+"\n")
//
// 	if urlHeaderString == "" {
// 		fmt.Println("null urlHeader")
// 	}
//
// 	// we only want the base file of the url (sometimes the segment media url has multiple folders)
// 	base := path.Base(fileBaseURL)
//
// 	// we need to create a file to save for the byte-range content
// 	if isByteRangeMPD {
// 		s := strings.Split(base, ".")
// 		base = s[0] + "_segment" + strconv.Itoa(segmentNumber) + ".m4s"
// 	}
//
// 	// create the new file location, or not
// 	if addSegDuration {
// 		createFile = fileLocation + "/" + strconv.Itoa(segmentDuration) + "sec_" + base
// 	} else {
// 		createFile = fileLocation + "/" + base
// 	}
//
// 	// save the file to the provided file location
// 	out, err := os.Create(createFile)
// 	if err != nil {
// 		fmt.Println("*** " + createFile + " cannot be downloaded ***")
// 		// stop the app
// 		utils.StopApp()
// 	}
// 	defer out.Close()
//
// 	//request the URL with GET
// 	rtt := getURLProgressively(urlHeaderString, isByteRangeMPD, startRange, endRange, createFile)
//
// 	fi, err := os.Stat(createFile)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
//
// 	size := strconv.FormatInt(fi.Size(), 10)
// 	segSize, err := strconv.Atoi(size)
// 	if err != nil {
// 		logging.DebugPrint(glob.DebugFile, debugLog, "Error : ", "Cannot convert the size to an int when getting a file")
// 		utils.StopApp()
// 	}
//
// 	return rtt, segSize
// }
// */
