const functionsModel = require("../models/user_model")
const { response, comparePass } = require("../utils/functions");
const AuthService = require("../services/auth_service");
const jwt = require("jsonwebtoken");

async function signup(req, res) {

    const { name, email, tel, password } = req.body

    const data_signup = {
        name,
        email,
        tel,
        password,
    };

    try {

        const searchEmail = await functionsModel.findByEmail(email);

        if (searchEmail.data.length > 0) {
            return res.status(401).json(
                response("error", "This email already exists", null, 0, null)
            );
        }

        const { query_sql, affectedRows, data, insertId } =
            await functionsModel.signup(data_signup);

        if (!insertId)
            return res
                .status(500)
                .json({ status: "error", message: "Erro ao criar usuário" });


        const accessToken = jwt.sign(
            { id: insertId },
            process.env.JWT_SECRET,
            { expiresIn: "30m" }
        );

        const refreshToken = jwt.sign(
            { id: insertId },
            process.env.JWT_REFRESH_SECRET,
            { expiresIn: "7d" }
        );

        await AuthService.signup({
            id: insertId,
            refreshToken
        });
        
        res.json(
            response("success", "User added successfully", query_sql, affectedRows, {
                data,
                insertId,
                accessToken,
                refreshToken
            })
        )


    } catch (error) {
        res.status(500).json(
            response("error", error.message, null, 0, null)
        );
    }
}

async function login(req, res) {
    try {
        const { email, password } = req.body;

        const user = await functionsModel.findByEmail(email);
        if (!user) {
            return res.status(401).json(
                response("error", "Credenciais inválidas", user, 0, null)
            );
        }

        const valid = await comparePass(password, user.data[0].password);
        if (!valid) {
            return res.status(401).json(
                response("error", "Credenciais inválidas", null, 0, null)
            );
        }

        const id = user.data[0].id
        const accessToken = jwt.sign(
            { id: id },
            process.env.JWT_SECRET,
            { expiresIn: "30m" }
        );

        const refreshToken = jwt.sign(
            { id: id },
            process.env.JWT_REFRESH_SECRET,
            { expiresIn: "7d" }
        );

        await AuthService.login({
            id: id,
            refreshToken
        });

        const user_data = user.data[0]
        res.json(
            response("success", "Login realizado", null, 1, {
                user_data,
                accessToken,
                refreshToken
            })
        );

    } catch (error) {
        res.status(500).json(
            response("error", error.message, null, 0, null)
        );
    }
}

async function get_user_data(req, res) {

    const userId = req.userId;

    try {
        const { query_sql, affectedRows, data } =
            await functionsModel.get_user_data(userId);
        res.json(
            response(
                "success",
                "Got user data successfully",
                query_sql,
                affectedRows,
                { "data": data, "token": 123 }
            )
        )
    } catch (error) {
        res.status(500).json(
            response("error", error.message, null, 0, null)
        )
    }
}

async function logout(req, res) {
    try {
        const userId = req.userId;
        await AuthService.logout(userId);

        res.json(
            response("success", "Logout realizado", userId, 1, null)
        );
    } catch (error) {
        res.status(500).json(
            response("error", error.message, null, 0, null)
        );
    }
}

module.exports = {
    login,
    signup,
    get_user_data,
    logout
}