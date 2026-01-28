const Joi = require('joi');

const schemas = {
  register: Joi.object({
    username: Joi.string().alphanum().min(3).max(30).required(),
    password: Joi.string().min(6).required(),
    email: Joi.string().email().required(),
    phone: Joi.string().pattern(/^[0-9+]+$/).required(),
    // Allow other fields if necessary or strip them
  }),
  login: Joi.object({
    username: Joi.string().required(),
    password: Joi.string().required()
  }),
  updateAdmin: Joi.object({
    username: Joi.string().alphanum().min(3).max(30),
    oldPassword: Joi.string().allow(''), // Optional if not changing password, but needed if password is set
    password: Joi.string().min(6).allow(''), // New password
    email: Joi.string().email(),
    phone: Joi.string().pattern(/^[0-9+]+$/)
  }),
  createAlbum: Joi.object({
    name: Joi.string().required(),
    description: Joi.string().allow('', null)
  }),
  updateAlbum: Joi.object({
    name: Joi.string(),
    description: Joi.string().allow('', null)
  }),
  createCollection: Joi.object({
    name: Joi.string().required(),
    description: Joi.string().allow('', null)
  }),
  updateCollection: Joi.object({
    name: Joi.string(),
    description: Joi.string().allow('', null)
  }),
  addToCollection: Joi.object({
    collectionId: Joi.number().required(),
    itemType: Joi.string().valid('album', 'collection').required(),
    itemName: Joi.string().required()
  })
};

const validate = (schema) => {
  return async (ctx, next) => {
    const { error, value } = schema.validate(ctx.request.body);
    if (error) {
      ctx.status = 400;
      ctx.body = { error: error.details[0].message };
      return;
    }
    // Replace body with validated value (useful for type conversion or stripping unknowns)
    ctx.request.body = value;
    await next();
  };
};

module.exports = {
  schemas,
  validate
};
