#include "KV_service_component.h"

#include <userver/yaml_config/merge_schemas.hpp>

namespace etcdserverpb {
void KVServiceComponent::Range(RangeCall& call, RangeRequest&& request) {
  RangeResponse response{};
  call.Finish(response);
}

userver::yaml_config::Schema KVServiceComponent::GetStaticConfigSchema() {
  return userver::yaml_config::MergeSchemas<
      userver::ugrpc::server::ServiceComponentBase>(R"(
type: object
description: gRPC KV service component
additionalProperties: false
properties:
    greeting-prefix:
        type: string
        description: greeting prefix
)");
}
}  // namespace etcdserverpb
