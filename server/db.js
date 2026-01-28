const knex = require('knex');
const config = require('./knexfile');

// Select environment, default to development
const environment = process.env.NODE_ENV || 'development';
const db = knex(config[environment]);

module.exports = db;
