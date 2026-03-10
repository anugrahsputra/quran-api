# Quran API - RapidAPI Documentation

This document provides technical details for integrating the **Quran API** via RapidAPI.

## 🌟 Overview

The Quran API is a high-performance, production-ready RESTful service providing access to:

- **Al-Quran Text**: All 114 Surahs with full verse details.
- **Advanced Search**: Full-text search across translations, Tafsir, and thematic topics.
- **Tafsir**: Detailed interpretations for every verse.
- **Prayer Times**: Accurate prayer times for any location globally.

## 🔗 RapidAPI Link

[Quran API on RapidAPI](https://rapidapi.com/downormal/api/quran-api-v1) _(Replace with actual link if available)_

## 🔐 Authentication

All requests through RapidAPI must include your RapidAPI Key. The RapidAPI gateway will handle the translation of this key to our backend.

**Required Headers:**

- `X-RapidAPI-Key`: Your unique RapidAPI Subscription Key.
- `X-RapidAPI-Host`: `quran-api.p.rapidapi.com` _(Example host)_

## 🚀 Endpoints

### 1. Health Check (Ping)

Verify the API is alive and responsive. Recommended for uptime monitoring.

- **Method**: `GET`
- **Path**: `/ping`
- **Response**: `{"message": "pong"}`

---

### 2. List All Surahs

Retrieve a complete list of all 114 chapters (Surahs).

- **Method**: `GET`
- **Path**: `/api/v1/surah/`
- **Response Structure**:
  ```json
  {
    "status": 200,
    "message": "success",
    "data": [
      {
        "id": 1,
        "name": "Al-Fatihah",
        "transliteration": "Al-Fatihah",
        "translation": "Pembukaan",
        "total_verses": 7,
        "type": "Makkiyah"
      },
      ...
    ]
  }
  ```

---

### 3. Get Surah Detail

Retrieve specific information and verses for a Surah with pagination support.

- **Method**: `GET`
- **Path**: `/api/v1/surah/:surah_id/`
- **Parameters**:
  - `surah_id` (Path, required): ID of the Surah (1-114).
  - `page` (Query, optional): Page number (default: 1).
  - `limit` (Query, optional): Verses per page (default: 10, max: 100).
- **Example**: `/api/v1/surah/1/?page=1&limit=5`

---

### 4. Get Ayah Detail

Retrieve detailed information for a specific absolute Ayah ID.

- **Method**: `GET`
- **Path**: `/api/v1/ayah/:ayah_id/`
- **Parameters**:
  - `ayah_id` (Path, required): Absolute ID of the verse (1-6236).
- **Example**: `/api/v1/ayah/1/`

---

### 5. Advanced Search

Search the Quran using full-text indexing across Arabic text, translations, Tafsir, and topics.

- **Method**: `GET`
- **Path**: `/api/v1/search`
- **Parameters**:
  - `q` (Query, required): Search query (e.g., "patience", "iman").
  - `page` (Query, optional): Page number.
  - `limit` (Query, optional): Results per page.
- **Example**: `/api/v1/search?q=faith&limit=10`

---

### 6. Prayer Times

Get prayer times for a specific city and country.

- **Method**: `GET`
- **Path**: `/api/v1/prayer-time/`
- **Parameters**:
  - `city` (Query, required): City name (e.g., "Jakarta").
  - `country` (Query, required): Country name (e.g., "Indonesia").
  - `date` (Query, optional): Date in `YYYY-MM-DD` format (default: today).
- **Example**: `/api/v1/prayer-time/?city=Jakarta&country=Indonesia`

## 📦 Response Formats

### Success Response

```json
{
  "status": 200,
  "message": "success",
  "data": { ... }
}
```

### Error Response

```json
{
  "status": 400,
  "message": "Error description here"
}
```

## ⏱️ Rate Limits

The API is subject to rate limiting through the RapidAPI gateway. Standard limits apply based on your subscription plan. If you receive a `429 Too Many Requests` response, please slow down your request rate.

## 🛠️ Best Practices

1. **Caching**: Cache Surah lists and Ayah details on your client-side as they are static data.
2. **Pagination**: Use `limit` and `page` parameters for Surah details to ensure fast response times and lower data consumption.
3. **Timeouts**: We recommend setting a client-side timeout of at least 10 seconds for search and prayer time requests.

---

_Generated for Quran API Integration._
