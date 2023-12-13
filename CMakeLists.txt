cmake_minimum_required(VERSION 3.22)
project(etcd-ydb LANGUAGES CXX)

set(BUILD_SHARED_LIBS OFF)
set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_EXTENSIONS OFF)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

message(STATUS "Using toolchain file: ${CMAKE_TOOLCHAIN_FILE}")

include(cmake/conan.cmake)
conan_cmake_install(
    PATH_OR_REFERENCE ${CMAKE_SOURCE_DIR}
    INSTALL_FOLDER ${CMAKE_BINARY_DIR}
    BUILD missing
    ENV "CONAN_CMAKE_GENERATOR=${CMAKE_GENERATOR}"
    CONF "tools.cmake.cmaketoolchain:generator=${CMAKE_GENERATOR}"
    PROFILE "etcd-ydb/default"
    PROFILE_BUILD "etcd-ydb/default"
    PROFILE_HOST "etcd-ydb/default"
    SETTINGS "build_type=${CMAKE_BUILD_TYPE}"
    SETTINGS_BUILD "build_type=${CMAKE_BUILD_TYPE}"
    SETTINGS_HOST "build_type=${CMAKE_BUILD_TYPE}"
)
include(${CMAKE_BINARY_DIR}/conan_toolchain.cmake)

find_package(userver 1.0.0 EXACT REQUIRED)

file(GLOB_RECURSE PROTOS ${CMAKE_SOURCE_DIR}/proto/*.proto)
add_grpc_library(etcd-api-proto PROTOS ${PROTOS})

add_executable(${PROJECT_NAME}
        src/components/etcd/api/KV_service_component.cpp
        src/components/etcd/api/KV_service_component.h
        src/main.cpp)
target_link_libraries(${PROJECT_NAME}
    PRIVATE
       etcd-api-proto
)
target_compile_options(${PROJECT_NAME}
    PRIVATE
        -Werror
        -Wall
        -Wextra
        -Wassign-enum
        -Wcast-align
        -Wcast-qual
        -Wconversion
        -Wextra-semi
        -Wsign-compare
        -Wsign-conversion
        -pedantic
        -pedantic-errors
)
