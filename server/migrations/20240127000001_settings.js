exports.up = function(knex) {
  return knex.schema
    .createTable('hmimg_settings', function(table) {
      table.string('key').primary();
      table.string('value').notNullable();
      table.timestamps(true, true);
    })
    .then(() => {
      // Seed default settings
      return knex('hmimg_settings').insert([
        { key: 'max_users', value: '100' },
        { key: 'allow_registration', value: 'true' }
      ]);
    });
};

exports.down = function(knex) {
  return knex.schema.dropTableIfExists('hmimg_settings');
};
