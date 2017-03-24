#!/bin/bash

# This script sets up the environment to run headless-chromium using Golang.
#
# Usage:
#   1. ./setup.sh
#      Install hc_server.
#   2. ./setup.sh --install_dir=$DIR
#      Change the installation dir.
#   3. ./setup.sh --protocol_langs=golang
#      Generate library for the specified language. Only golang is supported now.

RED='\033[0;31m'
NC='\033[0m'
BASEDIR=$(readlink -f $(dirname $0)/..)

# Parse command line.
OPTS=`getopt -n "setup.sh" -o dl:t: -l debug,install_dir:,protocol_langs: -- "$@"`
rc=$?
if [ $rc != 0 ] ; then
    echo "Options:
    -d | --debug            Whether headless Chromium was compiled using debug mode.
    -t | --install_dir      Install dir.
    -l | --protocol_langs   Languages (separated by comma) to generate protocol libraries.
    "
    exit 1
fi
eval set -- "$OPTS"

debug=
install_dir=/usr/local/headless_chromium
protocol_langs=
while true; do
    case "$1" in
        -d | --debug )           debug=1 ; shift ;;
        -t | --install_dir)      install_dir="$2"; shift 2 ;;
        -l | --protocol_langs)   protocol_langs="$2"; shift 2 ;;
        -- )                     shift; break ;;
        * )                      echo -e "${RED}Invalid option: $1${NC}" >&2 ; exit 1 ;;
    esac
done

dir=$(pwd)

# Install hc_server
cd ${BASEDIR}/cc/hc_server
if [ -n "$debug" ] ; then
    DEBUG_CHROMIUM=1 make 
else
    make
fi
rc=$?
if [ $rc != 0 ] ; then
    echo -e "${RED}Failed to build hc_server.${NC}"
    exit 1
fi
sudo install hc_server $install_dir/bin
cd $dir

# Generate protocol libraries.
if [ -n "$protocol_langs" ] ; then
    GOPATH=$(readlink -f ${BASEDIR}/../../../..)
    cd ${GOPATH}
    GOPATH=${GOPATH} go install github.com/yijinliu/headless-chromium/go/protocol_parser &&
    ./bin/protocol_parser --output-langs=${protocol_langs} ${install_dir}/protocol/*.json
    rc=$?
    if [ $rc != 0 ] ; then
        echo -e "${RED}Failed to generate protocol libraries.${NC}"
        exit 1
    fi
    cd $dir
fi
