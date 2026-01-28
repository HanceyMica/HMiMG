exports.up = function(knex) {
  return knex.schema
    .createTable('hmimg_users', function(table) {
      table.increments('id').primary();
      table.string('username').notNullable().unique();
      table.string('password').notNullable();
      table.string('email').notNullable();
      table.string('phone').notNullable();
      table.string('role').defaultTo('user'); // 'admin', 'user'
      table.timestamps(true, true);
    })
    .createTable('hmimg_albums', function(table) {
      table.increments('id').primary();
      table.string('name').notNullable();
      table.text('description');
      table.integer('created_by').unsigned().references('id').inTable('hmimg_users').onDelete('SET NULL');
      table.string('cover_image').nullable(); // Store path or ID
      table.timestamps(true, true);
    })
    .createTable('hmimg_collections', function(table) {
      table.increments('id').primary();
      table.string('name').notNullable();
      table.text('description');
      table.integer('created_by').unsigned().references('id').inTable('hmimg_users').onDelete('SET NULL');
      // We will handle hierarchy via the join table or parent_id if strict tree.
      // Given "Many-to-Many" for Collection-Album, and "Collection contains Collection", 
      // let's use a join table for maximum flexibility.
      table.timestamps(true, true);
    })
    .createTable('hmimg_images', function(table) {
      table.increments('id').primary();
      table.string('filename').notNullable();
      table.string('original_name').notNullable();
      table.string('path').notNullable();
      table.integer('size').unsigned();
      table.string('mimetype');
      table.integer('album_id').unsigned().notNullable().references('id').inTable('hmimg_albums').onDelete('CASCADE');
      table.integer('uploaded_by').unsigned().references('id').inTable('hmimg_users').onDelete('SET NULL');
      table.timestamps(true, true);
    })
    .createTable('hmimg_collection_items', function(table) {
      table.increments('id').primary();
      table.integer('collection_id').unsigned().notNullable().references('id').inTable('hmimg_collections').onDelete('CASCADE');
      
      // Target can be album or collection
      table.string('item_type').notNullable(); // 'album' or 'collection'
      table.integer('item_id').unsigned().notNullable();
      
      // We can't easily enforce foreign key on item_id because it points to different tables.
      // Logic must handle validity.
      
      table.unique(['collection_id', 'item_type', 'item_id']);
    });
};

exports.down = function(knex) {
  return knex.schema
    .dropTableIfExists('hmimg_collection_items')
    .dropTableIfExists('hmimg_images')
    .dropTableIfExists('hmimg_collections')
    .dropTableIfExists('hmimg_albums')
    .dropTableIfExists('hmimg_users');
};
