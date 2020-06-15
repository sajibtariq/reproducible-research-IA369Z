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
	glob "goDASH/global"
	"goDASH/logging"
)

// CreateQoE : get the P1203 and clae QoE values
func CreateQoE(log *map[int]logging.SegPrintLogInformation, debugLog bool, initBuffer int, maxRepRate int) {

	// *log does not support indexing :(
	logMap := *log

	// the P1203 standard only works for H264 (encoder) and up to resolutions of 1920x1080
	// so make sure the received segments are compliant
	logging.DebugPrint(glob.DebugFile, debugLog, "\nDEBUG: ", "checking for P1203 compatibility")
	for a := 1; a <= len(*log); a++ {

		// get the encoder and resolution
		width := logMap[a].RepWidth
		height := logMap[a].RepHeight
		codec := logMap[a].RepCodec

		if codec != glob.RepRateCodecAVC || width > glob.P1203maxWidth || height > glob.P1203maxHeight {
			logging.DebugPrint(glob.DebugFile, debugLog, "\nDEBUG: ", "Downloaded segments are not P1203 compliant")
			return
		}
	}

	// create channels, so the output is in the right order
	P1023Results := make(chan float64)
	claeResults := make(chan float64)
	duanmuResults := make(chan float64)
	yinResults := make(chan float64)
	yuResults := make(chan float64)

	// create the P1203 value
	go createP1203(log, P1023Results)

	// create the Claye value
	go getClaye(*log, claeResults, maxRepRate, false)

	// create the Duanmu value
	go getDuanmu(*log, duanmuResults, initBuffer, false)

	// create the Yin value
	go getYin(*log, yinResults, initBuffer, false)

	// create the Yu value
	go getYu(*log, yuResults, false)

	// create a local copy of the log and allocate the QoE values
	locallogMap := *log
	locallog := locallogMap[len(locallogMap)]
	// calculate the P1203, Claye, Duanmu, Yin and Yu values and
	// save to the last log as 3 decimal floats
	locallog.Yin = <-yinResults
	locallog.Yu = <-yuResults
	locallog.Duanmu = <-duanmuResults
	locallog.Clae = <-claeResults
	locallog.P1203 = <-P1023Results

	locallogMap[len(locallogMap)] = locallog
	*log = locallogMap
}
