const AuthService = require("../services/auth_service");

async function lifeCheck(req, res) {
    res.json("Api is ok")
}

async function redisDebug(req, res) {
    try {
        const { data, keys } = await AuthService.get_keys()

        res.json({
            status: "success",
            keys,
            data
        });
    } catch (error) {
        res.status(500).json({
            status: "error",
            message: error?.message || "Erro no Redis"
        });
    }   
}

module.exports = { 
    redisDebug, 
    lifeCheck 
}