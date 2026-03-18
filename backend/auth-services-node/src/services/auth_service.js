const redis = require("../config/redis")

async function login(user) {
    await redis.set(`refresh:user:${user.id}`, user.refreshToken, {
        EX: 60 * 60 * 24 * 7
    });
}

async function signup(user) {
    await redis.set(`refresh:user:${user.id}`, user.refreshToken, {
        EX: 60 * 60 * 24 * 7
    });
}

async function logout(userId) {
    await redis.del(`refresh:user:${userId}`);
}

async function get_keys() {
    const keys = await redis.keys("*");

    const data = {};

    for (const key of keys) {
        data[key] = await redis.get(key);
    }

    return { data, keys }
}

module.exports = {
    login,
    logout,
    get_keys,
    signup
}