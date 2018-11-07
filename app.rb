require 'sinatra'
require 'redis'
require 'securerandom'

set :key_size, ENV['KEY_SIZE'] || 8
set :redis, Redis.new(url: ENV['REDIS_URL'])

not_found { redirect to '/' }

get '/' do
  text = params['t']

  if text
    key = SecureRandom.urlsafe_base64(settings.key_size)
    settings.redis.set(key, text)
    redirect to "/#{key}"
  else
    erb :index
  end
end

get '/:key' do |key|
  text = settings.redis.get(key)

  if text
    @text = text.gsub("\n", '<br/>')
    erb :main
  else
    redirect to '/'
  end
end
