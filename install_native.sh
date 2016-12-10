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
        version=55.0.2883.91
        git checkout -b ${version} ${version} &&
        gclient sync --with_branch_heads
        rc=$?
        if [ $rc != 0 ]; then
            echo -e "${RED}Failed to checkout branch $BRANCH.${NC}"
            return 1
        fi
        # To update source code manually:
        #   git rebase-update
        #   gclient sync
        # To list all tags:
        #   git tag -l

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

    # Set is_clang to false to compile with gcc.
    # NOTE: There is a memory leak bug or infinite loop if compiled with gcc 5.4.0.
    # To compile use system clang, change two lines in build/toolchain/gcc_toolchain.gni
    #   cc = "$prefix/clang"     => cc = "clang"
    #   cxx = "$prefix/clang++"  => cc = "clang++"
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

    sudo mkdir -p /usr/local/lib &&
    sudo install out/Default/obj/headless/libheadless_lib.a /usr/local/lib/libheadless_chromium.a
    rc=$?
    if [ $rc != 0 ]; then
        echo -e "${RED}Failed to install libheadless_chromium.a${NC}"
        return 1
    fi

    if [ -n "$debug" ]; then
        sudo install out/Default/lib*.so /usr/local/lib/
        rc=$?
        if [ $rc != 0 ]; then
            echo -e "${RED}Failed to install component libs${NC}"
            return 1
        fi
    fi

    dest=/usr/local/include/headless_chromium
    folders="base base/containers base/debug base/files base/memory base/message_loop \
             base/numerics base/profiler base/strings base/synchronization base/task_scheduler \
             base/threading base/time build content/common content/public/common \
             headless/public headless/public/internal headless/public/util \
             mojo/public/c/system mojo/public/cpp/bindings mojo/public/cpp/bindings/lib \
             mojo/public/cpp/system net/base net/url_request testing/gtest/include/gtest \
             ui/base ui/base/touch ui/gfx ui/gfx/geometry url url/third_party/mozilla"
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

    cd ../..
}

install_tools &&
install_headless_chromium
