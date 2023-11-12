set(USERVER_FEATURE_MONGODB OFF)
set(USERVER_FEATURE_POSTGRESQL OFF)
set(USERVER_FEATURE_REDIS OFF)
set(USERVER_FEATURE_CLICKHOUSE OFF)
set(USERVER_FEATURE_RABBITMQ OFF)
set(USERVER_FEATURE_MYSQL OFF)
set(USERVER_USE_LD lld-14)

include(${CMAKE_SOURCE_DIR}/third_party/userver/cmake/SetupEnvironment.cmake)
include(GNUInstallDirs)

add_subdirectory(${CMAKE_SOURCE_DIR}/third_party/userver)
