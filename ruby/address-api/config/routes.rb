Rails.application.routes.draw do
  Rails.application.routes.draw do
    get 'address/:zip_code', to: 'address#address_by_zipcode'
  end
end
