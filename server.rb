require 'dotenv/load'
require 'sinatra'
require 'sinatra/namespace'
require 'sinatra/json'
require 'sinatra/reloader' if development?
require 'google/cloud/firestore'
require 'date'
require 'json'

not_found do
  send_file 'public/index.html'
end

namespace '/api' do
  get '/ping' do
    'pong'
  end

  namespace '/jots' do
    # Create jot
    post '/?' do
      jot_id = new_jot_code
      now = Time.now
      jot_data = { id: jot_id, read_only_id: new_jot_code, created_at: now, updated_at: now }

      jots.doc(jot_id).set(jot_data)
      status 201
      json jot_data
    end

    namespace '/:jot_id' do
      # Get a jot
      get '/?' do |jot_id|
        json jots.doc(jot_id).get.data
      end

      # Update jot in db
      patch '/?' do |jot_id|
        payload = JSON.parse(request.body.read)
        jots.doc(jot_id).update(title: payload['title'], delta: payload['delta'], updated_at: Time.now)
        status :accepted
      end

      # Delete a jot
      delete '/?' do |jot_id|
        jots.doc(jot_id).delete
        status :accepted
      end
    end
  end

  get '/bulk/jots' do
    jot_ids = params[:jot_ids]
    halt json([]) unless jot_ids

    jot_ids = jot_ids.split(',')
    jot_docs = firestore.transaction do |tx|
      jot_ids.map { |jot_id| tx.get("#{jots_collection}/#{jot_id}") }
    end

    jots = jot_docs
      .map(&:data)
      .compact
      .select { |jot| jot[:title] && !jot[:title].empty? }

    json jots
  end
end

# Purge job
get '/jobs/purge' do
  empty_jots = jots.where(:title, :==, '').get
  firestore.batch { |batch| empty_jots.each(&batch.method(:delete)) }
  status :accepted
end

helpers do
  # Codes

  def new_jot_code
    "#{random_string(3)}-#{random_string(4)}-#{random_string(3)}"
  end

  def random_string(length = 3)
    Array.new(length) { alphabet.sample }.join
  end

  def alphabet
    @alphabet ||= ('a'..'z').to_a.freeze
  end

  # DB

  def jots
    @jots ||= firestore.collection(jots_collection)
  end

  def jots_collection
    @jots_collection ||= "#{collection_prefix}.jots"
  end

  def collection_prefix
    @collection_prefix ||= settings.production? ? 'production' : "development.#{`whoami`.chomp}"
  end

  def firestore
    @firestore ||= Google::Cloud::Firestore.new
  end
end
