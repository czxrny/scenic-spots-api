# Scenic Spots API

The Scenic Spots API project provides an intuitive interface for client applications to access data about various scenic locations within a selected area. This includes viewpoints, lakes, rivers, trails, routes with beautiful scenery, and other hidden natural gems.

---

## About the Project

The API enables seamless retrieval of detailed information on natural attractions, making it easier to discover and explore beautiful spots.

## How it works

The API provides the following main endpoints:

- `/spot` – Search and create scenic spots.

- `/spot/{id}` – Search, update, and delete scenic spot specified by the ID.

- `/spot/{id}/review` – Submit, list, and delete reviews.

- `/spot/{id}/review/{rId}` – Search, update, and delete reviews specified by of a scenic spot specified the ID.

- `/user` – User registration, login, profile management, and account updates.

For detailed information on request methods, parameters, request/response formats, and authorization, please refer to the [swagger.yaml](docs/swagger.yaml) file.


## Project Structure

- `cmd/` – Contains the main application entry point (`main.go`).

- `internal/` – Includes initialization logic, HTTP handlers, database connection and functions, authorization, and logging components.

- `docs/` – Holds the API endpoint and database documentation.

- `utils/` – Contains reusable utility tools.

- `assets/` – Database seed files.

## Technologies Used

- Firebase Firestore (Cloud)
- Firebase Emulator

---

## Getting started

### 1. Configuartion

Before running the backend server, make sure to configure the necessary environment variables inside the `.env.example` file, then rename it to `.env`

### 2. Firebase Emulator (Optional)

For local development, you can use the Firebase emulator to simulate Firestore and Storage services:

- Ensure you have the Firebase CLI installed on your system.
- Initialize the emulator in the project root directory:

```bash
firebase init
```

During setup, select Firestore and Storage to configure the emulator.

Then, start the emulator with:
```bash
firebase emulators:start
```

> Running the emulator allows the backend to connect to local instances of Firestore and Storage instead of the live Firebase services.

### 3. Starting the API

To run the backend server, execute the following commands:

```bash
cd scenic-spots-api
go run cmd/main.go
```

This will start the API server on the configured port, making the endpoints available for client use.

---

## Postman tests

To run the Postman test:

1. Run the `tests/postman/setup.py` python script to initalize the enviroment variables needed for the tests.
2. Import the newly created `tests.postman_enviroment.json` and the [tests.postman_collection.json](tests/postman/postman-files/tests.postman_collection.json) into your Postman workspace.
3. Right click on the `tests` collection and select `Run` to run all of the tests.

> Remember to launch the backend before running the tests!