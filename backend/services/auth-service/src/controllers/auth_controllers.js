const authService = require("../internal/service/auth_service");
const { response } = require("../utils/functions");

function buildHttpInfo(req) {
  return {
    method: req.method,
    url: req.originalUrl || req.url,
  };
}

async function signup(req, res) {
  const startedAt = Date.now();

  try {
    const result = await authService.signup(req.body);
    return res.status(201).json(
      response("success", "Usuario criado com sucesso", {
        user: result.user,
        accessToken: result.accessToken,
        refreshToken: result.refreshToken,
      }, buildHttpInfo(req), null, startedAt)
    );
  } catch (error) {
    return res.status(error.status || 500).json(
      response("error", error.message || "Erro ao criar usuario", null, buildHttpInfo(req), error.message || "Erro ao criar usuario", startedAt)
    );
  }
}

async function login(req, res) {
  const startedAt = Date.now();

  try {
    const result = await authService.login(req.body);
    return res.json(
      response("success", "Login realizado", {
        user: result.user,
        accessToken: result.accessToken,
        refreshToken: result.refreshToken,
      }, buildHttpInfo(req), null, startedAt)
    );
  } catch (error) {
    return res.status(error.status || 500).json(
      response("error", error.message || "Erro ao realizar login", null, buildHttpInfo(req), error.message || "Erro ao realizar login", startedAt)
    );
  }
}

async function logout(req, res) {
  const startedAt = Date.now();

  try {
    const userId = req.userId;
    await authService.revokeRefreshToken(userId);
    return res.json(response("success", "Logout realizado", null, buildHttpInfo(req), null, startedAt));
  } catch (error) {
    return res.status(error.status || 500).json(
      response("error", error.message || "Erro ao fazer logout", null, buildHttpInfo(req), error.message || "Erro ao fazer logout", startedAt)
    );
  }
}

async function refresh(req, res) {
  const startedAt = Date.now();

  try {
    const { refreshToken } = req.body;
    const result = await authService.refreshAccessToken(refreshToken);
    return res.json(response("success", "Token renovado", result, buildHttpInfo(req), null, startedAt));
  } catch (error) {
    return res.status(error.status || 500).json(
      response("error", error.message || "Erro ao renovar token", null, buildHttpInfo(req), error.message || "Erro ao renovar token", startedAt)
    );
  }
}

async function getUserData(req, res) {
  const startedAt = Date.now();
  return res.json(response("success", "Usuario autenticado", { id: req.userId }, buildHttpInfo(req), null, startedAt));
}

module.exports = {
  signup,
  login,
  logout,
  refresh,
  getUserData,
};
