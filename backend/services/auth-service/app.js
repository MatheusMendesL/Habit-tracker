require("dotenv").config();
const express = require("express");
const cors = require("cors");

const app = express();

const userRouter = require("./src/routers/user_router");
const helperRouter = require("./src/routers/helpers_router");
const authRouter = require("./src/routers/auth_router");

app.use(cors());
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

app.use("/users", userRouter);
app.use("/helpers", helperRouter);
app.use("/auth", authRouter);

module.exports = app;
