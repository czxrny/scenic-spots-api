# 📚 Firestore Database Documentation for Scenic Spots Api

## 🔑 Overview
This document describes the structure of key collections in the Firestore database that *Scenic-Spots-Api* implements. The following provides the fields and expected data types for each collection.

---

## 📍 Collection: **Spots**
- **Description**: The **Spots** collection contains information about various places posted by users.
- **Documents**:
    - `id` (string): Unique identifier for the spot.
    - **Fields**:
        - `name` (string): Name of the spot. 
        - `description` (string): A short description of the spot.
        - `latitude` (float): Latitude of the spot’s location.
        - `longitude` (float): Longitude of the spot’s location.
        - `category` (string): Category of the spot (e.g., `restaurant`, `tourist attraction`).
        - `photos` (array of strings): A list of URLs to photos of the spot.
        - `addedBy` (string): User ID of the person who added the spot.
        - `createdAt` (timestamp): Timestamp indicating when the spot was added.

#### Example Document in JSON:
```json
{
  "id": "spot123",
  "name": "Central Park",
  "description": "A large park in New York City, perfect for outdoor activities.",
  "latitude": 40.785091,
  "longitude": -73.968285,
  "category": "park",
  "photos": [
    "https://example.com/images/central_park_1.jpg",
    "https://example.com/images/central_park_2.jpg"
  ],
  "addedBy": "user456",
  "createdAt": "2025-05-13T10:00:00Z"
}
```

## ⭐ Collection: **Reviews**
- **Description**: The **Reviews** collection stores reviews submitted by users for different spots. Each review includes a rating, content, and the user who submitted it.
- **Documents**:
    - `id` (string): Unique identifier for the review.
    - **Fields**:
        - `spotId` (string): ID of the spot being reviewed.
        - `rating` (float): Rating given to the spot (between 0 and 5).
        - `content` (string): The user comment on the rating.
        - `addedBy` (string): User ID of the person who added the review.
        - `createdAt` (timestamp): Timestamp indicating when the review was created.

#### Example Document in JSON:
```json
{
  "id": "review789",
  "spotId": "spot123",
  "rating": 4.5,
  "content": "Great place to relax and enjoy nature. Highly recommend!",
  "addedBy": "user123",
  "createdAt": "2025-05-13T11:00:00Z"
}
```

## 🧑‍💻 Collection: **User**
`TO BE IMPLEMENTED`