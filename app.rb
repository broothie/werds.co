require 'sinatra'
require 'base64'

not_found { redirect to '/' }

get '/' do
  if text = params['t']
    redirect to "/#{Base64.urlsafe_encode64(text)}"
  else
    erb :index
  end
end

get '/:hash' do |hash|
  begin
    @text = Base64.urlsafe_decode64(hash)
  rescue
    redirect to '/'
  end

  erb :main
end
