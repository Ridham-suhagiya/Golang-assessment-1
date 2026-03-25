# Country Search API - Code Explanation

This document explains how the code for the Country Search API works, from start to finish. We've built an HTTP server in Go that searches for countries using the REST Countries API and caches the results so that repeated searches are lightning-fast.

---

## 🏗️ Architecture Layers
The application is built using a clean, layered architecture. This means each part has a specific job and doesn't do anything else.

1. **`main.go` (The Wire-Up):** The starting point. It initializes all components and starts the server.
2. **`handler` (The Traffic Cop):** Takes incoming web requests (HTTP), reads the variables, and returns JSON. They don't do any complex logic; they just deal with web traffic.
3. **`service` (The Brain):** The core business logic. It coordinates the cache and the external API.
4. **`client` (The Messenger):** Knows how to talk to external websites over the internet (REST Countries API).
5. **`cache` (The Memory):** A fast, temporary storage area to remember previous search results.

---

## 📖 Let's walk through an example request!

Imagine a user makes the following web request to our server:
`GET http://localhost:8000/api/countries/search?name=canada`

Here's exactly what happens, step-by-step:

### Step 1: The Handler (`internal/handler/country.go`)
- The Go HTTP router (`http.ServeMux` in `main.go`) sees the `/api/countries/search` path and routes the traffic to our `CountryHandler.Search` function.
- The handler looks at the URL and extracts `name=canada`.
- The handler then asks the service: *"Hey Service, can you find me a country called 'canada'?"*

### Step 2: The Service (`internal/service/country.go`)
- The `CountryService` receives the request for `"canada"`.
- It cleans up the word (e.g., `" CAnada "` becomes `"canada"`) so the cache key is consistent.
- **Cache Check:** The service checks the `cache` first. *"Do we already know the details for 'canada'?"*
    - **It's a Cache Miss:** If this is the first time anyone searched for Canada, the cache replies: *"No, I don't have it."*
- **External Request:** Because we don't have it, the service asks the client: *"Hey Client, please go to the internet and fetch me 'canada'."*

### Step 3: The Client (`internal/client/restcountries.go`)
- The `RESTCountriesClient` builds a URL for the external API: `https://restcountries.com/v3.1/name/canada`.
- It shoots an HTTP request over the internet to that URL.
- When the external API responds with a giant block of JSON data, the client receives it.
- **Data Cleanup:** The client reads the JSON and picks out only what we need (Name, Capital, Currency, Population). It packages this into our neat `models.Country` struct and gives it back to the Service.

### Step 4: Back to the Service (`internal/service/country.go`)
- The service receives the packaged `models.Country` from the client.
- **Save for later:** The service pushes this data into the memory `cache`. Now, if someone searches for "canada" again in 5 seconds, the service will stop at the cache and skip Step 3 entirely! (Cache HIT).
- The service returns the data back to the Handler.

### Step 5: Back to the Handler (`internal/handler/country.go`)
- The handler receives the `models.Country` object.
- It prepares a good web response by setting the status code to `200 OK` and Content-Type to `application/json`.
- Finally, it uses `json.NewEncoder` to convert our Go struct back into text-based JSON and sends it over the network to the user's browser.

**Final JSON Output:**
```json
{
  "name": "Canada",
  "capital": "Ottawa",
  "currency": "$",
  "population": 41651653
}
```

---

## ✨ Code Simplification Efforts
To make things "more human", you'll notice:
1. **Meaningful variable names** (`request` instead of `r`, `service` instead of `svc`).
2. **Simple printing** using standard `fmt.Printf` instead of `log` packages.
3. **Easy-to-read server bootups** using a plain, blocking `http.ListenAndServe` instead of complex context channels.
