#!/usr/bin/env ruby
require 'set'
require_relative 'doc-commands'

# Documents are UTF-8; the C locale would raise "invalid byte sequence in US-ASCII".
Encoding.default_external = Encoding::UTF_8

text = File.read('.docs/project-analysis.md')
errors = []
files = IO.popen(['git', 'ls-files', '--cached', '--others', '--exclude-standard'], &:read).lines.map(&:strip).uniq
actual_tests = files.select { |f| File.file?(f) && f.end_with?('_test.go', '.test.ts', '.spec.ts') }.sort
inventory = text[/<!-- test-inventory:start -->\n(.*?)<!-- test-inventory:end -->/m, 1]
listed = inventory.to_s.scan(/^- `([^`]+)`$/).flatten
errors << 'Test inventory must list every test path exactly once' unless listed.sort == actual_tests && listed.uniq == listed
ids = Dir.glob('.specs/domain/*.md').flat_map { |f| File.read(f).scan(/^\| ([A-Z]{1,3}\d{3}) \|/).flatten }
# Persistent implementation documents must not refer to constraints that only
# existed in discarded task notes. Cryptographic names such as SHA256 are not IDs.
persistent_docs = Dir.glob('.docs/**/*.md') + %w[
  AGENTS.md README.md README.en.md
  beacon/AGENTS.md beacon/web/AGENTS.md collia/AGENTS.md common/AGENTS.md
].select { |f| File.file?(f) }
non_constraint_tokens = %w[SHA256]
persistent_docs.each do |file|
  refs = File.read(file).scan(/\b[A-Z]{1,3}\d{3}\b/).uniq - non_constraint_tokens
  refs.each { |ref| errors << "Unknown Domain reference #{file}: #{ref}" unless ids.include?(ref) }
end
section = text[/^### Evidence Backlink Table\n(.*?)(?=^## |^### |\z)/m, 1].to_s
rows = section.lines.select { |line| line.match?(/^\| [A-Z]{1,3}\d{3} \|/) }
seen = []
rows.each do |line|
  cells = line.split('|').map(&:strip)[1..]
  id, refs, assertion, state, gap = cells
  seen << id
  errors << "Missing assertion/gap/state #{id}" unless assertion && !assertion.empty? && %w[部分覆盖 待专项验收].include?(state) && gap && !gap.empty?
  references = refs.to_s.scan(/`([^`]+)`/).flatten
  errors << "Missing evidence #{id}" if references.empty?
  references.each do |reference|
    file, name = reference.split('::', 2)
    unless File.file?(file)
      errors << "Missing evidence path #{id}: #{file}"
      next
    end
    if name
      source = File.read(file)
      declared = if file.end_with?('_test.go')
        source.match?(/^func #{Regexp.escape(name)}\(t \*testing\.T\)/)
      elsif file.end_with?('.rb')
        source.match?(/^  def #{Regexp.escape(name)}\s*$/)
      else
        source.match?(/\b(?:it|test)\(['"]#{Regexp.escape(name)}['"]/)
      end
      errors << "Missing test declaration #{id}: #{reference}" unless declared
    end
  end
end
errors << 'Domain evidence IDs must match specs exactly once' unless seen.sort == ids.sort && seen.uniq == seen
command_section = text[/<!-- command-evidence:start -->\n(.*?)<!-- command-evidence:end -->/m, 1].to_s
recorded = {}
command_section.lines.grep(/^\| [a-f0-9]{12} \|/).each do |line|
  id, doc, cwd, command, result, detail = line.split('|').map(&:strip)[1..].map { |v| v.delete_prefix('`').delete_suffix('`').gsub('&#124;', '|') }
  errors << "Duplicate command evidence #{id}" if recorded.key?(id)
  recorded[id] = [doc, cwd, command]
  errors << "Missing command result/detail #{id}" unless %w[通过 失败 未执行].include?(result) && detail && !detail.empty?
end
commands = module_commands
expected = commands.to_h { |c| [c['id'], [c['doc'], c['cwd'], c['command']]] }
errors << 'Command evidence must match normalized module commands (path, directory and command)' unless recorded == expected
abort errors.join("\n") unless errors.empty?
puts "#{actual_tests.size} test paths and #{ids.size} domain IDs mapped exactly; #{commands.size} command results accounted for"
states = rows.group_by { |line| line.split('|')[4].strip }.transform_values(&:size)
puts "Domain acceptance remains explicit: #{states}"
counts = command_section.lines.grep(/^\| [a-f0-9]{12} \|/).group_by { |line| line.split('|')[5].strip }.transform_values(&:size)
puts "Command snapshots (not a claim that all passed): #{counts}"
