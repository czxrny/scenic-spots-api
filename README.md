# SCENIC SPOTS API

The “Scenic Spots API” project aims to enable intuitive communication with a client application by providing data about picturesque locations in a selected area — such as viewpoints, lakes, rivers, trails, and routes with beautiful scenery, as well as other types of “hidden natural gems“.

### ❗Please check /doc/swagger.yaml for endpoint documentation❗

# ✅ TODO – Scenic Spots API (Project Roadmap)

> A structured development plan for the REST API project written in Go, using Swagger for documentation and Firebase (Firestore + Storage) for data and media management.

---

## 🔧 Project Setup & Structure

- [x] Initialize Git repository and Go module (`go mod init`)
- [x] Define clean folder structure (`cmd/`, `internal/`, `api/`, `docs/`, etc.)
- [x] Add `README.md` with project description
- [x] Create initial Swagger/OpenAPI file (`docs/swagger.yaml`)

---

## 📄 API Documentation (Swagger / OpenAPI)

- [x] Define main schemas:
  - [x] `Spot` – full spot data
  - [x] `NewSpot` – required fields for creation
  - [x] `Review` – reviews per spot
- [x] Spot endpoints:
  - [x] `GET /spot` – search or list all scenic spots with filters
  - [x] `GET /spot/{id}` – get a specific spot by ID
  - [x] `POST /spot` – create a new spot
  - [x] `PUT /spot/{id}` – update spot data
  - [x] `DELETE /spot/{id}` – delete a spot
- [x] Spot photos:
  - [x] `PUT /spot/{id}/addPhoto` – upload up to 3 images (multipart/form-data)
  - [x] Return uploaded image URLs (hosted on Firebase/GCP)
- [x] Review endpoints:
  - [x] `POST /spot/{id}/review` – submit a review
  - [x] `GET /spot/{id}/review` – list reviews for a spot
  - [x] `DELETE /spot/{id}/review/{reviewId}` – delete a specific review

---

## 🧠 Backend Implementation (Go)

- [ ] App entrypoint (`main.go`), router and API handler setup
- [ ] Firestore integration for data storage:
  - [ ] Collections: `spots`, `reviews`, `users`
- [ ] Firebase Storage integration for image hosting
- [ ] Data validation and error handling (400, 404, 500, etc.)

---

## 🔐 Authentication & User Management

- [ ] Plan and define `User` model
- [ ] Choose authentication method (e.g. Firebase Auth or JWT)
- [ ] Add security schemes to Swagger (`securitySchemes`)
- [ ] Design user-related endpoints (login, register, user info)

---

## 🌐 Deployment & Final Touches

- [ ] Deploy backend (e.g. Google Cloud Run or App Engine)
- [ ] Add integration tests (Postman / Go test)
- [ ] Code cleanup, documentation, refactor

---

**🔄 Status**: API structure and documentation in progress. Core endpoints drafted. Backend development ongoing.