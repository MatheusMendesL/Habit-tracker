const conn = require("../config/database")
const { hashPass } = require("../utils/functions")

function get_user_data(data) {
    return new Promise((resolve, reject) => {
        if (!data) return reject(new Error("You need an id"));

        const query_sql = "SELECT * FROM users WHERE id = ?";

        conn.query(query_sql, [data.id], (error, results) => {
            if (error) return reject(error);
            resolve({
                query_sql,
                affectedRows: results.length,
                data: results,
            });
        });
    });
}

async function signup(data) {
    return new Promise(async (resolve, reject)  => {
        if (!data) return reject(new Error("You need to put data"));

        var password = await hashPass(data.password)

        const query_sql =
            "INSERT INTO users(name, email, tel, password, created_at ) VALUES (?, ?, ?, ?, NOW())";
        conn.query(
            query_sql,
            [data.name, data.email, data.tel, password],
            async (error, results) => {
                if (error) return reject(error);
                try {
                    const data_id = {
                        id: results.insertId,
                    };
                    const data_user = await get_user_data(data_id);
                    resolve({
                        query_sql,
                        affectedRows: results.affectedRows,
                        data: data_user.data[0],
                        insertId: results.insertId,
                    });
                } catch (error) {
                    reject(error);
                }
            }
        );
    });
}

function findByEmail(email) {
    return new Promise((resolve, reject) => {
        if (!email) return reject(new Error("You need to put an email"));

        const query_sql = "SELECT * FROM users WHERE email = ?"
        conn.query(query_sql, [email], (error, results) => {
            if (error) return reject(error)
            resolve({
                query_sql,
                affectedRows: results.length,
                data: results,
            })
        })
    })
}

module.exports = {
    get_user_data,
    signup,
    findByEmail
}