#ifndef ETCD_YDB_KV_SERVICE_COMPONENT_H
#define ETCD_YDB_KV_SERVICE_COMPONENT_H

#include <string>
#include <string_view>

#include <userver/components/component_config.hpp>
#include <userver/components/component_context.hpp>

#include <etcd/api/etcdserverpb/rpc_service.usrv.pb.hpp>

namespace etcdserverpb {
class KVServiceComponent final : public KVBase::Component {
 public:
  static constexpr std::string_view kName = "kv-service";

  KVServiceComponent(const userver::components::ComponentConfig& config,
                     const userver::components::ComponentContext& context)
      : KVBase::Component(config, context) {}

  // TODO [pavelbezpravel]: stub.
  void Range(RangeCall& call, RangeRequest&& request) override;
};
}  // namespace etcdserverpb

#endif  // ETCD_YDB_KV_SERVICE_COMPONENT_H
