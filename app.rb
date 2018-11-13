require 'sinatra'
require 'redis'
require 'securerandom'

set :key_size, ENV['KEY_SIZE'] || 8
set :text_limit, ENV['TEXT_SIZE'] || 1000
set :redis, Redis.new(url: ENV['REDIS_URL'])

not_found { redirect to '/' }

get '/' do
  text = params['t']

  if text
    key = SecureRandom.urlsafe_base64(settings.key_size)
    settings.redis.set(key, text[0...settings.text_limit])
    redirect to "/#{key}"
  else
    @text_limit = settings.text_limit
    erb :'index.html'
  end
end

get '/:key' do |key|
  text = settings.redis.get(key)

  if text
    @text = CGI.escape_html(text).gsub("\n", '<br/>')
    erb :'main.html'
  else
    redirect to '/'
  end
end
