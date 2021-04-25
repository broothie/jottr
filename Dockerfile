# Frontend
FROM node:15.14.0 AS frontend

WORKDIR /usr/src/app
COPY package.json yarn.lock ./
COPY frontend frontend

RUN yarn
RUN yarn build

# Backend
FROM ruby:2.7.1

WORKDIR /usr/src/app
COPY Gemfile Gemfile.lock puma.rb config.ru server.rb ./
COPY --from=frontend /usr/src/app/public public

RUN gem install bundler -v 2.1.4
RUN bundle config set without development
RUN bundle

CMD ["bundle", "exec", "puma", "-C", "puma.rb"]
