const Router = require('koa-router');
const multer = require('@koa/multer');
const path = require('path');
const config = require('../config/config');
const authController = require('../controllers/auth');
const storageController = require('../controllers/storage');
const settingsController = require('../controllers/settings');
const authMiddleware = require('../middleware/auth');
const { schemas, validate } = require('../middleware/validation');
const fs = require('fs');

const router = new Router({ prefix: '/api' });

// Ensure upload dir exists
const uploadDir = path.join(__dirname, '..', config.uploadDir);
if (!fs.existsSync(uploadDir)) {
    fs.mkdirSync(uploadDir, { recursive: true });
}

const storage = multer.diskStorage({
  destination: function (req, file, cb) {
    cb(null, uploadDir)
  },
  filename: function (req, file, cb) {
    const uniqueSuffix = Date.now() + '-' + Math.round(Math.random() * 1E9)
    cb(null, uniqueSuffix + path.extname(file.originalname))
  }
});

const upload = multer({ 
    storage: storage,
    fileFilter: (req, file, cb) => {
        const allowedMimes = ['image/jpeg', 'image/png', 'image/gif', 'image/webp'];
        if (allowedMimes.includes(file.mimetype)) {
            cb(null, true);
        } else {
            cb(new Error('Invalid file type. Only JPEG, PNG, GIF, and WebP are allowed.'));
        }
    }
});

// Auth Routes
router.post('/login', validate(schemas.login), authController.login);
router.post('/register', validate(schemas.register), authController.register);
router.put('/admin/update', authMiddleware, validate(schemas.updateAdmin), authController.updateAdmin);

// Settings Routes
router.get('/settings/public', settingsController.getPublicSettings);
router.get('/settings', authMiddleware, settingsController.getSettings);
router.put('/settings', authMiddleware, settingsController.updateSettings);

// Storage Routes
router.post('/albums', authMiddleware, validate(schemas.createAlbum), storageController.createAlbum);
router.get('/albums', authMiddleware, storageController.getAlbums);
router.get('/albums/:id', authMiddleware, storageController.getAlbum);
router.put('/albums/:id', authMiddleware, validate(schemas.updateAlbum), storageController.updateAlbum);
router.delete('/albums/:id', authMiddleware, storageController.deleteAlbum);

router.post('/collections', authMiddleware, validate(schemas.createCollection), storageController.createCollection);
router.get('/collections', authMiddleware, storageController.getCollections);
router.get('/collections/:id', authMiddleware, storageController.getCollection);
router.get('/collections/:id/random', authMiddleware, storageController.getRandomImageFromCollection);
router.put('/collections/:id', authMiddleware, validate(schemas.updateCollection), storageController.updateCollection);
router.delete('/collections/:id', authMiddleware, storageController.deleteCollection);
router.post('/collections/add', authMiddleware, validate(schemas.addToCollection), storageController.addToCollection);

router.post('/upload', authMiddleware, upload.array('images', 20), storageController.uploadImages);
router.get('/images', authMiddleware, storageController.getImages);
router.get('/images/:id', authMiddleware, storageController.getImage);

module.exports = router;
