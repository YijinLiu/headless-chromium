#!/bin/bash

RED='\033[0;31m'
NC='\033[0m'
BASEDIR=$(readlink -f $(dirname $0))

# Parse command line.
OPTS=`getopt -n "install_native.sh" -o d -l debug -- "$@"`
rc=$?
if [ $rc != 0 ] ; then
    echo "install_native.sh [--debug]"
    exit 1
fi
eval set -- "$OPTS"

debug=
while true; do
    case "$1" in
        -d | --debug )           debug=1 ; shift ;;
        -- )                     shift; break ;;
        * )                      echo -e "${RED}Invalid option: $1${NC}" >&2 ; exit 1 ;;
    esac
done


install_tools() {
    sudo apt update && sudo apt install -y bzip2 curl file git patch python tar
    rc=$?
    if [ $rc != 0 ]; then
        echo -e "${RED}Failed to install tools${NC}"
        return 1;
    fi
}


install_headless_chromium() {
    if [ ! -d "chromium" ]; then
        mkdir -p chromium
    fi
    cd chromium
    version=52.0.2743.116

    if [ ! -d "depot_tools" ]; then
        git clone --depth=1 https://chromium.googlesource.com/chromium/tools/depot_tools.git
        rc=$?
        if [ $rc != 0 ]; then
            echo -e "${RED}Failed to download depot_tools${NC}"
            return 1
        fi
    fi
    export PATH=$PATH:`pwd`/depot_tools

    if [ ! -d "src" ]; then
	    echo "Downloading Chromium source code ..."
        fetch chromium --nosvn=True
        rc=$?
        if [ $rc != 0 ]; then
            echo -e "${RED}Failed to fetch chromium source code${NC}"
            return 1
        fi

        cd src
        git fetch --tags
        git checkout -b ${version} ${version} &&
        gclient sync --with_branch_heads
        rc=$?
        if [ $rc != 0 ]; then
            echo -e "${RED}Failed to checkout branch $BRANCH.${NC}"
            return 1
        fi

        patch -p1 < ${BASEDIR}/chromium.patch
        rc=$?
        if [ $rc != 0 ]; then
            echo -e "${RED}Failed to patch chromium${NC}"
            return 1
        fi
    else
        cd src
    fi

    sudo mkdir -p /usr/local/share/fonts &&
    ./build/install-build-deps.sh --chromeos-fonts --nacl --no-arm --no-syms --no-prompt
    rc=$?
    if [ $rc != 0 ]; then
        echo -e "${RED}Failed to install chromium dependencies${NC}"
        return 1
    fi

    args="import(\"//build/args/headless.gn\")\n
is_debug = false\n
remove_webcore_debug_symbols = true\n
symbol_level = 0"
    if [ -n "$debug" ]; then
        args="${args}\nis_component_build = true"
    fi

    gclient runhooks -force && mkdir -p out/Default &&
    echo -e $args > out/Default/args.gn &&
    gn gen out/Default && ninja -C out/Default headless_shell
    rc=$?
    if [ $rc != 0 ]; then
        echo -e "${RED}Failed to build chromium headless${NC}"
        return 1
    fi

    sudo mkdir -p /usr/local/lib &&
    sudo install -C out/Default/obj/headless/libheadless_lib.a /usr/local/lib/libheadless_chromium.a
    rc=$?
    if [ $rc != 0 ]; then
        echo -e "${RED}Failed to install libheadless_chromium.a${NC}"
        return 1
    fi

    if [ -n "$debug" ]; then
        sudo install -C out/Default/lib*.so /usr/local/lib/
        rc=$?
        if [ $rc != 0 ]; then
            echo -e "${RED}Failed to install component libs${NC}"
            return 1
        fi
    fi

    dest=/usr/local/include/headless_chromium
    folders="base base/debug base/memory base/numerics base/strings base/synchronization \
             base/threading base/time build headless/public headless/public/internal \
             headless/public/util net/base net/url_request ui/gfx ui/gfx/geometry url \
             url/third_party/mozilla"
    for folder in ${folders}
    do
        sudo mkdir -p ${dest}/${folder} &&
        sudo cp -f ${folder}/*.h ${dest}/${folder}/
        rc=$?
        if [ $rc != 0 ]; then
            echo -e "${RED}Failed to install ${folder}${NC}"
            return 1
        fi
    done
    sudo mkdir -p ${dest}/headless/public/domains &&
    sudo cp -f out/Default/gen/headless/public/domains/*.h ${dest}/headless/public/domains/
    rc=$?
    if [ $rc != 0 ]; then
        echo -e "${RED}Failed to install headless/public/domains${NC}"
        return 1
    fi
}

install_tools &&
install_headless_chromium
