#!/bin/bash

# This script installs headless Chromium headers and libraries in your system.
# Tested on Ubuntu 16.04 / x86_64.
# NOTE: It takes very long time to download and compile the Chromium source code. Be patient.
#
# Usage:
#   1. ./install_chromium.sh
#      Download source code to ./chromium, compile and install to /usr/local/headless_chromium
#   2. ./install_chromium.sh --install_dir=$DIR
#      Change the installation dir.
#   3. ./install_chromium.sh --debug
#      Compile and install debug version Chromium.
#   4. ./install_chromium.sh --version=$VERSION
#      Sync Chromium to the specified version.
#      NOTE: Chromium source code is updated very frequently. There is no guarantee that this works.
#
#   To find the available Chromium versions:
#      cd chromium/src
#      git fetch --tags
#      git tag -l

RED='\033[0;31m'
NC='\033[0m'
BASEDIR=$(readlink -f $(dirname $0))

# Parse command line.
OPTS=`getopt -n "install_chromium.sh" -o dt:v: -l debug,install_dir:,version: -- "$@"`
rc=$?
if [ $rc != 0 ] ; then
    echo "install_chromium.sh [--debug]"
    exit 1
fi
eval set -- "$OPTS"

debug=
install_dir=/usr/local/headless_chromium
version="57.0.2987.110"
while true; do
    case "$1" in
        -d | --debug )           debug=1 ; shift ;;
        -t | --install_dir)      install_dir="$2"; shift 2 ;;
        -v | --version)          version="$2"; shift 2 ;;
        -- )                     shift; break ;;
        * )                      echo -e "${RED}Invalid option: $1${NC}" >&2 ; exit 1 ;;
    esac
done


install_tools() {
    sudo apt update && sudo apt install -y bzip2 clang curl file git patch python tar
    rc=$?
    if [ $rc != 0 ]; then
        echo -e "${RED}Failed to install tools${NC}"
        return 1
    fi
}


sync_to_version() {
    version="$1"

    git show-branch --list | grep "[$version]"
    rc=$?
    if [ $rc != 0 ]; then
        git checkout -b ${version} ${version}
    else
        git checkout ${version}
    fi
    rc=$?
    if [ $rc != 0 ]; then
        echo -e "${RED}Failed to checkout branch ${version}.${NC}"
        return 1
    fi

    gclient sync --with_branch_heads
    rc=$?
    if [ $rc != 0 ]; then
        echo -e "${RED}Failed to sync branch ${version}.${NC}"
        return 1
    fi

    patch --forward -p1 < ${BASEDIR}/chromium.patch
    rc=$?
    if [ $rc -gt 1 ]; then
        echo -e "${RED}Failed to patch chromium${NC}"
        return 1
    fi
}

install_headless_chromium() {
    if [ ! -d "chromium" ]; then
        mkdir -p chromium
    fi
    cd chromium

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
        sync_to_version "${version}"
        rc=$?
        if [ $rc != 0 ]; then
            return 1
        fi
    else
        cd src
        cur_version=`git status -uno | grep "On branch" | sed 's/On branch //'`
        if [ -n "$version" ] && [ "$version" != "$cur_version" ] ; then
            git checkout . &&
            git checkout master &&
            gclient sync &&
            sync_to_version "${version}"
            rc=$?
            if [ $rc != 0 ]; then
                return 1
            fi
            rm -rf out/Default/gen
        fi
    fi

    sudo mkdir -p /usr/local/share/fonts &&
    ./build/install-build-deps.sh --chromeos-fonts --nacl --no-arm --no-syms --no-prompt
    rc=$?
    if [ $rc != 0 ]; then
        echo -e "${RED}Failed to install chromium dependencies${NC}"
        return 1
    fi

    # Set is_clang to false to compile with gcc.
    # NOTE: There is a memory leak bug or infinite loop if compiled with gcc 5.4.0.
    # To compile use system clang, change two lines in build/toolchain/gcc_toolchain.gni
    #   cc = "$prefix/clang"     => cc = "clang"
    #   cxx = "$prefix/clang++"  => cxx = "clang++"
    args="import(\"//build/args/headless.gn\")\n
is_debug = false\n
remove_webcore_debug_symbols = true\n
symbol_level = 0\n
use_sysroot = false\n
is_clang = true\n
clang_use_chrome_plugins = false\n
v8_use_external_startup_data = false\n
icu_use_data_file = false\n"
    if [ -n "$debug" ]; then
        args="${args}\nis_component_build = true"
    fi

    gclient runhooks --force && mkdir -p out/Default &&
    echo -e $args > out/Default/args.gn &&
    gn gen out/Default && ninja -C out/Default headless_shell
    rc=$?
    if [ $rc != 0 ]; then
        echo -e "${RED}Failed to build chromium headless${NC}"
        return 1
    fi

    # Install binaries.
    sudo rm -rf ${install_dir} &&
    sudo mkdir -p ${install_dir}/bin &&
    sudo install out/Default/headless_shell out/Default/headless_lib.pak out/Default/libosmesa.so \
        ${install_dir}/bin/

    # Install libraries.
    sudo mkdir -p ${install_dir}/lib &&
    sudo install out/Default/obj/headless/libheadless_lib.a ${install_dir}/lib/
    rc=$?
    if [ $rc != 0 ]; then
        echo -e "${RED}Failed to install libheadless_lib.a${NC}"
        return 1
    fi

    if [ -n "$debug" ]; then
        sudo install out/Default/lib*.so ${install_dir}/lib/
        rc=$?
        if [ $rc != 0 ]; then
            echo -e "${RED}Failed to install component libs${NC}"
            return 1
        fi
    fi

    # Install headers.
    folders="base base/containers base/debug base/files base/memory base/message_loop \
             base/numerics base/process/ base/profiler base/strings base/synchronization \
             base/task_scheduler base/threading base/time build content/common \
             content/public/common headless/public headless/public/internal headless/public/util \
             mojo/public/c/system mojo/public/cpp/bindings mojo/public/cpp/bindings/lib \
             mojo/public/cpp/system net/base net/url_request testing/gtest/include/gtest \
             ui/base ui/base/touch ui/gfx ui/gfx/geometry url url/third_party/mozilla"
    for folder in ${folders}
    do
        sudo mkdir -p ${install_dir}/include/${folder} &&
        sudo cp -f ${folder}/*.h ${install_dir}/include/${folder}/
        rc=$?
        if [ $rc != 0 ]; then
            echo -e "${RED}Failed to install ${folder}.${NC}"
            return 1
        fi
    done
    folders="headless/public/domains headless/public/devtools/domains \
             headless/public/devtools/internal"
    for folder in ${folders}
    do
        sudo mkdir -p ${install_dir}/include/${folder} &&
        sudo cp -f out/Default/gen/${folder}/*.h ${install_dir}/include/${folder}/
        rc=$?
        if [ $rc != 0 ]; then
            echo -e "${RED}Failed to install ${folder}.${NC}"
            return 1
        fi
    done

    # Install protocols.
    sudo mkdir -p ${install_dir}/protocol &&
    sudo install third_party/WebKit/Source/core/inspector/browser_protocol.json \
        v8/src/inspector/js_protocol.json ${install_dir}/protocol

    cd ../..
}

install_tools &&
install_headless_chromium
