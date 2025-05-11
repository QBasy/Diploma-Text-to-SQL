# Diploma Work

## Text-to-SQL

### How to run?
#### You can run the .sh file by using CMD
```shell

```
#### or You can start the ```start.sh``` file by using your interpreter


### Structure

![image](https://github.com/user-attachments/assets/4bff83a0-f527-4107-9e81-3eac24398d97)


```go
Diploma-text-to-SQL/
	// A service that accepts requests from frontend 
	// and links Other services with one port
	// Go & Gin
    API/

	// Service working with user tokens
	// Login, Registration, Resetting password, SMTP
	// User information
    auth-service/
	
	// A service for direct work with the database
	// Go & GORM
    database-service/

	// Service working with AI, 
	// it translates natural language into a SQL queries
	// Go, Llama 4, Groqcloud
    text-to-SQL-service/

	// Visualising Data into Charts,
	// Pier, Scatter, Bar, Line
	// Go & SVG
    visualisation-service/
	
	// A frontend of this project
	// Vite, TypeScript, SvelteKit & TailwindCSS
    diploma-frontend/
```

### Client-Side Structure

```
src/
├── routes/
│   ├── +page.svelte
│   ├── +layout.svelte
│   ├── auth/
│   │   ├── +page.svelte
│   │   └── +layout.svelte
│   ├── documentation/
│   │   └── +page.svelte
│   ├── generate/
│   │   ├── complex/
│   │   │   └── +page.svelte
│   │   └── simple/
│   │       └── +page.svelte
│   └── profile/
│       ├── +page.svelte
│   	├── database/
│   	│   └── +page.svelte
│   	├── history/
│   	│   └── +page.svelte
│       └── settings/
│	    └── +page.svelte
├── lib/
│   ├── components/
│   │   ├── Navbar.svelte
│   │   ├── Footer.svelte
│   │   ├── Notification.svelte
│   │   └── LoadingSpinner.svelte
│   ├── stores/
│   │   ├── userStore.ts
│   │   ├── schemaStore.ts
│   │   ├── historyStore.ts
│   │   └── index.ts
│   ├── api/
│   │   ├── auth.ts
│   │   ├── database.ts
│   │   ├── history.ts
│   │   ├── text-to-sql.ts
│   │   ├── index.ts
│   └── types/
│       └── table.ts
└── app.html
```
