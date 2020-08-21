protoc --proto_path ../../../ -I=./proto --go_out=plugins=grpc:./proto proto/executor.proto
mv proto/github.com/brotherlogic/executor/proto/* ./proto
