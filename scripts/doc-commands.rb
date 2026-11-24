#!/usr/bin/env ruby
# Print normalized module validation commands. This inventory never executes them.
require 'json'
require 'shellwords'
require 'pathname'
require 'digest'

# Documents are UTF-8; the C locale would raise "invalid byte sequence in US-ASCII".
Encoding.default_external = Encoding::UTF_8

# Commands that need an external environment. The verification task skips them by
# default and prints the reason instead of reporting a false failure.
DOCKER_DEPENDENT_TEST_MODULES = %w[common/docker].freeze

def requirement(cwd, command)
  return 'docker' if DOCKER_DEPENDENT_TEST_MODULES.include?(cwd) && command.match?(/\Ago test\b/)
  return 'docker' if cwd == 'beacon' && (command.match?(/\Adocker buildx\b/) || command.match?(/\Atask (?:amd64|arm64)\b/))
  return 'service' if command.match?(/\Atask dev\b/) || command.match?(/\Acurl\b/)
  'self'
end

def module_commands
  commands = []
  Dir.glob('.docs/modules/*.md').sort.each do |doc|
    section = File.read(doc)[/^## 验证命令\n(.*?)(?=^## |\z)/m, 1]
    next unless section
    cwd = '.'
    section.scan(/```bash\n(.*?)```/m).flatten.each do |block|
      block.gsub(/\\\n\s*/, ' ').each_line do |raw|
        line = raw.sub(/\s+#.*$/, '').strip
        next if line.empty? || line.start_with?('#')
        # Each explicit `cd ... && ...` example starts from the workspace root.
        cwd = '.' if line.include?('&&')
        line.split('&&').each do |part|
          part = part.strip
          if part.start_with?('cd ')
            cwd = Pathname.new(File.join(cwd, Shellwords.split(part)[1])).cleanpath.to_s
            raise "#{doc}: missing directory #{cwd}" unless File.directory?(cwd)
          else
            # Package scripts are validated by doc:check's static script check.
            if part.match?(/\A(?:pnpm|npm|yarn) /)
              words = Shellwords.split(part)
              script = words[1] == 'run' ? words[2] : words[1]
              package = JSON.parse(File.read(File.join(cwd, 'package.json')))
              unless script == 'install' || package.fetch('scripts', {}).key?(script)
                raise "#{doc}: missing package script #{cwd}:#{script}"
              end
              next
            end
            next if part == 'task --list'
            key = Digest::SHA256.hexdigest([doc, cwd, part].join("\n"))[0, 12]
            commands << { 'id' => key, 'doc' => doc, 'cwd' => cwd, 'command' => part, 'requires' => requirement(cwd, part) }
          end
        end
      end
    end
  end
  commands
end
puts JSON.pretty_generate(module_commands) if $PROGRAM_NAME == __FILE__
