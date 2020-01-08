#!/bin/bash
#Quick build script for the ix-alertme utility
# INPUTS
OUTDIR="${1}"
if [ -n "${OUTDIR}" ] ; then
  OUTDIR=`realpath -q "${OUTDIR}"`
fi
#Verify that the output directory exists
if [ ! -d "${OUTDIR}" ] && [ -n "${OUTDIR}" ] ; then
  mkdir -p "${OUTDIR}"
fi

#Get the current directory of this build script
cdir=`dirname "${0}"`
cdir=`realpath -q "${cdir}"`
export GOPATH="/tmp/.gopath"
## TEMPORARY FIX: package source fetch issue from upstream
if [ -d "${GOPATH}/src/github.com/pierrec/lz4/v3" ] ; then
  rm -rf "${GOPATH}/src/github.com/pierrec/lz4/v3"
fi
mkdir -p "${GOPATH}/src/github.com/pierrec/lz4/v3"
git clone -q --depth=1 --branch v3.1.0 "https://github.com/pierrec/lz4" "${GOPATH}/src/github.com/pierrec/lz4/v3"

#Move into the source dir and build the tool
cd "${cdir}/src-go/ix-alertme"
go get
go build
if [ $? -eq 0 ] ; then
  unset GOPATH
  echo "ix-alertme built successfully"
  if [ -n "${OUTDIR}" ] ; then
    mv ix-alertme "${OUTDIR}/ix-alertme"
    echo " - Installed to ${OUTDIR}/ix-alertme"
  fi
else
  unset GOPATH
  echo "[ERROR] ix-alertme was not built"
  exit 1
fi
