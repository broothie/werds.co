require 'sinatra'
require 'redis'
require 'securerandom'

set :key_size, 8
set :redis, Redis.new(url: ENV['REDIS_URL'])

not_found { redirect to '/' }

get '/' do
  if text = params['text']
    key = SecureRandom.urlsafe_base64(settings.key_size)
    settings.redis.set(key, text)
    redirect to "/#{key}"
  else
    erb :index
  end
end

get '/:key' do |key|
  begin
    @text = settings.redis.get(key).gsub("\n", '<br/>')
  rescue
    redirect to '/'
  end

  erb :main
end
