const express = require("express");
const functionsControllers = require("../controllers/auth_controllers");
const router = express.Router();

router.get("/refresh", functionsControllers.refresh)

module.exports = router 