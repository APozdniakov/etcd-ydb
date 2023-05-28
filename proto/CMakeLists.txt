cmake_minimum_required(VERSION 3.15)
project(proto)

get_target_property(grpc_cpp_plugin_location gRPC::grpc_cpp_plugin LOCATION)

# etcd-api-proto

add_library(etcd-api-proto)
target_sources(etcd-api-proto
    PRIVATE
        gogoproto/gogo.proto
        google/api/http.proto
        google/api/annotations.proto
        protoc-gen-swagger/options/openapiv2.proto
        protoc-gen-swagger/options/annotations.proto
        etcd/api/mvccpb/kv.proto
        etcd/api/authpb/auth.proto
        etcd/api/versionpb/version.proto
        etcd/api/etcdserverpb/rpc.proto
        etcd/server/lease/leasepb/lease.proto
)
target_link_libraries(etcd-api-proto
    PUBLIC
        gRPC::grpc
        gRPC::grpc++
        protobuf::libprotobuf
)
target_include_directories(etcd-api-proto
    PUBLIC
        ${CMAKE_CURRENT_BINARY_DIR}
)
protobuf_generate(TARGET etcd-api-proto LANGUAGE cpp)
protobuf_generate(
    TARGET
        etcd-api-proto
    LANGUAGE
        grpc
    GENERATE_EXTENSIONS
        .grpc.pb.h
        .grpc.pb.cc
    PLUGIN
        "protoc-gen-grpc=${grpc_cpp_plugin_location}"
)

# etcd-api-concurrency-proto

add_library(etcd-api-concurrency-proto)
target_sources(etcd-api-concurrency-proto
    PRIVATE
        gogoproto/gogo.proto
        google/api/http.proto
        google/api/annotations.proto
        protoc-gen-swagger/options/openapiv2.proto
        protoc-gen-swagger/options/annotations.proto
        etcd/api/mvccpb/kv.proto
        etcd/api/authpb/auth.proto
        etcd/api/versionpb/version.proto
        etcd/api/etcdserverpb/rpc.proto
        etcd/server/etcdserver/api/v3election/v3electionpb/v3election.proto
        etcd/server/etcdserver/api/v3lock/v3lockpb/v3lock.proto
)
target_link_libraries(etcd-api-concurrency-proto
    PUBLIC
        gRPC::grpc
        gRPC::grpc++
        protobuf::libprotobuf
)
target_include_directories(etcd-api-concurrency-proto
    PUBLIC
        ${CMAKE_CURRENT_BINARY_DIR}
)
protobuf_generate(TARGET etcd-api-concurrency-proto LANGUAGE cpp)
protobuf_generate(
    TARGET
        etcd-api-concurrency-proto
    LANGUAGE
        grpc
    GENERATE_EXTENSIONS
        .grpc.pb.h
        .grpc.pb.cc
    PLUGIN
        "protoc-gen-grpc=${grpc_cpp_plugin_location}"
)
