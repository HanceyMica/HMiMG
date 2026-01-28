const Koa = require('koa');
const Router = require('koa-router');
const bodyParser = require('koa-bodyparser');
const serve = require('koa-static');
const cors = require('koa-cors');
const path = require('path');
const bcrypt = require('bcrypt');
const db = require('./db');
const config = require('./config/config');

const app = new Koa();
const router = new Router();

app.use(cors());
app.use(async (ctx, next) => {
  try {
    await next();
  } catch (err) {
    ctx.status = err.status || 500;
    ctx.body = { error: err.message };
    // Only log 500 errors to console
    if (ctx.status === 500) {
      console.error('Server Error:', err);
    }
  }
});
app.use(bodyParser());
app.use(serve(path.join(__dirname, config.uploadDir)));

// Check DB Connection and Init Admin
async function bootstrap() {
  try {
    await db.raw('SELECT 1');
    console.log('Database connected successfully.');

    // Run migrations automatically on start if needed, or user runs them manually.
    // Let's run them to ensure tables exist.
    try {
        await db.migrate.latest();
        console.log('Migrations run successfully.');
    } catch (e) {
        console.error('Migration failed:', e);
    }

    // Check Admin
    const admin = await db('hmimg_users').where({ username: 'admin' }).first();
    if (!admin) {
      const hashedPassword = await bcrypt.hash('admin', 10);
      await db('hmimg_users').insert({
        username: 'admin',
        password: hashedPassword,
        email: 'admin@yourdomaname.com',
        phone: '+8613200000000',
        role: 'admin'
      });
      console.log('Default admin account created.');
    }

  } catch (error) {
    console.error('Database connection failed:', error);
    // process.exit(1); // Keep running but log error
  }
}

bootstrap();

// Routes
router.get('/', async (ctx) => {
  ctx.body = 'HMiMG API Server';
});

// Import other routes here
const apiRoutes = require('./routes');
app.use(apiRoutes.routes());

app.use(router.routes()).use(router.allowedMethods());

const PORT = config.port;
app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
