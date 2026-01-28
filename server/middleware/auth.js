const jwt = require('jsonwebtoken');
const config = require('../config/config');

module.exports = async (ctx, next) => {
  const token = ctx.header.authorization?.replace('Bearer ', '');
  
  if (!token) {
    ctx.status = 401;
    ctx.body = { error: 'Authentication required' };
    return;
  }

  try {
    const decoded = jwt.verify(token, config.jwtSecret);
    ctx.state.user = decoded;
    await next();
  } catch (err) {
    ctx.status = 401;
    ctx.body = { error: 'Invalid token' };
  }
};
