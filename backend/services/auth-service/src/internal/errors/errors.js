const createError = (message, code, status = 500) => {
  const error = new Error(message);
  error.code = code;
  error.status = status;
  return error;
};

const ERR_INVALID_CREDENTIALS = createError("Invalid credentials", "INVALID_CREDENTIALS", 401);
const ERR_EMAIL_ALREADY_EXISTS = createError("This email is already in use", "EMAIL_ALREADY_EXISTS", 409);
const ERR_USER_NOT_FOUND = createError("User not found", "USER_NOT_FOUND", 404);
const ERR_INVALID_TOKEN = createError("Invalid token", "INVALID_TOKEN", 401);
const ERR_INVALID_ARGUMENT = createError("Invalid arguments", "INVALID_ARGUMENT", 400);

module.exports = {
  createError,
  ERR_INVALID_CREDENTIALS,
  ERR_EMAIL_ALREADY_EXISTS,
  ERR_USER_NOT_FOUND,
  ERR_INVALID_TOKEN,
  ERR_INVALID_ARGUMENT,
};
