# Diploma Work

## Text-to-SQL

### How to run?
#### You can run the .sh file by using CMD
```shell
"C:\Program Files\Git\bin\bash.exe" $PATHTOPROJECT"/Diploma-text-to-SQL/start.sh
```
#### or You can start the ```start.sh``` file by using your interpreter


### Structure

![image](https://github.com/user-attachments/assets/76d19dfa-8c7d-4f73-9271-80a9c865e593)


```go
Diploma-text-to-SQL/
	// A service that accepts requests from frontend 
	// and links database-service and text-to-sql
	// Go & Gin
    API/
        test/
            handlers_test.go
        main.go
        routes.go
        handlers.go
        go.mod
        .env
        Dockerfile
	
	// A service for direct work with the database
	// Go & GORM
    database-service/
        models/
            user.go
            item.go
        controllers/
            itemController.go
            userController.go
            customController.go
        main.go
        routes.go
        go.mod
        .env
        Dockerfile

	// A main service for this diploma work, 
	// it translates natural language into a SQL queries
	// Python, FastAPI & Transformers
    text-to-SQL-service/
        .venv/
        main.py
        requirements.txt
        text_to_sql.log
        Dockerfile
	
	// A frontend of this project
	// Vite, SvelteKit & TailwindCSS
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
│   │   ├── +layout.svelte
│   │   └── +page.server.js
│   ├── documentation/
│   │   └── +page.svelte
│   ├── generate/
│   │   ├── complex/
│   │   │   ├── +page.svelte
│   │   │   └── +page.server.js
│   │   └── simple/
│   │       ├── +page.svelte
│   │       └── +page.server.js
│   └── profile/
│       ├── +page.svelte
│       ├── +page.server.js
│   	├── database/
│   	│   ├── +page.svelte
│   	│   └── +page.server.js
│   	├── history/
│   	│   ├── +page.svelte
│   	│   └── +page.server.js
│       └── settings/
│	    ├── +page.svelte
│           └── +page.server.js
├── lib/
│   ├── components/
│   │   ├── Navbar.svelte
│   │   ├── Footer.svelte
│   │   ├── Notification.svelte
│   │   └── LoadingSpinner.svelte
│   └── stores/
│       ├── userStore.js
│       ├── schemaStore.js
│       └── historyStore.js
└── app.html
```
