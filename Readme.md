# 📊 Automatic Data Visualization from Natural Language Queries

This project enables users to type natural language questions like  
**"Show total revenue by month for 2023"**  
and instantly receive:
- An **auto-generated SQL query**
- An **interactive visualization** (SVG chart)

It is designed to support **non-technical users** with no SQL knowledge and is powered by a scalable **microservices architecture** and an **LLM-based query engine**.

---

### 📝 Natural Language Query Input

![Снимок экрана 2025-05-16 202918](https://github.com/user-attachments/assets/6e58a0de-87d3-4d34-b96f-ae398b269502)

---

### 🧠 Auto-generated SQL Query

![Снимок экрана 2025-05-16 202932](https://github.com/user-attachments/assets/bd28253b-996a-4198-bc72-0acb6af8d36f)

---

### 📈 Real-time Visualization

![Снимок экрана 2025-05-16 141739](https://github.com/user-attachments/assets/39f496e2-a965-4db2-8452-43941689ff38)
![Снимок экрана 2025-05-16 202954](https://github.com/user-attachments/assets/d8cd9c44-564e-43ae-9c65-ba45e4b43c1f)

---

## 🧠 How It Works

- Text queries are processed using an **LLM (meta-llama/llama-4-scout-17b-16e-instruct)** hosted via **GroqCloud**
- The system selects **simple or complex prompts** dynamically depending on query complexity
- SQL is executed against a PostgreSQL or SQLite database
- If results are visualizable, the system returns a **scalable SVG chart** for the frontend

---

## 🖥 Frontend

### ➤ `diploma-frontend`
- **Framework**: [SvelteKit](https://kit.svelte.dev/)
- **Styling**: [TailwindCSS](https://tailwindcss.com/)
- **UI/Icons**: [Lucide](https://lucide.dev/)
- **Features**:
  - Natural language input
  - Schema uploading
  - Visual feedback
  - Authentication and history tracking

---

## ⚙️ Backend (Microservices)

| Service                | Tech Stack                               | Description                                                                 |
|------------------------|-------------------------------------------|-----------------------------------------------------------------------------|
| `API Gateway`          | Go + Gin                                  | Routes requests, handles rate limiting and CORS                            |
| `text-to-sql-service`  | Go + Groqcloud + Llama                    | Converts natural language to SQL using LLaMA 4 Scout                       |
| `database-service`     | Go + Gin + PostgreSQL/SQLite + GORM       | Executes SQL, manages user DBs, exposes schema                             |
| `visualisation-service`| Go + gRPC + SVG                           | Converts SQL result into vector visualizations                             |
| `auth-service`         | Go + Gin + JWT + PostgreSQL               | Handles registration, login, token auth, Google OAuth                      |
| `history-service`      | Go + Gin + PostgreSQL                     | Saves and returns past user queries                                        |

  
---

## 📁 Project Structure

Diploma-text-to-SQL/

├── API/

├── auth-service/

├── database-service/

├── text-to-sql-service/

├── visualisation-service/

└── diploma-frontend/

---

### 🎨 Client Side Structure

src/

├── routes/

│ ├── auth/

│ ├── documentation/

│ ├── generate/simple/

│ ├── generate/complex/

│ └── profile/ (database, history, settings)

├── lib/

│ ├── components/

│ ├── stores/

│ ├── api/

│ └── types/

└── app.html

---

## 🛡️ Security

- **JWT-based authentication**
- **CORS** for safe frontend/backend interaction
- **Rate limiting** (100 requests/minute per IP)

---

## 🔧 Tools and Technologies Used

- **LLM**: meta-llama/llama-4-scout-17b-16e-instruct
- **Frontend**: SvelteKit + TailwindCSS + Lucide
- **Backend**: Go (Gin, GORM), FastAPI, gRPC
- **Database**: PostgreSQL, SQLite
- **Visualization**: SVG chart service (custom, gRPC-based)
- **Security**: JWT, HTTPS, CORS, Rate Limiting

---

## 🚀 Run the Project (Coming Soon)

> _Setup instructions or Dockerfile links can go here_





