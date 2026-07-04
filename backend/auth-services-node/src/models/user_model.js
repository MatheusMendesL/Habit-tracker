const pool = require("../config/database")
const { hashPass } = require("../utils/functions")

async function get_user_data(userId) {
    if (!userId) throw new Error("You need an id");

    console.log(userId);

    const query_sql = "SELECT * FROM users WHERE id = $1";
    const results = await pool.query(query_sql, [userId]);

    return {
        query_sql,
        affectedRows: results.rowCount,
        data: results.rows,
    };
}

async function signup(data) {
    if (!data) throw new Error("You need to put data");

    const password = await hashPass(data.password);

    const query_sql =
        "INSERT INTO users(name, email, tel, password, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id";
    const results = await pool.query(query_sql, [data.name, data.email, data.tel, password]);

    const insertId = results.rows[0].id;
    const data_id = { id: insertId };
    const data_user = await get_user_data(data_id);

    return {
        query_sql,
        affectedRows: results.rowCount,
        data: data_user.data[0],
        insertId: insertId,
    };
}

async function findByEmail(email) {
    if (!email) throw new Error("You need to put an email");

    const query_sql = "SELECT * FROM users WHERE email = $1";
    const results = await pool.query(query_sql, [email]);

    return {
        query_sql,
        affectedRows: results.rowCount,
        data: results.rows,
    };
}

module.exports = {
    get_user_data,
    signup,
    findByEmail
}