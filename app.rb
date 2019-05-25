require 'sinatra'
require 'json'

post '/' do
  File.open('data.json', 'w') { |f| f.write(request.body.read) }
end

