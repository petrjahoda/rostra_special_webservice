#!/usr/bin/env bash
cd linux
upx rostra_special_webservice_linux
cd ..
docker rmi -f petrjahoda/rostra_special_webservice:latest
docker  build -t petrjahoda/rostra_special_webservice:latest .
docker push petrjahoda/rostra_special_webservice:latest