#!/usr/bin/env ruby
# Regression tests use isolated fixtures, never mutate the working repository.
require 'minitest/autorun'
require 'tmpdir'
require 'fileutils'
require 'open3'
require 'yaml'
require 'json'

# Documents are UTF-8; the C locale would raise "invalid byte sequence in US-ASCII".
Encoding.default_external = Encoding::UTF_8

class DocumentationGateTest < Minitest::Test
  ROOT = File.expand_path('..', __dir__)

  def setup
    @dir = Dir.mktmpdir('beacon-doc-gate-')
    paths = %w[scripts/doc-evidence-check.rb scripts/doc-commands.rb scripts/doc-contract-check.rb scripts/doc-signal-check.py scripts/doc-schema/main.go scripts/doc-schema/main_test.go scripts/doc-schema/go.mod scripts/doc-schema/go.sum scripts/doc-schema/cmd/openapi-check/main.go .docs/project-analysis.md .docs/api/openapi.yml .docs/api/routes.md .docs/architecture.md beacon/service/router.go beacon/service/agent/agent.go beacon/service/agent/updater.go beacon/service/model/agent.go beacon/pkg/fiberx/fiberx.go beacon/web/package.json]
    paths += Dir.glob(File.join(ROOT, '.docs/modules/*.md')).map { |p| p.delete_prefix(ROOT + '/') }
    paths += Dir.glob(File.join(ROOT, '.specs/domain/*.md')).map { |p| p.delete_prefix(ROOT + '/') }
    paths += Dir.glob(File.join(ROOT, 'beacon/service/schema/*.go')).map { |p| p.delete_prefix(ROOT + '/') }
    paths += Dir.glob(File.join(ROOT, 'common/rpc/schema/*.go')).map { |p| p.delete_prefix(ROOT + '/') }
    paths += Dir.glob(File.join(ROOT, 'beacon/service/*/api/*.go')).map { |p| p.delete_prefix(ROOT + '/') }
    paths += Dir.glob(File.join(ROOT, 'beacon/service/*/service/*.go')).map { |p| p.delete_prefix(ROOT + '/') }
    # 项目分析文档中反引号标注的仓库路径也是 fixture 的一部分；测试清单里的
    # 每条测试路径因此都会出现在 fixture 中，门禁的清单比对才有意义。
    text = File.read(File.join(ROOT, '.docs/project-analysis.md'))
    paths += text.scan(/`([^`]+)`/).flatten.map { |p| p.split('::', 2).first }.select { |p| File.file?(File.join(ROOT, p)) }
    paths.uniq.each do |p|
      FileUtils.mkdir_p(File.dirname(File.join(@dir, p)))
      FileUtils.cp(File.join(ROOT, p), File.join(@dir, p))
    end
    # Existence of module directories is part of command validation.
    %w[beacon/web collia common common/docker common/logger common/timex].each { |p| FileUtils.mkdir_p(File.join(@dir, p)) }
    run!('git', 'init', '-q')
    run!('git', 'add', '.')
  end

  def teardown
    FileUtils.remove_entry(@dir)
  end

  def run!(*command)
    out, status = Open3.capture2e(*command, chdir: @dir)
    assert status.success?, out
    out
  end

  def reject!(*command)
    out, status = Open3.capture2e(*command, chdir: @dir)
    refute status.success?, "Unexpected pass: #{out}"
    out
  end

  def change(path)
    file = File.join(@dir, path)
    File.write(file, yield(File.read(file)))
  end

  def test_baseline
    run!('ruby', 'scripts/doc-evidence-check.rb')
    run!('ruby', 'scripts/doc-contract-check.rb')
  end

  def test_rejects_missing_domain_id
    change('.docs/project-analysis.md') { |s| s.lines.reject { |line| line.start_with?('| IA001 |') }.join }
    assert_includes reject!('ruby', 'scripts/doc-evidence-check.rb'), 'Domain evidence IDs'
  end

  def test_rejects_renamed_test
    change('.docs/project-analysis.md') { |s| s.sub('::TestUserDeleteArgsRequiresExplicitUUIDTargets', '::TestRemovedEvidence') }
    assert_includes reject!('ruby', 'scripts/doc-evidence-check.rb'), 'Missing test declaration'
  end

  def test_rejects_duplicate_test_inventory
    change('.docs/project-analysis.md') { |s| s.sub('<!-- test-inventory:start -->', "<!-- test-inventory:start -->\n- `beacon/service/schema/account_test.go`") }
    assert_includes reject!('ruby', 'scripts/doc-evidence-check.rb'), 'Test inventory'
  end

  def test_rejects_new_command_without_result
    change('.docs/modules/common.md') { |s| s.sub("cd common\n", "cd common\ngo test -count=1 ./...\n") }
    assert_includes reject!('ruby', 'scripts/doc-evidence-check.rb'), 'Command evidence'
  end

  def test_rejects_go_field_model_drift
    change('.docs/api/openapi.yml') do |s|
      doc = YAML.safe_load(s)
      doc['components']['schemas']['LoginArgs']['properties']['username']['type'] = 'integer'
      YAML.dump(doc)
    end
    assert_includes reject!('ruby', 'scripts/doc-contract-check.rb'), 'Go schema drift'
  end

  def test_rejects_unmodeled_success_response
    change('.docs/api/openapi.yml') do |s|
      doc = YAML.safe_load(s)
      doc['paths']['/api/v1/host/cpu_info']['get']['responses']['200'].delete('content')
      YAML.dump(doc)
    end
    assert_includes reject!('ruby', 'scripts/doc-contract-check.rb'), 'Missing response body model'
  end

  def test_rejects_handler_request_mismatch
    change('.docs/api/openapi.yml') do |s|
      doc = YAML.safe_load(s)
      doc['paths']['/api/v1/user/user_delete']['post']['requestBody']['content']['application/json']['schema']['$ref'] = '#/components/schemas/RoleDeleteArgs'
      YAML.dump(doc)
    end
    assert_includes reject!('ruby', 'scripts/doc-contract-check.rb'), 'Request type mismatch'
  end

  def test_normalizes_continuations_and_module_directories
    commands = JSON.parse(run!('ruby', 'scripts/doc-commands.rb'))
    docker = commands.select { |c| c['command'].start_with?('docker buildx') }
    assert_equal 1, docker.size
    assert_includes docker.first['command'], '-t beacon:size-check'
    assert_equal 'beacon', docker.first['cwd']
    assert commands.select { |c| c['doc'].end_with?('/collia.md') }.all? { |c| c['cwd'] == 'collia' }
  end

  def test_rejects_placeholder_object
    change('.docs/api/openapi.yml') do |s|
      doc = YAML.safe_load(s)
      doc['paths']['/api/v1/host/cpu_info']['get']['responses']['200']['content']['application/json']['schema'] = { 'type' => 'object' }
      YAML.dump(doc)
    end
    assert_includes reject!('ruby', 'scripts/doc-contract-check.rb'), 'Missing object properties'
  end

  def test_rejects_invalid_response_keys_from_yaml_flow_mapping
    change('.docs/api/openapi.yml') do |s|
      doc = YAML.safe_load(s)
      doc['paths']['/api/v1/agent/upgrade']['post']['responses']['400']['unquoted comma clause'] = nil
      YAML.dump(doc)
    end
    assert_includes reject!('ruby', 'scripts/doc-contract-check.rb'), 'Invalid response properties'
  end

  def mutate_contract
    change('.docs/api/openapi.yml') do |text|
      doc = YAML.safe_load(text)
      yield doc
      YAML.dump(doc)
    end
    reject!('ruby', 'scripts/doc-contract-check.rb')
  end

  def test_rejects_wrong_response_binding
    assert_includes mutate_contract { |d|
      d['paths']['/api/v1/host/cpu_info']['get']['responses']['200']['content']['application/json']['schema'] = { '$ref'=>'#/components/schemas/LoginArgs' }
    }, 'Response type mismatch'
  end

  def test_rejects_empty_schema
    assert_includes mutate_contract { |d|
      d['paths']['/api/v1/host/cpu_info']['get']['responses']['200']['content']['application/json']['schema'] = {}
    }, 'Empty schema'
  end

  def test_rejects_wrong_query_type
    assert_includes mutate_contract { |d|
      d['paths']['/api/v1/container/containers']['get']['parameters'].find { |p| p['name']=='page' }['schema'] = { 'type'=>'string' }
    }, 'Query schema mismatch'
  end

  def test_rejects_old_query_names
    assert_includes mutate_contract { |d|
      d['paths']['/api/v1/host/cpu_trending']['get']['parameters'].find { |p| p['name']=='start_time' }['name'] = 'StartTime'
    }, 'Query names mismatch'
  end

  def test_rejects_wrong_upgrade_body
    assert_includes mutate_contract { |d|
      d['paths']['/api/v1/agent/upgrade']['post']['requestBody']['content']['application/json']['schema'] = { '$ref'=>'#/components/schemas/LoginArgs' }
    }, 'Request type mismatch'
  end

  def test_rejects_extra_success_on_no_content_handler
    assert_includes mutate_contract { |d|
      d['paths']['/api/v1/user/user_delete']['post']['responses']['200'] = { 'description'=>'wrong', 'content'=>{ 'application/json'=>{ 'schema'=>{ 'type'=>'string' } } } }
    }, 'Handler returns 204'
  end

  def test_rejects_service_return_type_drift
    change('beacon/service/host/service/host.go') { |s|
      s.sub('func (h *HostService) CPUInfo(ctx context.Context) (schema.CPUInfoReply, error)', 'func (h *HostService) CPUInfo(ctx context.Context) (schema.MemoryInfoReply, error)')
    }
    assert_includes reject!('ruby', 'scripts/doc-contract-check.rb'), 'Response type mismatch'
  end

  def test_rejects_wrong_query_requirement
    assert_includes mutate_contract { |d|
      d['paths']['/api/v1/container/containers']['get']['parameters'].find { |p| p['name']=='page' }['required'] = false
    }, 'Query required mismatch'
  end

  def test_rejects_deleted_tracked_inventory_file
    FileUtils.rm(File.join(@dir, 'beacon/service/schema/account_test.go'))
    assert_includes reject!('ruby', 'scripts/doc-evidence-check.rb'), 'Test inventory'
  end

  def test_rejects_unknown_security_scheme
    assert_includes mutate_contract { |d|
      d['paths']['/api/v1/user/user_delete']['post']['security'] = [{ 'DoesNotExist'=>[] }]
    }, 'undefined security scheme'
  end

  def test_rejects_invalid_standard_schema_type
    assert_includes mutate_contract { |d|
      # HealthResponse is a reviewed inline-map model, outside the Go struct exporter.
      d['components']['schemas']['HealthResponse']['properties']['status']['type'] = 'not-a-json-type'
    }, 'OpenAPI standard validation failed'
  end

  def test_rejects_special_model_drift
    assert_includes mutate_contract { |d|
      d['components']['schemas']['UpgradeResult']['properties']['triggered']['type'] = 'string'
    }, 'Go schema drift: UpgradeResult'
  end

  def test_rejects_special_source_field_drift
    change('beacon/service/agent/updater.go') { |s| s.sub('Triggered      bool', 'Triggered      string') }
    assert_includes reject!('ruby', 'scripts/doc-contract-check.rb'), 'Go schema drift: UpgradeResult'
  end

  def test_rejects_refresh_auth_and_csrf_drift
    original = File.read(File.join(@dir, '.docs/api/openapi.yml'))
    [
      ->(d) { d['paths']['/api/v1/auth/token_update']['post'].delete('security') },
      ->(d) { d['components']['securitySchemes']['RefreshCookieAuth']['name'] = 'beacon_access_token' },
      ->(d) { d['paths']['/api/v1/auth/token_update']['post']['parameters'].reject! { |p| p['name'] == 'X-CSRF-Token' } }
    ].each do |mutate|
      doc = YAML.safe_load(original)
      mutate.call(doc)
      File.write(File.join(@dir, '.docs/api/openapi.yml'), YAML.dump(doc))
      assert_includes reject!('ruby', 'scripts/doc-contract-check.rb'), 'Refresh'
    end
  end

  def test_rejects_route_index_count_drift
    change('.docs/api/routes.md') { |s| s.sub('77 条业务路由', '79 条业务路由') }
    assert_includes reject!('ruby', 'scripts/doc-contract-check.rb'), 'Route index count mismatch'
  end

  def test_rejects_route_index_path_drift
    change('.docs/api/routes.md') { |s| s.lines.reject { |line| line.include?('`/api/v1/host/cpu_info`') }.join }
    assert_includes reject!('ruby', 'scripts/doc-contract-check.rb'), 'Route index missing'
  end

  def test_rejects_route_index_handler_drift
    change('.docs/api/routes.md') { |s| s.sub('`wecomAPI.Query`', '`wrongAPI.Query`') }
    assert_includes reject!('ruby', 'scripts/doc-contract-check.rb'), 'Route handler mismatch'
  end

  def test_rejects_unknown_persistent_domain_reference
    change('.docs/modules/beacon.md') { |s| s.sub('（IA007）', '（IR999）') }
    assert_includes reject!('ruby', 'scripts/doc-evidence-check.rb'), 'Unknown Domain reference'
  end

  def test_signal_checker_baseline
    run!('python3', 'scripts/doc-signal-check.py')
  end

  def test_rejects_go_test_signal
    change('.docs/modules/common.md') do |s|
      s.sub('- `interface AgentLifecycle', "- `func Fake (common/rpc/tunnel/server_test.go)`\n- `interface AgentLifecycle")
    end
    assert_includes reject!('python3', 'scripts/doc-signal-check.py'), 'Test file in entry signal'
  end

  def test_rejects_common_child_module_signal
    change('.docs/modules/common.md') do |s|
      s.sub('- 配置：`Option` in `common/database/options.go`', "- 配置：`Option` in `common/database/options.go`\n- 配置：`Option` in `common/logger/config.go`")
    end
    assert_includes reject!('python3', 'scripts/doc-signal-check.py'), 'Child module signal'
  end

  def test_business_contract_samples
    binary = File.join(@dir, 'openapi-check')
    output, status = Open3.capture2e({'GOWORK'=>'off'}, 'go', 'build', '-o', binary, './cmd/openapi-check', chdir: File.join(@dir, 'scripts/doc-schema'))
    assert status.success?, output
    document = JSON.generate(YAML.safe_load(File.read(File.join(@dir, '.docs/api/openapi.yml'))))
    password = { 'username'=>'admin', 'old_password'=>'old-password' }
    cases = [
      ['PasswordUpdateArgs', password.merge('new_password'=>'x'), false],
      ['PasswordUpdateArgs', password.merge('new_password'=>'abcDEF123456'), true],
      ['PasswordUpdateArgs', password.merge('new_password'=>'密码abcDEF123456'), true],
      ['ReportMonitorReportArgs', {}, false],
      ['ReportMonitorReportArgs', {'agent_id'=>''}, false],
      ['ReportMonitorReportArgs', {'agent_id'=>'agent/other'}, false],
      ['ReportMonitorReportArgs', {'agent_id'=>'a'*129}, false],
      ['ReportMonitorReportArgs', {'agent_id'=>'agent-a.1'}, true],
      ['FilesSearchArgs', {'path'=>''}, false],
      ['FilesSearchArgs', {'path'=>'/'}, true],
      ['UserCreateArgs', {'username'=>'audituser','password'=>'Strong-Pass9!','status'=>0,'role_ids'=>['r']}, false],
      ['UserCreateArgs', {'username'=>'audituser','password'=>'Strong-Pass9!','status'=>1,'role_ids'=>['r']}, true]
    ]
    [-1, 0, 1, 2, 3].each do |state|
      cases << ['UserCreateArgs', {'username'=>'audituser','password'=>'Strong-Pass9!','status'=>state,'role_ids'=>['r']}, [1,2].include?(state)]
      cases << ['RoleCreateArgs', {'name'=>'operator','status'=>state}, [1,2].include?(state)]
      %w[UserUpdateArgs RoleUpdateArgs].each do |name|
        cases << [name, {'id'=>'019f8dbc-8781-7af3-bf49-e7fafa912e87','status'=>state}, [0,1,2].include?(state)]
      end
    end
    cases.each do |name, sample, accepted|
      output, status = Open3.capture2e(binary, '--schema', name, '--sample', JSON.generate(sample), stdin_data: document, chdir: @dir)
      assert_equal accepted, status.success?, "#{name} sample acceptance mismatch: #{output}"
    end
  end
end
