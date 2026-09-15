const express = require("express");
const userController = require("../controllers/user_controller");
const jwt = require("../middlewares/jwt");
const router = express.Router();

router.get("/me", jwt, userController.getUserData);
router.get("/data/:id", userController.getUserDataById);

module.exports = router;