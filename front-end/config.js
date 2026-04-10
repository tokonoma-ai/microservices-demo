(function (){
  'use strict';

  var session      = require("express-session"),
      RedisStore   = require('connect-redis')(session)

  var secret = process.env.SESSION_SECRET;
  if (!secret) {
    if (process.env.NODE_ENV === 'production') {
      throw new Error('SESSION_SECRET environment variable must be set in production');
    }
    secret = 'dev-only-insecure-secret';
  }

  module.exports = {
    session: {
      name: 'md.sid',
      secret: secret,
      resave: false,
      saveUninitialized: true
    },

    session_redis: {
      store: new RedisStore({host: "session-db"}),
      name: 'md.sid',
      secret: secret,
      resave: false,
      saveUninitialized: true
    }
  };
}());
