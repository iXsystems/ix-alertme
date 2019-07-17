#!/bin/sh
# Helper script which will create the port / distfiles
# from a checked out git repo

# Set the port
port="misc/ix-alertme"
dfile="ix-alertme"

massage_subdir() {
  cd "$1"
  if [ $? -ne 0 ] ; then
     echo "SKIPPING $i"
     continue
  fi

comment="`cat Makefile | grep 'COMMENT ='`"

  echo "# \$FreeBSD\$
#

$comment
" > Makefile.tmp

  for d in `ls`
  do
    if [ "$d" = ".." ]; then continue ; fi
    if [ "$d" = "." ]; then continue ; fi
    if [ "$d" = "Makefile" ]; then continue ; fi
    if [ ! -f "$d/Makefile" ]; then continue ; fi
    echo "    SUBDIR += $d" >> Makefile.tmp
  done
  echo "" >> Makefile.tmp
  echo ".include <bsd.port.subdir.mk>" >> Makefile.tmp
  mv Makefile.tmp Makefile

}

if [ -z "$1" ] ; then
   echo "Usage: ./mkport.sh <portstree> <distfiles>"
   exit 1
fi

if [ ! -d "${1}/Mk" ] ; then
   echo "Invalid directory: $1"
   exit 1
fi

portsdir="${1}"
if [ -z "$portsdir" -o "${portsdir}" = "/" ] ; then
  portsdir="/usr/ports"
fi
#Set the specific env variable to use the custom ports dir location
PORTSDIR=${portsdir}

if [ -z "$2" ] ; then
  distdir="${portsdir}/distfiles"
else
  distdir="${2}"
fi
if [ ! -d "$distdir" ] ; then
  mkdir -p ${distdir}
fi

# Get the GIT tag
ghtag=`git log -n 1 | grep '^commit ' | awk '{print $2}'`

# Get the version
if [ -e "version" ] ; then
  verTag=$(cat version)
else
  verTag=$(date '+%Y%m%d%H%M')
fi

# Cleanup old distfiles
rm ${distdir}/${dfile}-* 2>/dev/null

# Copy ports files
if [ -d "${portsdir}/${port}" ] ; then
  rm -rf ${portsdir}/${port} 2>/dev/null
fi
cp -r ${dfile} ${portsdir}/${port}

# Set the version numbers
sed -i '' "s|%%CHGVERSION%%|${verTag}|g" ${portsdir}/${port}/Makefile
sed -i '' "s|%%GHTAG%%|${ghtag}|g" ${portsdir}/${port}/Makefile

# Create the makesums / distinfo file
cd "${portsdir}/${port}"
make makesum
if [ $? -ne 0 ] ; then
  echo "Failed makesum"
  exit 1
fi

# Remove the pkg-plist file (auto-generated at build time now)
#if [ -e "pkg-plist" ] ; then
#  rm "pkg-plist"
#fi
#make stage
#make makeplist | grep -v "check/what/makeplist/gives/you" > pkg-plist
#make clean

# Update port cat Makefile
tcat=$(echo $port | cut -d '/' -f 1)
massage_subdir ${portsdir}/${tcat}
