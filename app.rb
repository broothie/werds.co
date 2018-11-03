require 'sinatra'
require 'base64'

not_found { redirect '/' }

get '/' do
  if params['t']
    redirect "/#{Base64.urlsafe_encode64(params['t'])}"
    return
  end

  erb :index
end

get '/:hash' do |hash|
  erb :main, locals: { text: Base64.urlsafe_decode64(hash) }
end
