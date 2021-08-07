#! /usr/bin/env bash
CURDIR=$(cd $(dirname $0); pwd)

if [ "X$1" != "X" ]; then
    RUNTIME_ROOT=$1
else
    RUNTIME_ROOT=${CURDIR}
fi

SERVER='give_me_awesome'

#export CONF_DIR="$CURDIR/conf"
#export LOG_DIR="$RUNTIME_ROOT/log"
#exec "$CURDIR/bin/give_me_awesome"
#nohup ./${CURDIR}/bin/give_me_awesome 2>&1 &>nohup.out &

getPid(){
    pid=`ps aux | grep ${SERVER} | grep -v 'grep' | awk '{print $2}'`
    echo ${pid}
}

start(){
    pid=$(getPid)
    if [[ ${pid} ]];then
        echo "SERVER in running... pid:$pid"&&exit 0
    fi
    pro=${CURDIR}/${SERVER}
    if [[ ! -x ${pro} ]];then
        echo "${pro} can not exec!!!"
        return
    fi
    nohup ${CURDIR}/${SERVER} 2>&1 &>nohup.out &
    sleep 1
    status
}

stop(){
    pid=$(getPid)
    if [[ ! ${pid} ]];then
        echo "$server is not running..."
        return
    fi
    kill ${pid}
}

status(){
    pid=$(getPid)
    if [[ ${pid} ]];then
        echo "$server is running..."
    else
        echo "$server is not running..."
    fi
}

restart(){
    stop
    sleep 1
    start
}

case $1 in
    start)
        start
        echo "start Done!"
        ;;
    stop)
        stop
        echo "stop Done!"
        ;;
    restart)
        restart
        echo "restart Done!"
        ;;
    status)
        status
        ;;
    pid)
        getPid
        ;;
esac
