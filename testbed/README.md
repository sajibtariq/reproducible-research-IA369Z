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




