#!/usr/bin/env bash
#Quick build script for all the plugins
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
export GOPATH="${cdir}/.gopath"

#Move into the src-plugins dir and start building all the individual plugins
cd "${cdir}/src-plugins"
for plugin in `ls`
do
  if [ ! -e "${cdir}/src-plugins/${plugin}/manifest.json" ] ; then continue ; fi
  cd "${cdir}/src-plugins/${plugin}"
  go get
  go build
  if [ $? -eq 0 ] ; then
    echo "Plugin Created: ${plugin}"
    if [ -n "${OUTDIR}" ] ; then
      mkdir -p "${OUTDIR}/${plugin}"
      cp "${cdir}/src-plugins/${plugin}/${plugin}" "${OUTDIR}/${plugin}/${plugin}"
      cp "${cdir}/src-plugins/${plugin}/manifest.json" "${OUTDIR}/${plugin}/manifest.json"
      echo " - Installed to ${OUTDIR}/${plugin}"
    fi
  else
    echo "[ERROR] Plugin Could not be created: ${plugin}"
  fi
done
unset GOPATH
