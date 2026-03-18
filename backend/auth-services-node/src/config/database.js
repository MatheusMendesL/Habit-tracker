// config mysql

const mysql = require('mysql2')
const db_data = require('./config_sql')

const conn = mysql.createConnection(db_data)

module.exports = conn