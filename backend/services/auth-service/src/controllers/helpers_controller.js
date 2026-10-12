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

module.exports = {
    lifeCheck
}