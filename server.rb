require 'dotenv/load'
require 'sinatra'
require 'sinatra/cookies'
require 'sinatra/reloader' if development?
require 'google/cloud/firestore'
require 'date'
require 'json'

set session_secret: ENV.fetch('SESSION_SECRET')
enable :sessions

# Ping
get '/ping' do
  'pong'
end

# About page
get '/home' do
  jot_ids = get_recent_jot_ids
  @jots = jot_ids.empty? ? [] : jots.where(:id, 'in', jot_ids).get
  erb :home
end

# SCSS
get '/style.css' do
  scss :style
end

# Root. Automatically redirects to brand new jot
get '/' do
  jot_id = random_string
  now = Time.now
  jots.doc(jot_id).set(id: jot_id, body: '', created_at: now, updated_at: now)

  redirect "/jots/#{jot_id}"
end

# Serve up a jot
get '/jots/:jot_id' do |jot_id|
  jot_doc = jots.doc(jot_id).get

  unless jot_doc.exists?
    @jot_id = jot_id
    halt erb :not_found unless jot_doc.exists?
  end

  set_recent_jot!(jot_id)
  @jot = jot_doc.data
  erb :jot
end

# Update jot in db
put '/api/jots/:jot_id' do |jot_id|
  payload = JSON.parse(request.body.read)
  jots.doc(jot_id).update(body: payload['body'], updated_at: Time.now)
end

helpers do
  ALPHABET = ('a'..'z').to_a.freeze

  def random_string(length = 10)
    Array.new(length) { ALPHABET.sample }.join
  end

  def set_recent_jot!(jot_id)
    jot_ids = Set.new(get_recent_jot_ids)
    jot_ids << jot_id
    set_recent_jot_cookie!(jot_ids.to_a)
  end

  def get_recent_jot_ids
    (cookies[:jot_ids] || '').split('|')
  end

  def set_recent_jot_cookie!(jot_ids)
    cookies[:jot_ids] = jot_ids.join('|')
  end

  def jots
    @jots ||= firestore.collection("#{collection_prefix}.jots")
  end

  def firestore
    @firestore ||= Google::Cloud::Firestore.new
  end

  def collection_prefix
    @collection_prefix ||= settings.production? ? 'production' : "development.#{`whoami`.chomp}"
  end
end
