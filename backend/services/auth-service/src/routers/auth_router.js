const express = require("express");
const authController = require("../controllers/auth_controllers");
const jwt = require("../middlewares/jwt");
const router = express.Router();

router.post("/signup", authController.signup);
router.post("/login", authController.login);
router.post("/logout", jwt, authController.logout);
router.post("/refresh", authController.refresh);
router.get("/me", jwt, authController.getUserData);

module.exports = router;