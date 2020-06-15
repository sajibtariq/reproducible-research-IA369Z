from mininet.log import setLogLevel, info
from mininet.node import Controller, RemoteController
from mininet.link import TCLink
from mn_wifi.cli import CLI
from mn_wifi.net import Mininet_wifi
from multiprocessing import Process
import os
import time
import  random
import sys


station = []

def topology():

    mod = str(sys.argv[1]) # network type

    nett = str(sys.argv[2]) # mobility

    doc = str(sys.argv[3]) # operator

    num = int(sys.argv[4]) # traces

    host = int(sys.argv[5]) # num of host

    algo = str(sys.argv[6]) # name of adaptation algorithm


###################################
    "Create a network."
    net = Mininet_wifi()

    info("*** Creating nodes\n")


    for i in range(host):
        m='sta%s' % (i+1)
        j=i+1
        station.insert(i, net.addStation(m, ip='10.0.0.%s'%(j)))


    ap1 = net.addAccessPoint('ap1', ssid="simpletopo", mode="g", channel="5")

    server = net.addHost('server',ip='10.0.0.150')
    s3 = net.addSwitch('s3')
    c0 = net.addController('c0')

    info("*** Configuring wifi nodes\n")
    net.configureWifiNodes()

    info("*** Associating Stations\n")
    for i in range(host):
        m='sta%s' % (i+1)
        net.addLink(m, ap1)
    net.addLink(s3, ap1, 1, 10)  # initial link parameter default according to mininet
    net.addLink(server, s3)

    info("*** Starting network\n")
    net.build()
    c0.start()
    ap1.start([c0])
    s3.start([c0])


##############################
    time.sleep(10)
    info("*** Running CLI\n")
    #CLI(net)
    #print net.get('ap1')
    #print net['s3']

    os.system('cd ~/reproducible-research-IA369Z/testbed/qomex_godash/goDash/DashApp/src/files && sudo rm -R *')


    subfolder = 'mode_'+(mod)+ '_net_' + str(nett) +'_doc_' + str(doc) + '_num_' + str(num)+'_host_'+ str(host)+ '_algo_' + str(algo)
    
    global  path
    path='~/reproducible-research-IA369Z/testbed/experiment'
    os.system('mkdir -p '+ path+'/'+subfolder)
    #os.system('mkdir -p ~/reproducible-research-IA369Z/data/testbed/experiment/'+ subfolder)

    st=[]
    for i in range(host):
        m1='sta%s'%(i+1)
        m2=net[m1]
        st.insert(i, m2)

    switch=net['s3']
    server=net['server']
    ap= net['ap1']


    return st, switch, server, ap, host , algo, nett, doc, num, mod



# caddy server on
def server(sr):
    print sr
    print sr.cmd('cd ~/reproducible-research-IA369Z/testbed/caddy && ./caddy')

# pcap collect by tcpdump  from Ap-interface
def cap2(mod, host, algo, nett, doc, num):

    print os.system('sudo tcpdump -i ap1-eth10 -U -w ' +path+'/mode_' + str(mod) + '_net_' + str(nett) +'_doc_' + str(doc) + '_num_' + str(num)+'_host_'+ str(host)+ '_algo_' + str(algo)+'/mode_' + str(mod) + '_net_' + str(nett) +'_doc_' + str(doc) + '_num_' + str(num)+'_host_'+ str(host)+ '_algo_' + str(algo)+'_ap.pcap')

#tc script
def netcon(mod, nett, doc, num, host):

       os.system('sudo ./ta1.sh %s %s %s %d %d'%(mod, nett, doc, num, host))

# run godash player
def test1(mod, client, host, algo, nett, doc, num):

    #if algo =='conv':
        print client
        print client.cmd('cd ~/reproducible-research-IA369Z/testbed/qomex_godash/goDash/DashApp/src/goDASH && ./goDASH -config ../config/configure.json > ' +path+'/mode_' + str(mod) +'_net_' + str(nett) +'_doc_' + str(doc) + '_num_' + str(num)+'_host_'+ str(host)+ '_algo_' + str(algo)+'/mode_' + str(mod) +'_net_' + str(nett) +'_doc_' + str(doc) + '_num_' + str(num) + '_client_' + str(client)+ '_host_' + str(host)+ '_algo_' + str(algo)+'.txt && echo done_' + str(client))



# stop caddy server

def tstop(host, algo):

     if (algo == 'arbiter' or algo == 'bba-2'):
        tt=1100*host
        time.sleep(tt)
        os.system ('sudo pkill -9 caddy')

     else:
    	tt=1100*host
    	time.sleep(tt)
    	os.system('sudo pkill -9 caddy')

# stop tcpdump

def tsstop(host, algo):

    if (algo== 'arbiter' or algo == 'bba-2'):
    	tt=1100*host
    	time.sleep(tt)
    	os.system('sudo pkill -9 tcpdump')
    else:
    	tt=1100*host
    	time.sleep(tt)
    	os.system('sudo pkill -9 tcpdump')



if __name__ == '__main__':
    setLogLevel( 'info' )

    station, switch, ser, ap, host, algo, nett, doc, num, mod =  topology()

    a=True;b=False;c=False;d=False; e=False; f=False;g=False; h=False;

    if a:
       y=Process(target=server, args=(ser,))
       y.start()
       b=True

    if b:
       n=Process(target=cap2, args=(mod,host,algo,nett,doc,num))
       n.start()
       c=True

    if c:
       nn=Process(target=netcon, args=(mod,nett,doc,num,host))
       nn.start()
       d=True
    if d:
       #print 'dashc'
       for k in range(host):
           print station[k]
           q = Process(target=test1, args=(mod,station[k],host,algo,nett,doc,num))
           q.start()
           q.join
           e=True

    if e:
       t = Process(target=tstop, args=(host,algo))
       t.start()
       f=True

    if f:
       tt = Process(target=tsstop, args=(host,algo))
       tt.start()
