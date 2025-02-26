# 🖼️ Image CDN Server

## 📖 About the Project

The **Image CDN Server** is a high-performance and secure content delivery network (CDN) designed to optimize **image storage and delivery**. Built with **Golang and the Gin framework**, it ensures **efficient image processing, fast access times, and robust security** using **JWT authentication**. Each client's images are stored in **separate directories**, providing a structured and isolated storage solution.

---

## ✨ Key Features

- 🔒 **JWT Authentication** – Secure access to upload and retrieve images.
- 📂 **Client-Specific Storage** – Images are stored in isolated directories per client.
- 🚀 **Optimized Image Delivery** – Designed for fast, scalable content serving.
- 📏 **Automatic Image Resizing** – Serve different image sizes dynamically.
- 🌍 **CDN Integration** – Can be used with global CDNs for enhanced performance.
- 🏗 **Scalable Architecture** – Supports large-scale deployments.
- 📑 **Detailed Logging** – Keep track of all image operations securely.

---

## 🛠 Technology Stack

- **Backend Framework**: [Gin](https://gin-gonic.com/) – High-performance HTTP web framework.
- **Storage**: Local file system or **cloud storage (S3, GCS)**.
- **Authentication**: [JWT](https://jwt.io/) for secure authentication.
- **Docker**: [Docker](https://www.docker.com/) – Containerized deployment.
- **Database**: PostgreSQL or SQLite (optional, for metadata tracking).

---

## 📦 Installation and Setup

### Prerequisites

- Install **[Go](https://golang.org/doc/install)**.
- Install **[Docker](https://www.docker.com/)** (optional, for containerized deployment).
- Set up **PostgreSQL** (if using metadata storage).

### Steps to Run

1. Clone the repository:
   ```bash
   git clone https://github.com/HieronimusMorgan/Image-CDN.git
   cd Image-CDN
   ```

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Configure environment variables:
   - Create a `.env` file in the root directory:
   ```env
   JWT_SECRET=your_secret_key
   STORAGE_PATH=/path/to/storage
   DB_HOST=localhost
   DB_PORT=5432
   DB_USER=your_db_user
   DB_PASSWORD=your_db_password
   DB_NAME=image_cdn
   ```

4. Start the application:
   ```bash
   go run main.go
   ```

---

## 🔗 API Endpoints

### 🔓 Public Routes
- `GET /health` → **Service health check**.

### 🔒 Protected Routes (Require Authentication)
- `POST /upload` → **Upload an image**.
- `GET /image/{client_id}/{filename}` → **Retrieve an image**.
- `DELETE /image/{client_id}/{filename}` → **Delete an image**.

---

## 📂 Project Structure

```
Image-CDN/
├── cmd/
│   ├── main.go          # Entry point of the application
├── config/              # Configuration files
├── internal/
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Authentication and security middleware
│   ├── services/        # Business logic
├── storage/             # Local storage for uploaded images
├── pkg/response/        # API response structures
└── README.md
```

---

## 🤝 Contributing

Contributions are **welcome**! Follow these steps:

1. Fork the repository.
2. Create a new branch (`git checkout -b feature/YourFeature`).
3. Commit your changes (`git commit -m 'Add YourFeature'`).
4. Push to the branch (`git push origin feature/YourFeature`).
5. Open a Pull Request.

For major updates, **open an issue** first to discuss your proposal.

---

## 📜 License

This project is licensed under the **MIT License**. See the `LICENSE` file for more details.

---

## 📧 Contact

- **Hieronimus Morgan** – morganhero35@gmail.com
