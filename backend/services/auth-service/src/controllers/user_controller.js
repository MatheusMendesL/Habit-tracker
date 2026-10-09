const userRepository = require("../internal/repository/user_repository");
const { response } = require("../utils/functions");

function buildHttpInfo(req) {
  return {
    method: req.method,
    url: req.originalUrl || req.url,
  };
}

async function getUserData(req, res) {
  const startedAt = Date.now();

  try {
    const userId = req.userId;
    const result = await userRepository.getUserById(userId);

    return res.json(
      response("success", "Usuario encontrado", result.data, buildHttpInfo(req), null, startedAt)
    );
  } catch (error) {
    return res.status(error.status || 500).json(
      response("error", error.message || "Erro ao buscar usuario", null, buildHttpInfo(req), error.message || "Erro ao buscar usuario", startedAt)
    );
  }
}

async function getUserDataById(req, res) {
  const startedAt = Date.now();

  try {
    const userId = req.params.id;
    const result = await userRepository.getUserById(userId);

    return res.json(
      response("success", "Usuario encontrado", result.data, buildHttpInfo(req), null, startedAt)
    );
  } catch (error) {
    return res.status(error.status || 500).json(
      response("error", error.message || "Erro ao buscar usuario", null, buildHttpInfo(req), error.message || "Erro ao buscar usuario", startedAt)
    );
  }
}

module.exports = {
  getUserData,
  getUserDataById,
};