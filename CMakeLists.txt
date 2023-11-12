cmake_minimum_required(VERSION 3.15)
project(etcd-ydb LANGUAGES CXX)

set(BUILD_SHARED_LIBS OFF)
set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_EXTENSIONS OFF)

message(STATUS "Using toolchain file: ${CMAKE_TOOLCHAIN_FILE}")

include(cmake/conan.cmake)

conan_cmake_autodetect(settings)
conan_cmake_install(
    PATH_OR_REFERENCE ${CMAKE_SOURCE_DIR}
    INSTALL_FOLDER ${CMAKE_BINARY_DIR}
    BUILD missing
    REMOTE conancenter
    SETTINGS ${settings}
    ENV "CONAN_CMAKE_GENERATOR=${CMAKE_GENERATOR}"
    CONF "tools.cmake.cmaketoolchain:generator=${CMAKE_GENERATOR}"
)

include(${CMAKE_BINARY_DIR}/conan_paths.cmake)
include(cmake/userver.cmake)

file(GLOB_RECURSE PROTOS ${CMAKE_CURRENT_SOURCE_DIR}/proto/*.proto) # TODO [pavelbezpravel]: bruh...

include(GrpcTargets)
add_grpc_library(etcd-api-proto PROTOS ${PROTOS})

add_executable(${PROJECT_NAME} src/main.cpp)
target_link_libraries(${PROJECT_NAME}
    PRIVATE
       etcd-api-proto
)
