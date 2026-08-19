const internalAuthService = require("../internal/service/auth_service");

module.exports = {
  login: internalAuthService.login,
  signup: internalAuthService.signup,
  logout: internalAuthService.revokeRefreshToken,
  get_keys: internalAuthService.getRedisKeys,
  refresh: internalAuthService.refreshAccessToken,
};