const express = require("express");
const functionsControllers = require("../controllers/user_controller");
const jwt = require("../middlewares/jwt")
const router = express.Router();


router.post("/login", functionsControllers.login)
router.post("/signup", functionsControllers.signup)
router.get("/get_user_data/", jwt, functionsControllers.get_user_data)
router.post("/logout", jwt, functionsControllers.logout)

module.exports = router;