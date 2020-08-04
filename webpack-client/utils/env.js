const packageJson = require("../package.json")

const env = {
  NODE_ENV: process.env.NODE_ENV || "development",
  PORT: process.env.PORT || 3000,
  __CWP_APP_VERSION: packageJson.version,
}

if (process.env.NODE_ENV === "production") {
  env.__API_BASE_URL = 'wss://crossword-party.herokuapp.com/ws'
} else {
  env.__API_BASE_URL = 'wss://localhost:8000/ws'
}

Object.keys(env).forEach(key => process.env[key] = env[key])

module.exports = env;
