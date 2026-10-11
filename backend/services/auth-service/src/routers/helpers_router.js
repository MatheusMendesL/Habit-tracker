const express = require("express");
const functionsControllers = require("../controllers/helpers_controller");
const router = express.Router();

router.get("/lifeCheck", functionsControllers.lifeCheck)

module.exports = router 