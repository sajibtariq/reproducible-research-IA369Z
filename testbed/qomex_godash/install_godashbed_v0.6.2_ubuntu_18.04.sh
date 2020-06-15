#!/bin/bash

#
#    This program is free software; you can redistribute it and/or
#    modify it under the terms of the GNU General Public License
#    as published by the Free Software Foundation; either version 2
#    of the License, or (at your option) any later version.
#
#    This program is distributed in the hope that it will be useful,
#    but WITHOUT ANY WARRANTY; without even the implied warranty of
#    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
#    GNU General Public License for more details.
#
#    You should have received a copy of the GNU General Public License
#    along with this program; if not, write to the Free Software
#    Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA
#    02110-1301, USA.
#

# version 1.0.0 - 01-02-2020

# updated for ubuntu 18.04
# due to mininet issues in ubuntu 19.10 - only 18.04 is currently approved for installation
sudo apt install net-tools build-essential git python3-pip unzip -y
sudo pip3 install pandas
sudo pip3 install numpy
sudo pip3 install matplotlib

#  download go version 1.13.5
wget https://dl.google.com/go/go1.13.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.13.5.linux-amd64.tar.gz
rm go1.13.5.linux-amd64.tar.gz

#if you want to install golang on OSX, use this:
# curl -o golang.pkg https://dl.google.com/go/go1.13.5.darwin-amd64.pkg
# sudo open golang.pkg
# and follow the install steps

SCRIPTPATH="$( cd "$(dirname "$0")" ; pwd -P )"
echo $SCRIPTPATH

# download godash version 0.6.2
wget https://www.dropbox.com/s/ziwpz2eqjrnqoys/goDash_v0.6.2.zip
unzip goDash_v0.6.2.zip
rm -rf ./__MACOSX
cd goDash

# install itu 1203
git clone http://github.com/itu-p1203/itu-p1203.git
cd itu-p1203/
pip3 install .
python3 -m itu_p1203 examples/mode0.json

cd ..

# add path to go in bashrc
echo export GOROOT=/usr/local/go >> $HOME/.bashrc
echo export GOPATH=$SCRIPTPATH/goDash/DashApp >> $HOME/.bashrc
echo export PATH='''$PATH:$GOROOT/bin:$GOPATH/bin''' >> $HOME/.bashrc
source $HOME/.bashrc

export GOROOT=/usr/local/go
export GOPATH=$SCRIPTPATH/goDash/DashApp
export PATH=$PATH:$GOROOT/bin:$GOPATH/bin

cd DashApp/src/goDASH
go get -u gonum.org/v1/gonum/...
go get github.com/cavaliercoder/grab
go get github.com/lucas-clemente/quic-go/http3
go get github.com/marten-seemann/qpack
go get github.com/onsi/ginkgo
go get github.com/onsi/gomega
go get github.com/marten-seemann/qtls

# build player
go build

cd $SCRIPTPATH
# clone mininet repo
git clone git://github.com/mininet/mininet

# cd into folder
cd mininet/util

# install mininet with all options
./install.sh -a

# we need to setup server for quic protocol
# these include building a quic server with certificates
cd $SCRIPTPATH

# download goDASHbed
wget https://www.dropbox.com/s/0bedfwfoszvnbdu/goDASHbed.zip
unzip goDASHbed.zip
rm -rf ./__MACOSX

# copy cert.pem and key.pem
cp goDash/DashApp/src/goDASH/http/certs/cert.pem goDash/DashApp/src/github.com/lucas-clemente/quic-go/internal/testdata/cert.pem
cp goDash/DashApp/src/goDASH/http/certs/key.pem goDash/DashApp/src/github.com/lucas-clemente/quic-go/internal/testdata/priv.key
cd goDash/DashApp/src/github.com/lucas-clemente/quic-go/example
go build
cp example $SCRIPTPATH/goDASHbed/example

cd $SCRIPTPATH

# link godashbed.org to IP address in ubuntu
echo "10.0.0.1   www.godashbed.org" | sudo tee -a /etc/hosts

# install apache, so we can host content in /var/www
sudo apt install apache2 -y

# add the voip generator for goDAShbed
wget https://www.dropbox.com/s/8b4ymyxrt78ggfq/D-ITG-2.8.1-r1023.zip
unzip D-ITG-2.8.1-r1023.zip
cd D-ITG-2.8.1-r1023/src
make
sudo make install

# update source for this terminal
. ~/.bashrc


# note: log out and log in again may be needed for terminal to link to both P.1203 and go executibles...
