const db = require('../db');

exports.getPublicSettings = async (ctx) => {
  const settings = await db('hmimg_settings').where('key', 'allow_registration').first();
  ctx.body = {
    allow_registration: settings ? settings.value === 'true' : false
  };
};

exports.getSettings = async (ctx) => {
  if (ctx.state.user.role !== 'admin') {
    return ctx.throw(403, 'Access denied');
  }
  const settings = await db('hmimg_settings').select('*');
  const settingsMap = {};
  settings.forEach(s => {
    settingsMap[s.key] = s.value;
  });
  ctx.body = settingsMap;
};

exports.updateSettings = async (ctx) => {
  if (ctx.state.user.role !== 'admin') {
    return ctx.throw(403, 'Access denied');
  }
  
  const { max_users, allow_registration } = ctx.request.body;
  
  if (max_users !== undefined) {
      await db('hmimg_settings').insert({ key: 'max_users', value: String(max_users) })
        .onConflict('key').merge();
  }
  
  if (allow_registration !== undefined) {
      await db('hmimg_settings').insert({ key: 'allow_registration', value: String(allow_registration) })
        .onConflict('key').merge();
  }

  const { website_title } = ctx.request.body;
  if (website_title !== undefined) {
      await db('hmimg_settings').insert({ key: 'website_title', value: String(website_title) })
        .onConflict('key').merge();
  }
  
  ctx.body = { message: 'Settings updated' };
};
