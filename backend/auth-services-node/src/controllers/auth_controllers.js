const jwt = require("jsonwebtoken");
const redis = require("../config/redis"); 

async function refresh(req, res) {
  const { refreshToken } = req.body;

  if (!refreshToken) {
    return res.status(401).json({ message: "Refresh token ausente" });
  }

  try {
    const payload = jwt.verify(
      refreshToken,
      process.env.JWT_REFRESH_SECRET
    );

    const userId = payload.id;

    const savedRefresh = await redis.get(`refresh:user:${userId}`);

    if (!savedRefresh || savedRefresh !== refreshToken) {
      return res.status(403).json({ message: "Refresh inválido" });
    }

    const newAccessToken = jwt.sign(
      { id: userId },
      process.env.JWT_SECRET,
      { expiresIn: "30m" }
    );

    return res.json({ accessToken: newAccessToken });

  } catch (err) {
    return res.status(403).json({ message: "Refresh inválido" });
  }
}

module.exports = { refresh };
