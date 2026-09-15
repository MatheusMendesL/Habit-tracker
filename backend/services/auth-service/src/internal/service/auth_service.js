const jwt = require("jsonwebtoken");
const redis = require("../../config/redis");
const userRepository = require("../repository/user_repository");
const { hashPass, comparePass } = require("../../utils/functions");
const {
  createError,
  ERR_INVALID_CREDENTIALS,
  ERR_EMAIL_ALREADY_EXISTS,
  ERR_INVALID_TOKEN,
  ERR_USER_NOT_FOUND,
} = require("../errors/errors");

async function persistRefreshToken(userId, refreshToken) {
  if (!userId || !refreshToken) {
    throw createError("User id and refresh token are required", "INVALID_ARGUMENT", 400);
  }

  await redis.set(`refresh:user:${userId}`, refreshToken, {
    EX: 60 * 60 * 24 * 7,
  });
}

async function revokeRefreshToken(userId) {
  if (!userId) {
    throw createError("User id is required", "INVALID_ARGUMENT", 400);
  }

  await redis.del(`refresh:user:${userId}`);
}

function generateTokens(userId) {
  if (!userId) {
    throw createError("User id is required", "INVALID_ARGUMENT", 400);
  }

  const accessToken = jwt.sign({ id: userId }, process.env.JWT_SECRET, { expiresIn: "30m" });
  const refreshToken = jwt.sign({ id: userId }, process.env.JWT_REFRESH_SECRET, { expiresIn: "7d" });

  return { accessToken, refreshToken };
}

async function login({ email, password }) {
  if (!email || !password) {
    throw createError("Email and password are required", "INVALID_ARGUMENT", 400);
  }

  const userResult = await userRepository.findByEmail(email);
  if (!userResult.data.length) {
    throw ERR_USER_NOT_FOUND;
  }

  const user = userResult.data[0];
  const isValid = await comparePass(password, user.password);
  if (!isValid) {
    throw ERR_INVALID_CREDENTIALS;
  }

  const tokens = generateTokens(user.id);
  await persistRefreshToken(user.id, tokens.refreshToken);

  return {
    user,
    ...tokens,
  };
}

async function signup({ name, email, tel, password }) {
  if (!name || !email || !password) {
    throw createError("Name, email and password are required", "INVALID_ARGUMENT", 400);
  }

  const existingUser = await userRepository.findByEmail(email);
  if (existingUser.data.length > 0) {
    throw ERR_EMAIL_ALREADY_EXISTS;
  }

  const hashedPassword = await hashPass(password);
  const createdUser = await userRepository.createUser({
    name,
    email,
    tel,
    password: hashedPassword,
  });

  const tokens = generateTokens(createdUser.insertId);
  await persistRefreshToken(createdUser.insertId, tokens.refreshToken);

  return {
    user: createdUser.data,
    ...tokens,
  };
}

async function refreshAccessToken(refreshToken) {
  if (!refreshToken) {
    throw ERR_INVALID_TOKEN;
  }

  try {
    const payload = jwt.verify(refreshToken, process.env.JWT_REFRESH_SECRET);
    const storedRefreshToken = await redis.get(`refresh:user:${payload.id}`);

    if (!storedRefreshToken || storedRefreshToken !== refreshToken) {
      throw ERR_INVALID_TOKEN;
    }

    const accessToken = jwt.sign({ id: payload.id }, process.env.JWT_SECRET, { expiresIn: "30m" });
    return { accessToken };
  } catch (error) {
    throw ERR_INVALID_TOKEN;
  }
}

async function getRedisKeys() {
  const keys = await redis.keys("*");
  const data = {};

  for (const key of keys) {
    data[key] = await redis.get(key);
  }

  return { data, keys };
}

module.exports = {
  persistRefreshToken,
  revokeRefreshToken,
  generateTokens,
  login,
  signup,
  refreshAccessToken,
  getRedisKeys,
};
