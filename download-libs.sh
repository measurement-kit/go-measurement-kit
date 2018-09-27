#!/bin/bash
set -e
LIBS_VERSION="0.9.0-alpha.2-1"

TARGET_PLATFORM="all"
PKG_TOPDIR=$(cd $(dirname $0) && pwd -P)

download_libs()
{
    PLATFORM=$1
    DOWNLOAD_URL="https://github.com/measurement-kit/golang-prebuilt/releases/download/v${LIBS_VERSION}/libs_${PLATFORM}_${LIBS_VERSION}.tar.gz"

    echo "Downloading libs for $PLATFORM"

    cd $PKG_TOPDIR/libs
    echo "  downloading $DOWNLOAD_URL into $PKG_TOPDIR/libs"
    curl -LsO $DOWNLOAD_URL
    ARCHIVE=libs_${PLATFORM}_${LIBS_VERSION}.tar.gz

    if [ -d $PKG_TOPDIR/libs/$PLATFORM ];then
        rm -rf $PKG_TOPDIR/libs/$PLATFORM
    fi

    tar xzf $ARCHIVE
    rm $ARCHIVE
}

if [ "$1" != "" ];then
    TARGET_PLATFORM=$1
fi

mkdir -p $PKG_TOPDIR/libs
if [ "$TARGET_PLATFORM" == "all" ];then
    download_libs macos
    download_libs mingw
    download_libs linux
elif [ "$TARGET_PLATFORM" == "macos" ];then
    download_libs macos
elif [ "$TARGET_PLATFORM" == "linux" ];then
    download_libs linux
elif [ "$TARGET_PLATFORM" == "mingw" ];then
    download_libs mingw
else
    echo "Error: Unsupported platform $TARGET_PLATFORM"
    exit 1
fi
exit 0
