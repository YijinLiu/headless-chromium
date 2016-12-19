#!/bin/bash

# This script installs headless chromium headers and libraries in your system.
# NOTE: Only tested on Ubuntu 16.04 / x86_64.

RED='\033[0;31m'
NC='\033[0m'
BASEDIR=$(readlink -f $(dirname $0))

# Parse command line.
OPTS=`getopt -n "install_chromium.sh" -o dt:u: -l debug,install_dir:,upgrade_to: -- "$@"`
rc=$?
if [ $rc != 0 ] ; then
    echo "install_chromium.sh [--debug]"
    exit 1
fi
eval set -- "$OPTS"

debug=
install_dir=/usr/local/headless_chromium
upgrade_to=
while true; do
    case "$1" in
        -d | --debug )           debug=1 ; shift ;;
        -t | --install_dir)      install_dir="$2"; shift 2 ;;
        -u | --upgrade_to)       upgrade_to="$2"; shift 2 ;;
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
        git fetch --tags
        # To list all tags:
        #   git tag -l
        sync_to_version "55.0.2883.91"
        rc=$?
        if [ $rc != 0 ]; then
            return 1
        fi
    else
        cd src

        if [ -n "$upgrade_to" ]; then
            git rebase-update &&
            gclient sync &&
            sync_to_version "${upgrade_to}"
            rc=$?
            if [ $rc != 0 ]; then
                return 1
            fi
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
             base/numerics base/profiler base/strings base/synchronization base/task_scheduler \
             base/threading base/time build content/common content/public/common \
             headless/public headless/public/internal headless/public/util \
             mojo/public/c/system mojo/public/cpp/bindings mojo/public/cpp/bindings/lib \
             mojo/public/cpp/system net/base net/url_request testing/gtest/include/gtest \
             ui/base ui/base/touch ui/gfx ui/gfx/geometry url url/third_party/mozilla"
    for folder in ${folders}
    do
        sudo mkdir -p ${install_dir}/include/${folder} &&
        sudo cp -f ${folder}/*.h ${install_dir}/include/${folder}/
        rc=$?
        if [ $rc != 0 ]; then
            echo -e "${RED}Failed to install ${folder}${NC}"
            return 1
        fi
    done
    sudo mkdir -p ${install_dir}/include/headless/public/domains &&
    sudo cp -f out/Default/gen/headless/public/domains/*.h \
        ${install_dir}/include/headless/public/domains/
    rc=$?
    if [ $rc != 0 ]; then
        echo -e "${RED}Failed to install headless/public/domains${NC}"
        return 1
    fi

    # Install protocols.
    sudo mkdir -p ${install_dir}/protocol &&
    sudo install third_party/WebKit/Source/core/inspector/browser_protocol.json \
        v8/src/inspector/js_protocol.json ${install_dir}/protocol

    cd ../..
}

install_tools &&
install_headless_chromium
