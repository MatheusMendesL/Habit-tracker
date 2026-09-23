const AuthService = require("../services/auth_service");
const { response } = require("../utils/functions");

function buildHttpInfo(req) {
    return {
        method: req.method,
        url: req.originalUrl || req.url,
    };
}

async function lifeCheck(req, res) {
    const startedAt = Date.now();
    return res.json(response("success", "Api is ok", { status: "ok" }, buildHttpInfo(req), null, startedAt));
}

async function redisDebug(req, res) {
    const startedAt = Date.now();

    try {
        const { data, keys } = await AuthService.get_keys();

        return res.json(response("success", "Redis ok", { keys, data }, buildHttpInfo(req), null, startedAt));
    } catch (error) {
        return res.status(500).json(
            response("error", error?.message || "Erro no Redis", null, buildHttpInfo(req), error?.message || "Erro no Redis", startedAt)
        );
    }
}

module.exports = {
    redisDebug,
    lifeCheck
}