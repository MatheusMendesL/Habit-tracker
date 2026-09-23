const pool = require("../../config/database");

async function getUserById(userId) {
  if (!userId) {
    throw new Error("You need an id");
  }

  const querySql = "SELECT * FROM users WHERE id = $1";
  const results = await pool.query(querySql, [userId]);

  return {
    querySql,
    affectedRows: results.rowCount,
    data: results.rows,
  };
}

async function findByEmail(email) {
  if (!email) {
    throw new Error("You need to put an email");
  }

  const querySql = "SELECT * FROM users WHERE email = $1";
  const results = await pool.query(querySql, [email]);

  return {
    querySql,
    affectedRows: results.rowCount,
    data: results.rows,
  };
}

async function createUser(data) {
  if (!data) {
    throw new Error("You need to put data");
  }

  const querySql =
    "INSERT INTO users(name, email, tel, password, created_at) VALUES ($1, $2, $3, $4, NOW()) RETURNING id";

  const results = await pool.query(querySql, [
    data.name,
    data.email,
    data.tel,
    data.password,
  ]);

  const insertId = results.rows[0]?.id;
  const user = await getUserById(insertId);

  return {
    querySql,
    affectedRows: results.rowCount,
    data: user.data[0],
    insertId,
  };
}

module.exports = {
  getUserById,
  findByEmail,
  createUser,
};
