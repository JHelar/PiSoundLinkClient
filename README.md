# PiSoundLinkClient

#!/bin/bash
kaka=`youtube-dl -g $1`
SAVEIFS=$IFS
IFS=$'\n'
kaka=($kaka)
IFS=$SAVEIFS
omxplayer -o local ${kaka[1]}
