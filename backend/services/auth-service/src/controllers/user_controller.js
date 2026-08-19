const userRepository = require("../internal/repository/user_repository");
const { response } = require("../utils/functions");

async function getUserData(req, res) {
  try {
    const userId = req.userId;
    const result = await userRepository.getUserById(userId);

    return res.json(
      response("success", "Usuario encontrado", null, result.affectedRows, result.data)
    );
  } catch (error) {
    return res.status(error.status || 500).json(
      response("error", error.message, null, 0, null)
    );
  }
}

async function getUserDataById(req, res) {
  try {
    const userId = req.params.id;
    const result = await userRepository.getUserById(userId);

    return res.json(
      response("success", "Usuario encontrado", null, result.affectedRows, result.data)
    );
  } catch (error) {
    return res.status(error.status || 500).json(
      response("error", error.message, null, 0, null)
    );
  }
}

module.exports = {
  getUserData,
  getUserDataById,
};