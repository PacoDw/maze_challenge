#!/bin/bash

create_topic() {
    TOPIC_NAME="test-topic"
    PARTITIONS=1
    REPLICATION_FACTOR=1
    LOCALHOST=2181

    for arg in $@
    do
        case "$arg" in
            --topic*|-t*)
            if [ "$arg" != *=* ]; then shift; fi
            TOPIC_NAME="${arg#*=}"
            ;;
            --partitions*|-p*)
            if [ "$arg" != *=* ]; then shift; fi
            PARTITIONS="${arg#*=}"
            ;;
            --replication_factor*|-rf*)
            if [ "$arg" != *=* ]; then shift; fi
            REPLICATION_FACTOR="${arg#*=}"
            ;;
            --localhost*|-lc*)
            if [ "$arg" != *=* ]; then shift; fi
            LOCALHOST="${arg#*=}"
            ;;
            --help|-h)
             >&2 printf "Parameters meaning and default values:
            -tn   --topic_name         By default is test-topic
            -p    --partitions         By default is 1
            -rf   --replication_factor By defaylt is 1
            -lc   --localhost          By default is 32181
            -h    --help \n" 
            echo ""
            return
            ;;
            *)
            >&2 printf "Error: Invalid argument ${arg}\nUse the option -h or --help to display the available commands"
            echo ""
            return
            ;;
        esac
    done

    echo "--topic ${TOPIC_NAME} --partitions ${PARTITIONS} --replication-factor ${REPLICATION_FACTOR} --if-not-exists --zookeeper localhost:${LOCALHOST}"

    return
}
