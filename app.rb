require 'sinatra'
require 'base64'

not_found { redirect '/' }

get '/gen/:text' do |text|
  redirect "/#{Base64.urlsafe_encode64(text)}"
end

get '/' do
  erb :index
end

post '/' do
  redirect "/#{Base64.urlsafe_encode64(params['text'])}"
end

get '/:hash' do |hash|
  erb :main, locals: { text: Base64.urlsafe_decode64(hash) }
end
