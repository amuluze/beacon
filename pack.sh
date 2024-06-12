#!/bin/bash
#####################################################################
# Author     : wangjialong
# Date       : 2024/06/12 15:33:01
# Description:
#####################################################################

set -ex

# web
rm -rf amprobe/dist
cd web
rm -rf dist
make build
cp -r dist ../amprobe/dist
cd ..

# amprobe
cd amprobe
make wire
make amd
cd ..

# amvector
cd amvector
rm amvector
make wire
make bin
cd ..

docker save amuluze/amprobe:v1.3.4 | gzip > amprobe.tar.gz

FILELIST="run.sh amvector/amvector amvector/configs amprobe.tar.gz amprobe/configs amprobe/nginx"
zip -0qr amprobe.installer.zip ${FILELIST}
rm amprobe.tar.gz
