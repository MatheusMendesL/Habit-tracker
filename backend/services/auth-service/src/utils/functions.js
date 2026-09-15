const bcrypt = require("bcrypt")

function response(status, message, query, affected_rows, data) {
    return {
        status: status,
        message: message,
        query: query,
        affected_rows: affected_rows,
        timestamp: new Date().toISOString(),
        data: data
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