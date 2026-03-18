const express = require("express");
const functionsControllers = require("../controllers/helpers_controller");
const router = express.Router();

router.get("/lifeCheck", functionsControllers.lifeCheck)
router.get("/redisKeys", functionsControllers.redisDebug)

module.exports = router 