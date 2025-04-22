<script lang="ts">
    import { PlayCircle, Globe,Settings,AlertTriangle, Terminal, GitBranch, Container, Database } from 'lucide-svelte';
</script>

<h1 class="text-3xl font-bold mb-4 text-gray-800">🚀 How to Start</h1>
<p class="text-gray-600 mb-8">
    This guide walks you through the process of setting up and running the automatic data visualization system locally on your machine.
</p>

<div class="space-y-12">
    <!-- Prerequisites Section -->
    <div class="bg-white rounded-xl shadow p-6">
        <div class="flex items-center gap-3 mb-4">
            <Terminal class="w-8 h-8 text-blue-600" />
            <h2 class="text-2xl font-semibold text-gray-800">Prerequisites</h2>
        </div>

        <p class="text-gray-600 mb-4">Before getting started, make sure you have the following tools installed on your system:</p>

        <ul class="list-disc list-inside space-y-1 text-gray-600 ml-4">
            <li>Docker (20.10.x or newer)</li>
            <li>Docker Compose (v2.x or newer)</li>
            <li>Git</li>
            <li>A GroqCloud API key (for text-to-SQL functionality)</li>
        </ul>

        <div class="mt-6 bg-yellow-50 p-4 rounded-lg">
            <h3 class="text-lg font-medium mb-2 text-yellow-800">Note on API Keys</h3>
            <p class="text-gray-700">
                To use the text-to-SQL translation feature powered by Llama models, you'll need to obtain an API key from GroqCloud. Visit <a href="https://console.groq.com/" class="text-blue-600 hover:underline" target="_blank" rel="noopener noreferrer">GroqCloud Console</a> to sign up and obtain your key.
            </p>
        </div>
    </div>

    <div class="bg-white rounded-xl shadow p-6">
        <div class="flex items-center gap-3 mb-4">
            <GitBranch class="w-8 h-8 text-green-600" />
            <h2 class="text-2xl font-semibold text-gray-800">Clone the Repository</h2>
        </div>

        <p class="text-gray-600 mb-4">Open your terminal and run the following commands to clone the repository and navigate to its directory:</p>

        <div class="bg-gray-800 text-green-400 p-4 rounded-md overflow-x-auto text-sm mb-4">
            <code>git clone https://github.com/QBasy/Diploma-Text-to-SQL.git</code><br>
            <code>cd Diploma-Text-to-SQL</code>
        </div>

        <p class="text-gray-600">The repository contains all necessary services, Dockerfiles, and a Docker Compose configuration for easy setup.</p>
    </div>

    <!-- Configuration Section -->
    <div class="bg-white rounded-xl shadow p-6">
        <div class="flex items-center gap-3 mb-4">
            <Settings class="w-8 h-8 text-purple-600" />
            <h2 class="text-2xl font-semibold text-gray-800">Environment Configuration</h2>
        </div>

        <p class="text-gray-600 mb-4">Create a <code class="bg-gray-100 px-1 rounded text-sm">.env</code> file in the root directory with the following variables:</p>

        <div class="bg-gray-800 text-green-400 p-4 rounded-md overflow-x-auto text-sm mb-4">
            <code># Database Configuration</code><br>
            <code>DB_USER=postgres</code><br>
            <code>DB_PASSWORD=postgres</code><br>
            <code>DB_NAME=main_db</code><br>
            <code>DB_HOST=postgres</code><br>
            <code>DB_PORT=5432</code><br>
            <code><br></code>
            <code># Service Ports</code><br>
            <code>API_GATEWAY_PORT=5001</code><br>
            <code>AUTH_SERVICE_PORT=5003</code><br>
            <code>DATABASE_SERVICE_PORT=5002</code><br>
            <code>HISTORY_SERVICE_PORT=5008</code><br>
            <code>TEXT_TO_SQL_SERVICE_PORT=5006</code><br>
            <code>VISUALIZATION_SERVICE_PORT=5007</code><br>
            <code>FRONTEND_PORT=4173</code><br>
            <code><br></code>
            <code># GroqCloud API Key</code><br>
            <code>GROQ_API_KEY=your_groq_api_key_here</code><br>
            <code><br></code>
            <code># JWT Secret</code><br>
            <code>JWT_SECRET=your_secure_jwt_secret_here</code>
            <code># SMTP Configuration</code>
            <code>SMTP_EMAIL</code>
            <code>SMTP_PASSWORD</code>
            <code>SMTP_HOST</code>
            <code>SMTP_PORT</code>
        </div>

        <div class="bg-yellow-50 p-4 rounded-lg">
            <h3 class="text-lg font-medium mb-2 text-yellow-800">Important Security Note</h3>
            <p class="text-gray-700">
                For production environments, use strong, unique passwords and JWT secrets. Never commit your <code class="bg-gray-100 px-1 rounded text-sm">.env</code> file to version control.
            </p>
        </div>
    </div>

    <!-- Running the System Section -->
    <div class="bg-white rounded-xl shadow p-6">
        <div class="flex items-center gap-3 mb-4">
            <Container class="w-8 h-8 text-red-600" />
            <h2 class="text-2xl font-semibold text-gray-800">Running the System</h2>
        </div>

        <p class="text-gray-600 mb-4">Start all services using Docker Compose:</p>

        <div class="bg-gray-800 text-green-400 p-4 rounded-md overflow-x-auto text-sm mb-4">
            <code>docker-compose up -d</code>
        </div>

        <p class="text-gray-600 mb-4">This command will:</p>
        <ul class="list-disc list-inside space-y-1 text-gray-600 ml-4">
            <li>Build all service images</li>
            <li>Create and configure the PostgreSQL database</li>
            <li>Start all microservices in detached mode</li>
            <li>Set up the necessary networking between containers</li>
        </ul>

        <p class="text-gray-600 mt-4">To check the status of your containers:</p>

        <div class="bg-gray-800 text-green-400 p-4 rounded-md overflow-x-auto text-sm mb-4">
            <code>docker-compose ps</code>
        </div>

        <p class="text-gray-600 mt-4">To view logs for all services or a specific service:</p>

        <div class="bg-gray-800 text-green-400 p-4 rounded-md overflow-x-auto text-sm mb-4">
            <code># All services</code><br>
            <code>docker-compose logs -f</code><br>
            <code><br></code>
            <code># Specific service</code><br>
            <code>docker-compose logs -f frontend</code>
        </div>
    </div>

    <!-- Accessing the Application Section -->
    <div class="bg-white rounded-xl shadow p-6">
        <div class="flex items-center gap-3 mb-4">
            <Globe class="w-8 h-8 text-blue-600" />
            <h2 class="text-2xl font-semibold text-gray-800">Accessing the Application</h2>
        </div>

        <p class="text-gray-600 mb-4">Once all services are running, you can access the application at:</p>

        <div class="bg-blue-50 p-4 rounded-lg text-center text-lg font-medium text-blue-800 mb-4">
            <a href="http://localhost:3000" target="_blank" rel="noopener noreferrer" class="hover:underline">http://localhost:4173</a>
        </div>

        <div class="grid md:grid-cols-2 gap-6">
            <div>
                <h3 class="text-lg font-medium mb-2 text-gray-700">Service Endpoints</h3>
                <ul class="space-y-1 text-gray-600">
                    <li><span class="font-medium">API Gateway:</span> http://localhost:5001</li>
                    <li><span class="font-medium">Auth Service:</span> http://localhost:5003</li>
                    <li><span class="font-medium">Database Service:</span> http://localhost:5002</li>
                    <li><span class="font-medium">History Service:</span> http://localhost:5008</li>
                    <li><span class="font-medium">Text-to-SQL Service:</span> http://localhost:5006</li>
                    <li><span class="font-medium">Visualisation Service:</span> http://localhost:5007</li>
                </ul>
            </div>
            <div>
                <h3 class="text-lg font-medium mb-2 text-gray-700">Initial Access</h3>
                <p class="text-gray-600">
                    When you first access the system, you'll need to:
                </p>
                <ol class="list-decimal list-inside space-y-1 text-gray-600 mt-2">
                    <li>Register a new user account</li>
                    <li>Create or connect to a sample database</li>
                    <li>Start querying and visualizing data</li>
                </ol>
            </div>
        </div>
    </div>

    <!-- Troubleshooting Section -->
    <div class="bg-white rounded-xl shadow p-6">
        <div class="flex items-center gap-3 mb-4">
            <AlertTriangle class="w-8 h-8 text-yellow-600" />
            <h2 class="text-2xl font-semibold text-gray-800">Troubleshooting</h2>
        </div>

        <div class="space-y-4">
            <div>
                <h3 class="text-lg font-medium text-gray-700">Service not starting</h3>
                <p class="text-gray-600">
                    Check the logs for the specific service using <code class="bg-gray-100 px-1 rounded text-sm">docker-compose logs -f service_name</code>. Common issues include port conflicts or missing environment variables.
                </p>
            </div>

            <div>
                <h3 class="text-lg font-medium text-gray-700">Database connection issues</h3>
                <p class="text-gray-600">
                    Ensure PostgreSQL container is running with <code class="bg-gray-100 px-1 rounded text-sm">docker-compose ps postgres</code>. If needed, try restarting the service: <code class="bg-gray-100 px-1 rounded text-sm">docker-compose restart postgres</code>.
                </p>
            </div>

            <div>
                <h3 class="text-lg font-medium text-gray-700">Text-to-SQL not working</h3>
                <p class="text-gray-600">
                    Verify your GroqCloud API key is correctly set in the <code class="bg-gray-100 px-1 rounded text-sm">.env</code> file and that the text-to-SQL service can access it.
                </p>
            </div>

            <div>
                <h3 class="text-lg font-medium text-gray-700">Complete reset</h3>
                <p class="text-gray-600">
                    If you need a fresh start, you can completely reset the system with:
                </p>
                <div class="bg-gray-800 text-green-400 p-3 rounded-md overflow-x-auto text-sm mt-2">
                    <code>docker-compose down -v</code><br>
                    <code>docker-compose up -d</code>
                </div>
                <p class="text-gray-600 mt-2">
                    Note: This will delete all data stored in the PostgreSQL database.
                </p>
            </div>
        </div>
    </div>
</div>

<div class="border-t border-gray-200 mt-12 pt-8 text-gray-500 text-sm">
    <p>Last updated: April 20, 2025</p>
</div>