const db = require('../db');
const bcrypt = require('bcrypt');
const jwt = require('jsonwebtoken');
const config = require('../config/config');

exports.login = async (ctx) => {
  const { username, password } = ctx.request.body;
  const user = await db('hmimg_users').where({ username }).first();

  if (!user || !(await bcrypt.compare(password, user.password))) {
    ctx.status = 401;
    ctx.body = { error: 'Invalid credentials' };
    return;
  }

  const token = jwt.sign(
    { id: user.id, username: user.username, role: user.role },
    config.jwtSecret,
    { expiresIn: '24h' }
  );

  ctx.body = { token, user: { id: user.id, username: user.username, role: user.role } };
};

exports.register = async (ctx) => {
  // Fetch settings from DB
  const settings = await db('hmimg_settings').select('*');
  const settingsMap = settings.reduce((acc, curr) => ({ ...acc, [curr.key]: curr.value }), {});
  
  const allowRegistration = settingsMap.allow_registration === 'true';
  const maxUsers = parseInt(settingsMap.max_users || '100', 10);

  if (!allowRegistration) {
    ctx.status = 403;
    ctx.body = { error: 'Registration is closed' };
    return;
  }

  const countResult = await db('hmimg_users').count('id as count').first();
  const userCount = countResult.count;

  if (userCount >= maxUsers) {
    ctx.status = 403;
    ctx.body = { error: 'User limit reached' };
    return;
  }

  const { username, password, email, phone } = ctx.request.body;
  
  // Check if exists
  const existing = await db('hmimg_users').where({ username }).first();
  if (existing) {
      ctx.status = 400;
      ctx.body = { error: 'Username taken' };
      return;
  }

  const hashedPassword = await bcrypt.hash(password, 10);
  
  await db('hmimg_users').insert({
    username,
    password: hashedPassword,
    email,
    phone,
    role: 'user'
  });

  ctx.body = { message: 'Registered successfully' };
};

exports.updateAdmin = async (ctx) => {
    const { username, password, oldPassword, email, phone } = ctx.request.body;
    const userId = ctx.state.user.id;

    // Check if user exists
    const user = await db('hmimg_users').where('id', userId).first();
    if (!user) {
        ctx.status = 404;
        ctx.body = { error: 'User not found' };
        return;
    }

    const updates = {};
    if (username) updates.username = username;
    if (email) updates.email = email;
    if (phone) updates.phone = phone;
    
    // Handle password change
    if (password) {
        // Require old password
        if (!oldPassword) {
            ctx.status = 400;
            ctx.body = { error: 'Old password is required to set a new password' };
            return;
        }

        // Verify old password
        const valid = await bcrypt.compare(oldPassword, user.password);
        if (!valid) {
            ctx.status = 401;
            ctx.body = { error: 'Invalid old password' };
            return;
        }

        updates.password = await bcrypt.hash(password, 10);
    }

    await db('hmimg_users').where('id', userId).update(updates);
    
    // If password changed, maybe invalidate tokens? For now just return success
    // Frontend should handle logout
    
    ctx.body = { message: 'Profile updated successfully', passwordChanged: !!password };
};
