require 'opentelemetry'
require "uri"
require "net/http"

class AddressController < ApplicationController
  def address_by_zipcode
    tracer = OpenTelemetry.tracer_provider.tracer('RubyApp', '1.0.0')
    tracer.in_span('get_address_via_cep', kind: :server) do |span|

      url = URI("https://viacep.com.br/ws/#{params[:zip_code]}/json/")
      https = Net::HTTP.new(url.host, url.port)
      https.use_ssl = true
      request = Net::HTTP::Get.new(url)
      response = https.request(request)

      render json: response.read_body.to_s
    end
  end
end
