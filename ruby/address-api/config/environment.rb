# Load the Rails application.
require 'opentelemetry/sdk'
require_relative "application"

OpenTelemetry::SDK.configure do |c|
  c.use 'OpenTelemetry::Instrumentation::HttpClient'
  c.use 'OpenTelemetry::Instrumentation::Net::HTTP'
  c.use 'OpenTelemetry::Instrumentation::PG', {
    # By default, this instrumentation includes the executed SQL as the `db.statement`
    # semantic attribute. Optionally, you may disable the inclusion of this attribute entirely by
    # setting this option to :omit or sanitize the attribute by setting to :obfuscate
    db_statement: :obfuscate,
  }
  c.use 'OpenTelemetry::Instrumentation::Rails'
end

# Initialize the Rails application.
Rails.application.initialize!
