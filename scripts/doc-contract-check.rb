#!/usr/bin/env ruby
# Check the contract's modeled shapes, references and handler request/status facts.
require 'yaml'
require 'json'
require 'open3'

# Documents are UTF-8; the C locale would raise "invalid byte sequence in US-ASCII".
Encoding.default_external = Encoding::UTF_8
doc = YAML.load_file('.docs/api/openapi.yml')
shapes, status = Open3.capture2('go', 'run', 'scripts/doc-schema/main.go')
abort 'Go schema export failed' unless status.success?
expected = JSON.parse(shapes)
schemas = doc.fetch('components').fetch('schemas')
if ARGV.include?("--update-go-schemas")
  expected.each { |name, shape| schemas[name] = shape }
end
errors = []
# Refresh uses a refresh cookie even when the access session has expired.
refresh = doc.fetch('paths').fetch('/api/v1/auth/token_update').fetch('post')
errors << 'Refresh security mismatch' unless refresh['security'] == [{'RefreshCookieAuth'=>[]}]
refresh_scheme = doc.fetch('components').fetch('securitySchemes').fetch('RefreshCookieAuth', {})
errors << 'Refresh cookie scheme mismatch' unless refresh_scheme.values_at('type', 'in', 'name') == ['apiKey', 'cookie', 'beacon_refresh_token']
[['cookie', 'beacon_refresh_token'], ['cookie', 'beacon_csrf_token'], ['header', 'X-CSRF-Token']].each do |location, name|
  params = refresh.fetch('parameters', []).select { |p| p['in'] == location && p['name'] == name }
  errors << "Refresh credential/CSRF mismatch: #{name}" unless params.size == 1 && params[0]['required'] == true && params[0].dig('schema', 'type') == 'string'
end
expected.each { |name, shape| errors << "Go schema drift: #{name}" unless schemas[name] == shape }
schemas.each { |name, shape| errors << "Stale Go schema: #{name}" if shape.key?('x-go-source') && !expected.key?(name) }
walk = lambda do |value, location|
  case value
  when Hash
    if value.key?('$ref')
      target = value['$ref']
      resolved = target.start_with?('#/') && target.delete_prefix('#/').split('/').reduce(doc) { |node, key| node.is_a?(Hash) ? node[key] : nil }
      errors << "Invalid reference #{location}: #{target}" unless resolved
    end
    errors << "Unmodeled object #{location}" if value['additionalProperties'] == true
    if value['type'] == 'object' && !value.key?('properties') && !value['additionalProperties'].is_a?(Hash)
      errors << "Missing object properties #{location}"
    end
    if value['type'] == 'array' && !value.key?('items')
      errors << "Missing array items #{location}"
    end
    Array(value['required']).each do |key|
      errors << "Missing required property #{location}.#{key}" unless value.fetch('properties', {}).key?(key)
    end if value['required'].is_a?(Array)
    value.each do |key, child|
      errors << "Empty schema #{location}" if key == 'schema' && child == {}
      walk.call(child, "#{location}/#{key}")
    end
  when Array
    value.each_with_index { |child, i| walk.call(child, "#{location}/#{i}") }
  end
end
walk.call(doc, '#')
handler_files = {
  'user'=>'account', 'role'=>'account', 'resource'=>'account', 'auth'=>'auth',
  'container'=>'container', 'host'=>'host', 'audit'=>'audit', 'mail'=>'mail',
  'alarm'=>'alarm', 'dingtalk'=>'dingtalk', 'feishu'=>'feishu', 'wecom'=>'wecom',
  'system'=>'system', 'agent'=>'agent'
}
# agent 模块的 handler 与 service 同处 beacon/service/agent/agent.go。
handler_sources = { 'agent' => 'beacon/service/agent/agent.go' }
# Special handlers have explicit reviewed bindings. They are never silently skipped.
schema_ref = ->(name) { { '$ref' => "#/components/schemas/#{name}" } }
special = {
  '/health' => ['get', '200', 'application/json', schema_ref.call('HealthResponse')],
  '/ready' => ['get', '200', 'application/json', schema_ref.call('HealthResponse')],
  '/ws/terminal' => ['get', '101', nil, nil],
  '/ws/{id}' => ['get', '101', nil, nil],
  '/api/v1/index/index' => ['get', '200', 'text/plain', { 'type' => 'string' }],
  '/api/v1/agent/list' => ['get', '200', 'application/json', { 'type'=>'array', 'nullable'=>true, 'items'=>schema_ref.call('ModelAgent') }],
  '/api/v1/agent/upgrade' => ['post', '200', 'application/json', schema_ref.call('UpgradeResult'), 'UpgradeArgs'],
  '/api/v1/host/report' => ['post', '200', 'application/json', schema_ref.call('ReportResponse'), 'ReportMonitorReportArgs'],
  '/api/v1/host/get_install_token' => ['post', '200', 'application/json', schema_ref.call('InstallTokenResponse')],
  '/api/v1/host/install/ca' => ['get', '200', 'application/x-pem-file', { 'type'=>'string', 'format'=>'binary' }]
}
%w[install install/config install/package install/certs].each do |endpoint|
  shape = { 'type'=>'string' }
  shape['format'] = 'binary' if endpoint.end_with?('/package', '/certs')
  special["/api/v1/host/#{endpoint}"] = ['get', '200', '*/*', shape]
end
router = File.read('beacon/service/router.go')
count = 0
doc['paths'].each do |path, item|
  item.each do |method, operation|
    next unless %w[get post put patch delete options head].include?(method)
    count += 1
    responses = operation.fetch('responses')
    responses.each do |code, response|
      unknown = response.keys.reject { |key| %w[description content headers links $ref].include?(key) || key.start_with?('x-') }
      errors << "Invalid response properties #{path} #{code}: #{unknown}" unless unknown.empty?
      errors << "Missing response description #{path} #{code}" unless response.key?('description') || response.key?('$ref')
    end
    success = responses.select { |code, _| code.to_s.match?(/\A2\d\d\z/) }
    unless path.start_with?('/ws/')
      errors << "Missing success response #{method} #{path}" if success.empty?
      success.each do |code, response|
        errors << "Missing response body model #{path}" if code.to_s != '204' && response.fetch('content', {}).empty?
        errors << "204 has content #{path}" if code.to_s == '204' && response.key?('content')
      end
    end
    if special.key?(path)
      verb, code, media, shape, input = special.fetch(path)
      codes = responses.keys.select { |c| c.match?(/\A(?:2\d\d|101)\z/) }
      response_matches = media ? responses.dig(code, 'content', media, 'schema') == shape : responses.fetch(code, {}).fetch('content', {}).empty?
      errors << "Special response mismatch #{path}" unless method == verb && codes == [code] && response_matches
      if input
        errors << "Request type mismatch #{path}: #{input}" unless operation.dig('requestBody','content','application/json','schema') == schema_ref.call(input) && operation.dig('requestBody','required') == true
      elsif operation.key?('requestBody')
        errors << "Unexpected request body #{path}"
      end
      next
    end
    if path.start_with?('/api/v1/')
      group, endpoint = path.delete_prefix('/api/v1/').split('/', 2)
      mod = handler_files[group]
      if mod && !endpoint.include?('/')
        api = mod == 'dingtalk' ? 'dingTalkAPI' : mod + 'API'
        registration = router.match(/g(?:User|Role|Resource)?\.#{method.capitalize}\("\/#{Regexp.escape(endpoint)}", a\.#{api}\.(\w+)\)/)
        if registration
          function = registration[1]
          source = handler_sources.fetch(mod) { "beacon/service/#{mod}/api/#{mod}.go" }
          text = File.read(source)
          body = text[/^func \([^\n]+\) #{Regexp.escape(function)}\([^\n]+\{\n(.*?)(?=^func |\z)/m, 1]
          if body
            if body.include?('fiberx.ParseBody(') && !body.match?(/\b(?:var args |args := )schema\.\w+/)
              errors << "Cannot resolve request type #{path}"
            end
            if body.include?('fiberx.ParseQuery(') && !body.match?(/\bvar args schema\.\w+/)
              errors << "Cannot resolve query type #{path}"
            end
            unless body.include?('fiberx.ParseBody(') || body.include?('ctx.FormFile(')
              errors << "Unexpected request body #{path}" if operation.key?('requestBody')
            end
            if body.include?('fiberx.NoContent(')
              errors << "Handler returns 204: #{path}" unless success.keys == ['204']
            end
            if body.include?('fiberx.ParseBody(') && (args = body[/\b(?:var args |args := )schema\.(\w+)/, 1])
              actual = operation.dig('requestBody', 'content', 'application/json', 'schema', '$ref')
              content = operation.dig('requestBody', 'content') || {}
              errors << "Request type mismatch #{path}: #{args}" unless actual == "#/components/schemas/#{args}" && operation.dig('requestBody', 'required') == true && content.values.all? { |entry| entry['schema'] == { '$ref'=>actual } }
            end
            if body.include?('fiberx.ParseQuery(') && (args = body[/\bvar args schema\.(\w+)/, 1])
              if ARGV.include?('--update-go-schemas')
                fields = expected.fetch(args)
                parameters = operation.fetch('parameters', []).reject { |p| p['in'] == 'query' }
                fields.fetch('properties').each do |name, prop|
                  parameters << { 'in'=>'query', 'name'=>prop.fetch('x-go-query', name), 'required'=>fields.fetch('required', []).include?(name), 'schema'=>prop }
                end
                operation['parameters'] = parameters
              end
              params = operation.fetch('parameters', []).select { |p| p['in'] == 'query' }
              props = expected.fetch(args).fetch('properties')
              wanted = props.map { |name, prop| prop.fetch('x-go-query', name) }
              errors << "Query names mismatch #{path}" unless params.map { |p| p['name'] }.sort == wanted.sort
              props.each do |name, prop|
                wire_name = prop.fetch('x-go-query', name)
                param = params.find { |p| p['name'] == wire_name }
                next unless param
                errors << "Query schema mismatch #{path}: #{wire_name}" unless param['schema'] == prop
                required = expected.fetch(args).fetch('required', []).include?(name)
                errors << "Query required mismatch #{path}: #{wire_name}" unless param.fetch('required', false) == required
              end
            end
            if body.include?('ctx.FormFile("file")')
              upload = { 'type'=>'object', 'properties'=>{ 'file'=>{ 'type'=>'string', 'format'=>'binary' } }, 'required'=>['file'] }
              errors << "Multipart request mismatch #{path}" unless operation.dig('requestBody','content','multipart/form-data','schema') == upload && operation.dig('requestBody','required') == true
            end
            if function == 'FileDownload'
              errors << "Download response mismatch #{path}" unless success.keys == ['200'] && responses.dig('200','content','*/*','schema') == { 'type'=>'string','format'=>'binary' }
            end
            unless body.include?('fiberx.NoContent(') || body.include?('fiberx.Success(') || function == 'FileDownload'
              errors << "Cannot resolve success response #{path}"
            end
            if (reply = body[/return fiberx\.Success\(\w+, (.+)\)/, 1])
              type = reply[/\Aschema\.(\w+)\{/, 1]
              unless type
                method_name = body[/\b#{Regexp.escape(reply)}, err := \w+\.\w+\.(\w+)\(/, 1]
                service = File.read(handler_sources.fetch(mod) { "beacon/service/#{mod}/service/#{mod}.go" })
                type = service[/^func \([^\n]+\) #{Regexp.escape(method_name.to_s)}\([^\n]*\) \(schema\.(\w+), error\)/, 1] if method_name
              end
              if type
                ref = { '$ref' => "#/components/schemas/#{type}" }
                actual = responses.dig('200', 'content', 'application/json', 'schema')
                errors << "Response type mismatch #{path}: #{type}" unless actual == ref && success.keys == ['200']
              else
                errors << "Cannot resolve response type #{path}: #{reply}"
              end
            end
          else
            errors << "Missing handler body #{source}:#{function}"
          end
        else
          errors << "Missing handler registration #{path}"
        end
      else
        errors << "Unbound handler #{method} #{path}"
      end
    else
      errors << "Unbound handler #{method} #{path}"
    end
  end
end
standard_output, standard_status = Open3.capture2e({'GOWORK'=>'off'}, 'go', 'run', './cmd/openapi-check',
  stdin_data: JSON.generate(doc), chdir: File.expand_path('scripts/doc-schema'))
errors << "OpenAPI standard validation failed: #{standard_output.strip}" unless standard_status.success?
abort errors.join("\n") unless errors.empty?
if ARGV.include?('--update-go-schemas')
  File.write('.docs/api/openapi.yml', YAML.dump(JSON.parse(JSON.generate(doc))).sub(/^---\n/, ''))
end
puts "#{count} operations: success bodies/204, references, handler request types and #{expected.size} Go field models checked"
