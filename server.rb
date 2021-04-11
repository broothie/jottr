require 'dotenv/load'
require 'sinatra'
require 'sinatra/cookies'
require 'sinatra/reloader' if development?
require 'google/cloud/firestore'
require 'date'
require 'json'

not_found do
  redirect '/home'
end

# Ping
get '/ping' do
  'pong'
end

# About page
get '/home' do
  @jots = get_recent_jots
  erb :home
end

# SCSS
get '/style.css' do
  scss :style
end

# Root. Automatically redirects to brand new jot
get '/' do
  jot_id = new_jot_code
  now = Time.now
  jots.doc(jot_id).set(id: jot_id, created_at: now, updated_at: now)

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

delete '/api/jots/:jot_id' do |jot_id|
  jots.doc(jot_id).delete
  redirect '/home'
end

# Update jot in db
put '/api/jots/:jot_id' do |jot_id|
  payload = JSON.parse(request.body.read)
  jots.doc(jot_id).update(title: payload['title'], body: payload['body'], updated_at: Time.now)
end

helpers do
  ALPHABET = ('a'..'z').to_a.freeze
  DIGITS = (0..9).to_a.freeze
  ENCODED_EMPTY_ARRAY = Base64.urlsafe_encode64([].to_json).freeze

  def new_jot_code
    "#{Array.new(3) { ALPHABET.sample }.join}-#{Array.new(3) { DIGITS.sample }.join}"
  end

  def set_recent_jot!(jot_id)
    jot_ids = Set.new(get_recent_jot_ids)
    jot_ids << jot_id
    cookies[:jot_ids] = cookie_encode(jot_ids.to_a)
  end

  def get_recent_jots
    return [] if get_recent_jot_ids.empty?

    jots.where(:id, :in, get_recent_jot_ids).get
  end

  def get_recent_jot_ids
    cookie_decode(cookies[:jot_ids] || ENCODED_EMPTY_ARRAY)
  end

  def cookie_encode(value)
    Base64.urlsafe_encode64(value.to_json)
  end

  def cookie_decode(raw)
    JSON.parse(Base64.urlsafe_decode64(raw))
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
