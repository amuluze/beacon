#!/usr/bin/env bash
# Explicit executable validation, separate from the fast documentation inventory gate.
# Does not start services, generate sources, mutate dependencies or build release images.
set -uo pipefail
cd "$(dirname "$0")/.." || exit 1
# Ruby reads UTF-8 documents; force the encoding so the result never depends on the
# caller's locale (a C locale raises "invalid byte sequence in US-ASCII").
export RUBYOPT="-EUTF-8${RUBYOPT:+ $RUBYOPT}"
result=0
# Module tests that need a Docker daemon and registry access are opt-in; the
# inventory classifies them so the skip is explicit rather than a silent pass.
docker_test_modules=$(ruby scripts/doc-commands.rb | ruby -rjson -e 'puts JSON.parse(STDIN.read).select { |c| c["requires"] == "docker" && c["command"] == "go test -race ./..." }.map { |c| c["cwd"] }.uniq.join(" ")')
# scripts/doc-schema 是独立的 Go module，仓库 go.work 不收录它；验证时切回它
# 自己的 module 上下文，避免 workspace 模式拒绝不在 use 列表中的模块。
run_module() {
  local module="$1"
  shift
  if [ "$module" = "scripts/doc-schema" ]; then
    (cd "$module" && GOWORK=off "$@")
  else
    (cd "$module" && "$@")
  fi
}
for module in beacon collia common common/docker common/logger common/timex scripts/doc-schema; do
  for action in build test vet; do
    printf '\n[%s] go %s\n' "$module" "$action"
    case "$action" in
      test)
        if [[ " $docker_test_modules " == *" $module "* && "${BEACON_DOC_VERIFY_DOCKER:-0}" != "1" ]]; then
          printf '  已跳过：%s 的测试依赖 Docker 守护进程与镜像仓库；设置 BEACON_DOC_VERIFY_DOCKER=1 开启\n' "$module"
          continue
        fi
        run_module "$module" go test -race ./... || result=1 ;;
      build)
        # A module with a single main package would drop a binary into its directory.
        if [ "$(run_module "$module" go list ./... | wc -l | tr -d ' ')" = "1" ]; then
          run_module "$module" go build -o /dev/null ./... || result=1
        else
          run_module "$module" go build ./... || result=1
        fi ;;
      vet) run_module "$module" go vet ./... || result=1 ;;
    esac
  done
done
ruby scripts/doc-contract-check.rb || result=1
ruby scripts/doc-check-test.rb || result=1
bash scripts/doc-check.sh || result=1
exit "$result"
