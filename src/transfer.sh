#/bin/bash

## Shell file for running with SSH

# first time running
#go build main.go # for å lage objektfil

#scp -r /home/student/Desktop/Sanntid/Elevator_project_gr_16/ student@129.241.187.153:gruppe16/ # kopierer over hele mappen til gitt ip-adresse
#scp -r /home/student/Desktop/Sanntid/Elevator_project_gr_16/ student@129.241.187.147:gruppe16/

#gnome-terminal --title "virtual_153: server" -x ssh student@129.241.187.153 & # for å kjøre på annen maskin
#gnome-terminal --title "virtual_155: server" -x ssh student@129.241.187.157 &

# second time
#scp -r /home/student/Desktop/Sanntid/Elevator_project_gr_16/main student@129.241.187.153:gruppe16/main # kopierer over objektfil til gitt ip-adresse
#scp -r /home/student/Desktop/Sanntid/Elevator_project_gr_16/main student@129.241.187.147:gruppe16/main
#scp -r /home/student/Desktop/Sanntid/Elevator_project_gr_16/main student@129.241.187.157:gruppe16/main
#scp -r /home/student/Desktop/Sanntid/Elevator_project_gr_16/main student@129.241.187.148:gruppe16/main




###################################################
# Vilde, Mia og Marie

# start running on this computer
go build main.go
scp -r /home/student/Desktop/VildeogMarie-master/Project student@129.241.187.147:gruppe79/

# open new terminal, connect to server
gnome-terminal --title "virtual_147: server" -x ssh student@129.241.187.147 &

scp -r /home/student/Desktop/VildeogMarie-master/Project/main student@129.241.187.147:gruppe79/main


#Package loss
#sudo iptables -D INPUT 1 # skru av
#sudo iptables -A INPUT -m statistic --mode random --probability 0.15 -j DROP # skru på





