# 📁 assets/seeds

This folder contains **test seed data** used for development and testing purposes.

## ✅ Contents

- **Spots** – sample tourist locations such as castles, mountains, and natural reserves.
- **Reviews** – example user reviews for each spot.
- **Users** – 3 predefined test accounts:

| Username | Password   | Role  |
|----------|------------|-------|
| `user1`  | `user123`  | user  |
| `user2`  | `user123`  | user  |
| `admin`  | `admin123` | admin |

> 🔐 Passwords are already hashed and ready for database seeding.

## 📌 Notes

- These records are meant **only for development or testing environments**.
- You can load this data during app via .env *DB_POPULATE* variable being set to *true*.
