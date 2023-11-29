#include "KV_service_component.h"

namespace etcdserverpb {
void KVServiceComponent::Range(RangeCall& call, RangeRequest&& request) {
  RangeResponse response{};
  call.Finish(response);
}
}  // namespace etcdserverpb
