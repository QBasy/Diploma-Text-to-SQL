# Diploma Work

## Text-to-SQL

### How to run?
#### You can run the .sh file by using CMD
```shell
"C:\Program Files\Git\bin\bash.exe" $PATHTOPROJECT"/Diploma-text-to-SQL/start.sh
```
#### or You can start the ```start.sh``` file by using your interpreter


### Structure

```assembly
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
        /models
            user.go
            item.go
        /controllers
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
