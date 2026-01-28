module.exports = {
  // Server Configuration
  port: process.env.PORT || 3001,
  jwtSecret: process.env.JWT_SECRET || 'YOUR_JWT_SECRET_HERE',
  
  // Database Configuration
  database: {
    client: 'mysql2', // 'client' can be 'mysql2' or 'pg'
    connection: {
      host: 'YOUR_DB_HOST', // e.g., '127.0.0.1'
      port: 3306, // Default MySQL port 3306. For PostgreSQL use 5432
      user: 'YOUR_DB_USER',
      password: 'YOUR_DB_PASSWORD',
      database: 'hmimg_db'
    },
    pool: {
      min: 2,
      max: 10
    },
    migrations: {
      tableName: 'knex_migrations'
    }
  },

  // Upload Configuration
  uploadDir: 'uploads',
  
  // User Limits
  maxUsers: 100, // Maximum number of users allowed
  allowRegistration: true // Whether admin allows registration
};
