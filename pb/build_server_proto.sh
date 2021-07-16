#! /bin/bash

OUT_PATH="."

protoc -I . --go_out=plugins=grpc:"$OUT_PATH" *.proto

if [ $? != 0 ]; then
  echo "build protobuf fail!"
  exit
fi

rm -rf $OUT_PATH/*.pb.go
protoc -I . --go_out=plugins=grpc:"$OUT_PATH" *.proto

# cur_pwd=$(pwd)
# cd $OUT_PATH
# ls *.pb.go | xargs -n1 -IX bash -c 'sed s/,omitempty// X > X.tmp && mv X{.tmp,}'
# cd $cur_pwd

echo "build protobuf success!"
