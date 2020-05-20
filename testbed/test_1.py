import os
import sys
import subprocess

mode=['3g','4g','5g'] #network type '5g',

host=[1] # number of host
algo=['average'] # adaptation algorithm 'conv'

net3= ['metro','bus', 'train', 'ferry','car'] # mobility for 3g
net4=['bus', 'train', 'static','car','pedestrian'] # mobility for 4g
net5=['A_A_Static','D_Driving', 'D_Static'] #'A_A_Static','A_A_Driving','N_A_Driving',, ,,'D_Driving' mobility for 5g


doc3=['Am'] # number of operator (1)
doc4=['Am','Bm']   # number of operator (2)
doc5=['Bm']   # number of operator (1)

num = [1,2,3] # number of traces 

count = 1

for curr in range(count):
    for md in mode:
        if md == '5g':
           for i in net5:
               for j in doc5:
                   for k in num:
                       for l in host:
                           for m in algo:
                               clear = 'sudo mn -c'
                               test3 = 'sudo python test3_1.py '+ str(md)+ ' ' + str(i) + ' ' + str(j)+ ' ' + str(k) + ' ' + str(l)+ ' ' + str(m)
                               subprocess.run(clear.split(' '))
                               print(test3)
                               subprocess.run(test3.split(' '))
        elif md =='4g':
            for i in net4:
               for j in doc4:
                   for k in num:
                       for l in host:
                           for m in algo:
                               clear = 'sudo mn -c'
                               test3 = 'sudo python test3_1.py '+ str(md)+ ' ' + str(i) + ' ' + str(j)+ ' ' + str(k) + ' ' + str(l)+ ' ' + str(m)
                               subprocess.run(clear.split(' '))
                               print(test3)
                               subprocess.run(test3.split(' '))
        else:
            for i in net3:
               for j in doc3:
                   for k in num:
                       for l in host:
                           for m in algo:
                               clear = 'sudo mn -c'
                               test3 = 'sudo python test3_1.py '+ str(md)+ ' ' + str(i) + ' ' + str(j)+ ' ' + str(k) + ' ' + str(l)+ ' ' + str(m)
                               subprocess.run(clear.split(' '))
                               print(test3)
                               subprocess.run(test3.split(' '))
