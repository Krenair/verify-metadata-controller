#!/bin/bash
cp ./build/distributions/mdgen.zip ./
cp -R docker/integration_tests ./
docker run -it -e HSM_USER -e HSM_PASSWORD -v $(pwd)/test:/test -v $(pwd)/mdgen.zip:/mdgen.zip --mount src=$(pwd)/integration_tests,target=/integration_tests,type=bind $1
