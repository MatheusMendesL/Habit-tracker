const authService = require("../internal/service/auth_service");
const { response } = require("../utils/functions");

async function signup(req, res) {
  try {
    const result = await authService.signup(req.body);
    return res.status(201).json(
      response("success", "Usuario criado com sucesso", null, 1, {
        user: result.user,
        accessToken: result.accessToken,
        refreshToken: result.refreshToken,
      })
    );
  } catch (error) {
    return res.status(error.status || 500).json(
      response("error", error.message, null, 0, null)
    );
  }
}

async function login(req, res) {
  try {
    const result = await authService.login(req.body);
    return res.json(
      response("success", "Login realizado", null, 1, {
        user: result.user,
        accessToken: result.accessToken,
        refreshToken: result.refreshToken,
      })
    );
  } catch (error) {
    return res.status(error.status || 500).json(
      response("error", error.message, null, 0, null)
    );
  }
}

async function logout(req, res) {
  try {
    const userId = req.userId;
    await authService.revokeRefreshToken(userId);
    return res.json(response("success", "Logout realizado", null, 1, null));
  } catch (error) {
    return res.status(error.status || 500).json(
      response("error", error.message, null, 0, null)
    );
  }
}

async function refresh(req, res) {
  try {
    const { refreshToken } = req.body;
    const result = await authService.refreshAccessToken(refreshToken);
    return res.json(response("success", "Token renovado", null, 1, result));
  } catch (error) {
    return res.status(error.status || 500).json(
      response("error", error.message, null, 0, null)
    );
  }
}

async function getUserData(req, res) {
  return res.json(response("success", "Usuario autenticado", null, 1, { id: req.userId }));
}

module.exports = {
  signup,
  login,
  logout,
  refresh,
  getUserData,
};
