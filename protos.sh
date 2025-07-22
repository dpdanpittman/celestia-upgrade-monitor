#!/bin/bash
set -e

echo "Setting up Celestia proto files and dependencies..."

# Create directory structure
mkdir -p celestia/{blob,signal,mint,minfee,qgb}/v1
mkdir -p celestia/core/v1/{gas_estimation,tx}
mkdir -p third_party/{gogoproto,cosmos_proto,google/api}
mkdir -p third_party/cosmos/{base/query/v1beta1,msg/v1}

# Download Celestia proto files
echo "Downloading Celestia proto files..."
git clone https://github.com/celestiaorg/celestia-app.git --depth=1 temp_celestia
cp -r temp_celestia/proto/celestia/* celestia/
rm -rf temp_celestia

# Download dependencies
echo "Downloading dependencies..."
git clone https://github.com/cosmos/gogoproto.git --depth=1 temp_gogoproto
cp temp_gogoproto/gogoproto/gogo.proto third_party/gogoproto/
rm -rf temp_gogoproto

git clone https://github.com/cosmos/cosmos-sdk.git --depth=1 temp_cosmos
cp -r temp_cosmos/proto/cosmos/base/query/v1beta1/* third_party/cosmos/base/query/v1beta1/
cp -r temp_cosmos/proto/cosmos/msg/v1/* third_party/cosmos/msg/v1/
cp -r temp_cosmos/proto/cosmos/* third_party/cosmos_proto/
rm -rf temp_cosmos

# Download cosmos_proto
git clone https://github.com/cosmos/cosmos-proto.git --depth=1 temp_cosmos_proto
cp temp_cosmos_proto/proto/cosmos_proto/cosmos.proto third_party/cosmos_proto/
rm -rf temp_cosmos_proto

git clone https://github.com/googleapis/googleapis.git --depth=1 temp_googleapis
cp temp_googleapis/google/api/annotations.proto third_party/google/api/
cp temp_googleapis/google/api/http.proto third_party/google/api/
rm -rf temp_googleapis

# Download protobuf dependencies
mkdir -p third_party/google/protobuf
curl -s https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/descriptor.proto > third_party/google/protobuf/descriptor.proto
curl -s https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/timestamp.proto > third_party/google/protobuf/timestamp.proto
curl -s https://raw.githubusercontent.com/protocolbuffers/protobuf/main/src/google/protobuf/field_mask.proto > third_party/google/protobuf/field_mask.proto

echo "Setup complete!"
