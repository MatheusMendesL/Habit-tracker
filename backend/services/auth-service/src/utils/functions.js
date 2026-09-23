const bcrypt = require("bcrypt")

function response(status, message, data = null, httpInfo = null, error = null, startedAt = Date.now()) {
    return {
        timestamp: new Date().toISOString(),
        api_version: "v1",
        duration_ms: startedAt ? Date.now() - startedAt : 0,
        error: error || "",
        message: message || "",
        data: data ?? null,
        status: status || "success",
        http: httpInfo || null,
    }
}

const SALT_ROUNDS = 10

async function hashPass(password) {
    return await bcrypt.hash(password, SALT_ROUNDS);
}

async function comparePass(pass, hash) {
    return await bcrypt.compare(pass, hash);
}

module.exports = {
    response,
    hashPass,
    comparePass
}