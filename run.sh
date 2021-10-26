#!/bin/bash

PRO_NAME="godaemon"
PORT="10600"
HOST="0.0.0.0"
BASEDIR=$(dirname $(realpath "$0"))

function start(){
    cd $BASEDIR
    if [ ! -d "log"  ]
    then
    	mkdir log
    fi

    if [ "`timeout 2 bash -c "</dev/tcp/${HOST}/${PORT}" &> /dev/null; echo $?`" -ne 0 ]     
    then
         chmod +x ${BASEDIR}/${PRO_NAME}
         nohup ./${PRO_NAME} &> /dev/null &
         sleep 1
    fi
    PID=`ps -ef | grep -v $0| grep "godaemon" | awk '{print $2}'`
    if ps -p $PID > /dev/null
    then
        echo -e "\033[42;37m ${PRO_NAME} is Running... \033[0m"
    fi
}

function stop(){
    PID=`ps -ef | grep -v $0| grep "godaemon" | awk '{print $2}'`
    kill $PID
    sleep 2
    if [ "`timeout 2 bash -c "</dev/tcp/${HOST}/${PORT}" &> /dev/null; echo $?`" -ne 0 ]
    then
        echo -e "\033[42;37m ${PRO_NAME} is exit... \033[0m"
    fi
}


case "$1" in 
    start)   start ;;
    stop)    stop ;;
    restart) stop; start ;;
    *) echo "usage: $0 start|stop|restart" >&2
       exit 1
       ;;
esac
