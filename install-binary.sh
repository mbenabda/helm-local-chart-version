#!/bin/bash

# Shamelessly copied from https://github.com/technosophos/helm-template
#
# Helm Template Plugin
# Copyright (C) 2016, Matt Butcher
# 
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
# 
# The above copyright notice and this permission notice shall be included in
# all copies or substantial portions of the Software.
# 
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
# THE SOFTWARE.

PROJECT_BIN_NAME="local-chart-version"
PROJECT_NAME="helm-${PROJECT_BIN_NAME}"
PROJECT_GH="mbenabda/${PROJECT_NAME}"

: ${HELM_PLUGIN_PATH:="$(helm home)/plugins/${PROJECT_NAME}"}

# Convert the HELM_PLUGIN_PATH to unix if cygpath is
# available. This is the case when using MSYS2 or Cygwin
# on Windows where helm returns a Windows path but we
# need a Unix path

if type cygpath > /dev/null 2>&1; then
  HELM_PLUGIN_PATH=$(cygpath -u $HELM_PLUGIN_PATH)
fi

if [[ $SKIP_BIN_INSTALL == "1" ]]; then
  echo "Skipping binary install"
  exit
fi

# Discover the architecture for this system.
initArch() {
  ARCH=$(uname -m)
  case $ARCH in
    x86) ARCH="386";;
    x86_64) ARCH="amd64";;
    i686) ARCH="386";;
    i386) ARCH="386";;
  esac
}

# Discover the operating system for this system.
initOS() {
  OS=$(echo `uname`|tr '[:upper:]' '[:lower:]')

  case "$OS" in
    # Msys support
    msys*) OS='windows';;
    # Minimalist GNU for Windows
    mingw*) OS='windows';;
  esac
}

# Figure out the download url for the latest available version.
getDownloadURL() {
  if ! type "curl" > /dev/null && ! type "wget" > /dev/null; then
    echo "Either curl or wget is required"
    exit 1
  fi

  local url="https://api.github.com/repos/${PROJECT_GH}/releases/latest"
  local version=$(git describe --tags --exact-match 2>/dev/null)
  if [ -n "$version" ]; then
    url="https://api.github.com/repos/${PROJECT_GH}/releases/tags/${version}"
  fi
  
  # Use the GitHub API to find the download url for this project.
  if type "curl" > /dev/null; then
    DOWNLOAD_URL=$(curl -v -s $url | grep "${OS}-${ARCH}" | awk '/\"browser_download_url\":/{gsub( /[,\"]/,"", $2); print $2}')
  elif type "wget" > /dev/null; then
    DOWNLOAD_URL=$(wget -q -O - $url | grep "${OS}-${ARCH}" | awk '/\"browser_download_url\":/{gsub( /[,\"]/,"", $2); print $2}')
  fi

  if ! echo "${DOWNLOAD_URL}" | grep -q "${OS}-${ARCH}"; then
    echo "No prebuilt binary for ${OS}-${ARCH}."
    exit 1
  fi
}

# Download the plugin package.
downloadFile() {
  PLUGIN_TMP_FILE="/tmp/${PROJECT_NAME}.tgz"
  echo "Downloading $DOWNLOAD_URL"
  if type "curl" > /dev/null; then
    curl -L "$DOWNLOAD_URL" -o "$PLUGIN_TMP_FILE"
  elif type "wget" > /dev/null; then
    wget -q -O "$PLUGIN_TMP_FILE" "$DOWNLOAD_URL"
  fi
}

# Unpack and install the helm plugin
installFile() {
  HELM_TMP="/tmp/$PROJECT_NAME"
  mkdir -p "$HELM_TMP"
  tar xf "$PLUGIN_TMP_FILE" -C "$HELM_TMP"
  echo "Preparing to install into ${HELM_PLUGIN_PATH}"
  cp -r "$HELM_TMP" "$HELM_PLUGIN_PATH"
}

# Executed if an error occurs.
fail_trap() {
  result=$?
  if [ "$result" != "0" ]; then
    echo "Failed to install $PROJECT_NAME"
    echo "\tFor support, go to https://github.com/${PROJECT_GH}."
  fi
  exit $result
}

# Use the installed plugin's binary to make sure it is working.
testVersion() {
  set +e
  echo "$PROJECT_NAME installed into $HELM_PLUGIN_PATH/$PROJECT_BIN_NAME"
  $HELM_PLUGIN_PATH/$PROJECT_BIN_NAME -h
  set -e
}

# Execution

#Stop execution on any error
trap "fail_trap" EXIT
set -e
initArch
initOS
getDownloadURL
downloadFile
installFile
testVersion
