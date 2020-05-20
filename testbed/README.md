# Testbed

## Requirements



Ubuntu 18.04 Operating Syatem

[Mininet-Wifi](https://github.com/intrig-unicamp/mininet-wifi)

[goDASH](https://github.com/uccmisl/goDASH)

[Caddy](https://caddyserver.com/)

[Tcpdump](https://www.tcpdump.org/)



## Installation
It might take longer. Please be patient!!!. It might also ask to put  yes/ y /accept etc during the installation. No worries whatever it ask just follow the instaruction.  
```bash
step 1: $ git clone https://github.com/sajibtariq/reproducible-research-IA369Z.git

step 2: $ cd reproducible-research-IA369Z

step 3: $ cd testbed

step 2: $ sudo chmood 777 build.sh

step 4: $ ./build.sh
```
After installation move caddy file from  /usr/local/bin/ directory to /reproducible-research-IA369Z/testbed/caddy/  directory
```bash
$ sudo mv  /usr/local/bin/caddy ../reproducible-research-IA369Z/testbed/caddy/
```
To download movie content inside caddy directory use ```bash dash_movie_content.sh script. ``` It might take longer. Please be patient!!!
```bash
$ cd caddy && sudo chmod 777 dash_movie_content.sh && ./dash_movie_content.sh
```
**Although goDASH player support several adaptive algorithms, for this class project we will use only conventional adaptive algorithm**

Go to ../reproducible-research-IA369Z/testbed/goDash/DashApp/src/ directory, open the config file and modify the file as follows-
```bash
{
        "adapt" : "conventional",
        "codec" : "h264",
        "debug" : "on",
        "initBuffer" : 2,
        "maxBuffer" : 60,
        "maxHeight" : 1080,
        "streamDuration" : 600,
        "logFile" : "log_file_2",
        "getHeaders" : "off",
        "terminalPrint" : "on",
        "printHeader" : "{\"Algorithm\":\"on\",\"Seg_Dur\":\"on\",\"Codec\":\"off\",\"Width\":\"on\",\"Height\":\"on\",\"FPS\":\"off\",\"Play_Pos\":\"on\",\"RTT\":\"on\",\"Seg_Repl\":\"off\",\"Protocol\":\"off\",\"P.1203\":\"on\",\"Clae\":\"off\",\"Duanmu\":\"off\",\"Yin\":\"off\",\"Yu\":\"off\"}",
        "useTestbed" : "off",
        "url" : "[http://10.0.0.150:2015/html/x264/bbb/DASH_Files/live/bbb_enc_x264_dash.mpd]",
        "QoE" : "on"
}
```
## Run Network Topology

Move to directory /reproducible-research-IA369Z/testbed/ and run the ```bash test_1.py``` script

```bash
$ sudo python3 test_1.py
```
 ```bash test_1.py``` script contains all the network utilization information as follows for network emulation and internally called ```bash test3_1.py``` main script to run the topology and generate the data
First Header  | Second Header
------------- | -------------
Content Cell  | Content Cell
Content Cell  | Content Cell
