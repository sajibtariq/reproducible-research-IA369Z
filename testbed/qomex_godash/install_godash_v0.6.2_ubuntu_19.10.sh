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

# version 1.0.0 - 14-01-2020

# updated for ubuntu 18.04 to 19.10
sudo apt install net-tools build-essential git python3-pip unzip -y
sudo pip3 install pandas
sudo pip3 install numpy
sudo pip3 install matplotlib

#  download go version 1.13.5
wget https://dl.google.com/go/go1.14.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.14.linux-amd64.tar.gz
rm go1.14.linux-amd64.tar.gz

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

# update source for this terminal
. ~/.bashrc


# note: log out and log in again may be needed for terminal to link to the go executible...
