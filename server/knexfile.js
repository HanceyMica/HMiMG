const config = require('./config/config');

module.exports = {
  development: config.database,
  production: config.database
};
