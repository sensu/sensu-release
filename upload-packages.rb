#!/usr/bin/env ruby
gem "aws-sdk-s3", "~> 1"
gem "mixlib-cli", "~> 1.7"
gem "packagecloud-ruby", "1.0.8"

require "mixlib/cli"
require "aws-sdk-s3"
require "packagecloud"
require "json"


# Deep merge two hashes. Nested hashes are deep merged, arrays are
# concatenated and duplicate array items are removed.
#
# @param hash_one [Hash]
# @param hash_two [Hash]
# @return [Hash] deep merged hash.
def deep_merge(hash_one, hash_two)
  merged = hash_one.dup
  hash_two.each do |key, value|
    merged[key] = case
    when hash_one[key].is_a?(Hash) && value.is_a?(Hash)
      deep_merge(hash_one[key], value)
    when hash_one[key].is_a?(Array) && value.is_a?(Array)
      (hash_one[key] + value).uniq
    else
      value
    end
  end
  merged
end

REQUIRED_ENV_VARS = [
  "AWS_ACCESS_KEY_ID",
  "AWS_SECRET_ACCESS_KEY",
  "PACKAGECLOUD_USER",
  "PACKAGECLOUD_TOKEN"
]

PACKAGE_TYPES = [
  "deb",
  "rpm"
]

# every deb & rpm distro/version supported by packagecloud.io (last updated 2018-10-24)
PACKAGECLOUD_PLATFORMS = {
  :linux => {
    :deb => {
      :sysvinit => {
        :ubuntu => [
          # "warty",
          # "hoary",
          # "breezy",
          # "dapper",
          # "edgy",
          # "feisty",
          # "gutsy",
          # "hardy",
          # "intrepid",
          # "jaunty",
          # "karmic",
          # "lucid",
          # "maverick",
          # "natty",
          # "oneiric",
          "precise",
          "quantal",
          "raring",
          "saucy",
          "trusty",
          "utopic"
        ],
        :debian => [
          "etch",
          "lenny",
          "squeeze",
          "wheezy"
        ],
        :raspbian => [
          # "wheezy"
        ],
        :elementaryos => [
          # "jupiter",
          # "luna",
          # "freya"
        ],
        :linuxmint => [
          # "petra",
          # "qiana",
          # "rebecca",
          # "rafaela",
          # "rosa"
        ]
      },
      :systemd => {
        :ubuntu => [
          "vivid",
          "wily",
          "xenial",
          "yakkety",
          "zesty",
          "artful",
          "bionic"
        ],
        :debian => [
          "jessie",
          "stretch",
          "buster"
        ],
        :raspbian => [
          # "jessie",
          # "stretch",
          # "buster"
        ],
        :elementaryos => [
          # "loki"
        ],
        :linuxmint => [
          # "sarah",
          # "serena",
          # "sonya",
          # "sylvia",
          # "tara"
        ]
      }
    },
    :rpm => {
      :sysvinit => {
        :el => [
          "5",
          "6"
        ],
        :ol => [
          "5",
          "6"
        ],
        :scientific => [
          "5",
          "6"
        ],
        :sles => [
          "11.4"
        ],
        :fedora => [
          "14"
        ],
        :poky => [
          # "jethro",
          # "krogoth"
        ]
      },
      :systemd => {
        :el => [
          "7"
        ],
        :ol => [
          "7"
        ],
        :scientific => [
          "7"
        ],
        :sles => [
          "12.0",
          "12.1",
          "12.2",
          "12.3",
          "15.0"
        ],
        :opensuse => [
          "13.1",
          "13.2",
          "42.1",
          "42.2",
          "42.3"
        ],
        :fedora => [
          # "15",
          # "16",
          # "17",
          # "18",
          # "19",
          # "20",
          # "21",
          # "22",
          # "23",
          "24",
          "25",
          "26",
          "27",
          "28"
        ]
      }
    }
  }
}

PACKAGECLOUD_PLATFORMS[:linux][:deb][:none] = deep_merge(PACKAGECLOUD_PLATFORMS[:linux][:deb][:sysvinit], PACKAGECLOUD_PLATFORMS[:linux][:deb][:systemd])
PACKAGECLOUD_PLATFORMS[:linux][:rpm][:none] = deep_merge(PACKAGECLOUD_PLATFORMS[:linux][:rpm][:sysvinit], PACKAGECLOUD_PLATFORMS[:linux][:rpm][:systemd])

class SensuPackageCLI
  include Mixlib::CLI

  option :s3_bucket,
    :long => "--s3-bucket BUCKET",
    :required => true,
    :description => "Amazon S3 bucket to pull artifacts from (e.g. sensu-ci-builds)"

  option :s3_prefix,
    :long => "--s3-prefix PREFIX",
    :required => true,
    :description => "Amazon S3 bucket prefix to deploy (e.g. HEAD/20181019-0935)"

  option :pc_repo,
    :long => "--pc-repo REPOSITORY",
    :required => true,
    :description => "Packagecloud repository to use (e.g. beta)"

  option :help,
    :short => "-h",
    :long => "--help",
    :description => "Show this message",
    :on => :tail,
    :boolean => true,
    :show_options => true,
    :exit => 0
end

def check_env_vars
  missing_env_vars = []

  REQUIRED_ENV_VARS.each do |env_var|
    unless ENV.has_key?(env_var)
      missing_env_vars << env_var
    end
  end

  unless missing_env_vars.empty?
    puts "Missing required environment variables: #{missing_env_vars.join(', ')}"
    exit 1
  end
end

cli = SensuPackageCLI.new
cli.parse_options
@s3_bucket = cli.config[:s3_bucket]
@s3_prefix = cli.config[:s3_prefix]
@pc_repo = cli.config[:pc_repo]

# verify that the required environment variables are set
check_env_vars

# connect to Amazon S3
@s3 ||= Aws::S3::Client.new()

# connect to packagecloud
credentials = Packagecloud::Credentials.new(ENV["PACKAGECLOUD_USER"], ENV["PACKAGECLOUD_TOKEN"])
@pc = Packagecloud::Client.new(credentials)

# fetch all object paths (rpms, debs) from s3
puts "fetching object paths from s3 ..."
object_paths = []
PACKAGE_TYPES.each do |package_type|
  @s3.list_objects(:bucket => @s3_bucket, :prefix => "#{@s3_prefix}/#{package_type}").contents.each do |object|
    object_paths << object.key
  end
end

if object_paths.empty?
  puts "error: no objects found in s3"
  exit 1
end

Dir.mktmpdir("sensu-artifacts-") do |tmpdir|
  object_paths.each do |object_path|
    destination_path = File.join(tmpdir, object_path)
    destination_dir = File.dirname(destination_path)
    object_dir = File.dirname(object_path)
    package_type = object_dir.split("/")[2].to_sym
    service_type = object_dir.split("/")[3].to_sym

    # create destination directory if it doesn't exist
    unless Dir.exist?(destination_dir)
      FileUtils.mkdir_p(destination_dir)
    end

    # download artifact from s3
    puts "downloading #{object_path} ..."
    @s3.get_object(:bucket => @s3_bucket, :response_target => destination_path, :key => object_path)

    PACKAGECLOUD_PLATFORMS[:linux][package_type][service_type].each do |distro, versions|
      versions.each do |version|
        begin
          puts "uploading #{object_path} to packagecloud #{distro}/#{version} ..."
          pkg = Packagecloud::Package.new(:file => destination_path)
          response = @pc.put_package(@pc_repo, pkg, "#{distro}/#{version}")
          unless response.succeeded
            messages = JSON.parse(response.response)
            messages.each do |message_type, sub_messages|
              sub_messages.each do |sub_message|
                puts "  error: #{message_type} #{sub_message}"
              end
            end
          end
        rescue => error
          puts "  error: #{error.message}"
        end
      end
    end
  end
end
