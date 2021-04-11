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
  @title = 'jottr - home'
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
  jots.doc(jot_id).set(id: jot_id, read_only_id: new_jot_code, created_at: now, updated_at: now)

  redirect "/jots/#{jot_id}"
end

# Serve up a jot
get '/jots/:jot_id' do |jot_id|
  jot_doc = jots.doc(jot_id).get

  unless jot_doc.exists?
    @jot_id = jot_id
    halt erb :not_found
  end

  set_recent_jot!(jot_id)
  @jot = jot_doc.data
  @title = "jottr - #{@jot[:title]}"
  erb :jot
end

# Serve up a readonly jot
get '/jots/:read_only_jot_id/readonly' do |read_only_jot_id|
  jot_docs = jots.where(:read_only_id, :==, read_only_jot_id).get

  if jot_docs.count.zero?
    @jot_id = read_only_jot_id
    halt erb :not_found
  end

  @jot = jot_docs.first.data
  @title = "jottr - #{@jot[:title]}"
  @read_only = true
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
  ENCODED_EMPTY_ARRAY = Base64.urlsafe_encode64([].to_json).freeze

  def new_jot_code
    "#{random_string(3)}-#{random_string(4)}-#{random_string(3)}"
  end

  def random_string(length = 3)
    Array.new(length) { ALPHABET.sample }.join
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
