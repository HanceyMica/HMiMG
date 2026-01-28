const db = require('../db');
const fs = require('fs');
const path = require('path');
const config = require('../config/config');

// Albums
exports.createAlbum = async (ctx) => {
  const { name, description } = ctx.request.body;
  const userId = ctx.state.user.id;

  const existing = await db('hmimg_albums').where('name', name).first();
  if (existing) return ctx.throw(400, 'Album name already exists');

  const [id] = await db('hmimg_albums').insert({
    name,
    description,
    created_by: userId
  });
  ctx.body = { id, name, description };
};

exports.getAlbums = async (ctx) => {
  const albums = await db('hmimg_albums').select('*');
  ctx.body = albums;
};

exports.getAlbum = async (ctx) => {
  const id = ctx.params.id;
  const album = await db('hmimg_albums').where('id', id).first();
  if (!album) return ctx.throw(404, 'Album not found');
  ctx.body = album;
};

exports.updateAlbum = async (ctx) => {
  const id = ctx.params.id;
  const { name, description } = ctx.request.body;
  const userId = ctx.state.user.id;

  // Optional: Check ownership or admin role
  const album = await db('hmimg_albums').where('id', id).first();
  if (!album) return ctx.throw(404, 'Album not found');

  if (ctx.state.user.role !== 'admin' && album.created_by !== userId) {
      return ctx.throw(403, 'Access denied');
  }

  if (name !== album.name) {
      const existing = await db('hmimg_albums').where('name', name).first();
      if (existing) return ctx.throw(400, 'Album name already exists');
  }

  await db('hmimg_albums').where('id', id).update({
      name,
      description,
      updated_at: db.fn.now()
  });

  ctx.body = { message: 'Album updated', id };
};

exports.deleteAlbum = async (ctx) => {
  const id = ctx.params.id;
  const userId = ctx.state.user.id;

  const album = await db('hmimg_albums').where('id', id).first();
  if (!album) return ctx.throw(404, 'Album not found');

  if (ctx.state.user.role !== 'admin' && album.created_by !== userId) {
      return ctx.throw(403, 'Access denied');
  }

  // Delete images files first (optional, but good practice to clean up)
  const images = await db('hmimg_images').where('album_id', id);
  images.forEach(img => {
      try {
          const filePath = path.join(__dirname, '..', config.uploadDir, img.path);
          if (fs.existsSync(filePath)) {
              fs.unlinkSync(filePath);
          }
      } catch (e) {
          console.error('Failed to delete file', e);
      }
  });

  // Delete DB records
  await db('hmimg_images').where('album_id', id).del();
  await db('hmimg_collection_items').where({ item_type: 'album', item_id: id }).del();
  await db('hmimg_albums').where('id', id).del();

  ctx.body = { message: 'Album deleted' };
};

// Collections
exports.createCollection = async (ctx) => {
  const { name, description } = ctx.request.body;
  const userId = ctx.state.user.id;

  const existing = await db('hmimg_collections').where('name', name).first();
  if (existing) return ctx.throw(400, 'Collection name already exists');

  const [id] = await db('hmimg_collections').insert({
    name,
    description,
    created_by: userId
  });
  ctx.body = { id, name, description };
};

exports.getCollections = async (ctx) => {
  const collections = await db('hmimg_collections').select('*');
  // Enhance with items if needed, or separate API
  ctx.body = collections;
};

exports.getCollection = async (ctx) => {
  const id = ctx.params.id;
  const collection = await db('hmimg_collections').where('id', id).first();
  if (!collection) return ctx.throw(404, 'Collection not found');

  // Get items
  const items = await db('hmimg_collection_items').where('collection_id', id);
  
  let children = [];
  if (items.length > 0) {
      const type = items[0].item_type;
      const ids = items.map(i => i.item_id);
      
      if (type === 'album') {
          children = await db('hmimg_albums').whereIn('id', ids);
          // Attach type for frontend convenience
          children = children.map(c => ({ ...c, type: 'album' }));
      } else if (type === 'collection') {
          children = await db('hmimg_collections').whereIn('id', ids);
          children = children.map(c => ({ ...c, type: 'collection' }));
      }
  }

  ctx.body = { ...collection, children };
};

exports.updateCollection = async (ctx) => {
  const id = ctx.params.id;
  const { name, description } = ctx.request.body;
  const userId = ctx.state.user.id;

  const collection = await db('hmimg_collections').where('id', id).first();
  if (!collection) return ctx.throw(404, 'Collection not found');

  if (ctx.state.user.role !== 'admin' && collection.created_by !== userId) {
      return ctx.throw(403, 'Access denied');
  }

  if (name !== collection.name) {
      const existing = await db('hmimg_collections').where('name', name).first();
      if (existing) return ctx.throw(400, 'Collection name already exists');
  }

  await db('hmimg_collections').where('id', id).update({
      name,
      description,
      updated_at: db.fn.now()
  });

  ctx.body = { message: 'Collection updated', id };
};

exports.deleteCollection = async (ctx) => {
  const id = ctx.params.id;
  const userId = ctx.state.user.id;

  const collection = await db('hmimg_collections').where('id', id).first();
  if (!collection) return ctx.throw(404, 'Collection not found');

  if (ctx.state.user.role !== 'admin' && collection.created_by !== userId) {
      return ctx.throw(403, 'Access denied');
  }

  // Cascading delete
  // 1. Remove this collection from any parent collections
  await db('hmimg_collection_items').where({ item_type: 'collection', item_id: id }).del();
  
  // 2. Remove items inside this collection (the links, not the actual items)
  await db('hmimg_collection_items').where('collection_id', id).del();

  // 3. Delete collection
  await db('hmimg_collections').where('id', id).del();

  ctx.body = { message: 'Collection deleted' };
};

// Add to Collection
exports.addToCollection = async (ctx) => {
  const { collectionId, itemType, itemName } = ctx.request.body; // itemType: 'album' or 'collection', itemName instead of itemId
  
  // 1. Check if collection exists
  const collection = await db('hmimg_collections').where('id', collectionId).first();
  if (!collection) return ctx.throw(404, 'Collection not found');

  // 2. Resolve Item ID from Name
  let itemId;
  if (itemType === 'album') {
      const album = await db('hmimg_albums').where('name', itemName).first();
      if (!album) return ctx.throw(404, 'Album not found');
      itemId = album.id;
  } else if (itemType === 'collection') {
      const col = await db('hmimg_collections').where('name', itemName).first();
      if (!col) return ctx.throw(404, 'Target collection not found');
      if (col.id === parseInt(collectionId)) return ctx.throw(400, 'Cannot add collection to itself');
      itemId = col.id;
  } else {
      return ctx.throw(400, 'Invalid item type');
  }

  // 3. Check mixed content rule
  const existingItems = await db('hmimg_collection_items')
    .where('collection_id', collectionId)
    .first();
  
  if (existingItems && existingItems.item_type !== itemType) {
    return ctx.throw(400, 'Collection cannot contain mixed types');
  }

  // 4. Check if already exists
  const exists = await db('hmimg_collection_items')
      .where({
          collection_id: collectionId,
          item_type: itemType,
          item_id: itemId
      })
      .first();
  
  if (exists) return ctx.throw(400, 'Item already in collection');

  await db('hmimg_collection_items').insert({
    collection_id: collectionId,
    item_type: itemType,
    item_id: itemId
  });

  ctx.body = { message: 'Added successfully' };
};

// Images (Upload)
exports.uploadImages = async (ctx) => {
  const files = ctx.request.files;
  const { albumId } = ctx.request.body;

  if (!files || files.length === 0) return ctx.throw(400, 'No files uploaded');
  if (!albumId) {
      // Cleanup files if no album
      files.forEach(f => fs.unlinkSync(f.path));
      return ctx.throw(400, 'Album ID required');
  }

  const album = await db('hmimg_albums').where('id', albumId).first();
  if (!album) {
      files.forEach(f => fs.unlinkSync(f.path));
      return ctx.throw(404, 'Album not found');
  }

  const insertedIds = [];
  
  for (const file of files) {
      const [id] = await db('hmimg_images').insert({
        filename: file.filename,
        original_name: file.originalname,
        path: file.filename,
        size: file.size,
        mimetype: file.mimetype,
        album_id: albumId,
        uploaded_by: ctx.state.user.id
      });
      insertedIds.push(id);
  }

  // Set cover if album has none (use first image)
  if (!album.cover_image && files.length > 0) {
      await db('hmimg_albums').where('id', albumId).update({ cover_image: files[0].filename });
  }

  ctx.body = { ids: insertedIds, count: files.length };
};

exports.getImages = async (ctx) => {
    const { albumId } = ctx.query;
    let query = db('hmimg_images').select('*');
    if (albumId) {
        query = query.where('album_id', albumId);
    }
    const images = await query;
    ctx.body = images;
};

// Helper function to get all album IDs recursively from a collection
async function getAlbumIdsFromCollection(collectionId) {
    const items = await db('hmimg_collection_items').where('collection_id', collectionId);
    let albumIds = [];

    for (const item of items) {
        if (item.item_type === 'album') {
            albumIds.push(item.item_id);
        } else if (item.item_type === 'collection') {
            const childIds = await getAlbumIdsFromCollection(item.item_id);
            albumIds = albumIds.concat(childIds);
        }
    }
    return albumIds;
}

exports.getRandomImageFromCollection = async (ctx) => {
    const id = ctx.params.id;
    const returnType = ctx.query.type || 'json'; // 'json' or 'redirect'
    
    // Check if collection exists
    const collection = await db('hmimg_collections').where('id', id).first();
    if (!collection) return ctx.throw(404, 'Collection not found');

    // Get all album IDs in this collection (recursive)
    const albumIds = await getAlbumIdsFromCollection(id);

    if (albumIds.length === 0) return ctx.throw(404, 'No albums in this collection');

    // Get a random image from these albums
    // SQLite uses RANDOM(), MySQL/PG uses RAND() or RANDOM(). Knex abstracts this but sometimes needs raw.
    // Since we are using better-sqlite3 (implied by file structure), we use RANDOM()
    const image = await db('hmimg_images')
        .whereIn('album_id', albumIds)
        .orderByRaw('RANDOM()')
        .first();

    if (!image) return ctx.throw(404, 'No images found in this collection');

    if (returnType === 'redirect') {
        ctx.redirect(`http://localhost:3001/${image.path}`);
    } else {
        ctx.body = image;
    }
};

exports.getImage = async (ctx) => {
  const id = ctx.params.id;
  const image = await db('hmimg_images').where('id', id).first();
  if (!image) return ctx.throw(404, 'Image not found');
  
  // Get album info
  const album = await db('hmimg_albums').where('id', image.album_id).first();
  
  ctx.body = { ...image, album_name: album ? album.name : null };
};
