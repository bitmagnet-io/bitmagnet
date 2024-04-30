require 'fileutils'

module Schemas
  class Generator < Jekyll::Generator
    def generate(site)
      FileUtils.mkdir_p File.join(site.dest, 'schemas')
      Dir.glob(File.join(site.source, 'schemas', '*.json')).each do |json_file|
        File.open(File.join(site.dest, 'schemas', File.basename(json_file)), 'w') do |file|
          file.write(File.read(json_file))
        end
      end
    end
  end
end
