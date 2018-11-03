require 'sinatra'
require 'base64'

not_found { redirect '/' }

get '/' do
  erb :index
end

post '/' do
  redirect "/#{Base64.encode64(params['text'])}"
end

get '/:hash' do |hash|
  erb :main, locals: { text: Base64.decode64(hash) }
end
