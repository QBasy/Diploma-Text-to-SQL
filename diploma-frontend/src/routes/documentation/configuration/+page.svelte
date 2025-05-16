<script lang="ts">
    import { ChevronRight, HammerIcon, Server, Globe, ShieldAlert, Database } from 'lucide-svelte';

    let POSTGRES_USER = "";
    let POSTGRES_PASSWORD = "";
    let POSTGRES_DB = "";
    let API_GATEWAY_PORT = "";
</script>

<h1 class="text-3xl font-bold mb-4 text-gray-800">⚙️ Configuration</h1>
<p class="text-gray-600 mb-8">
    The system uses environment variables and Docker Compose for configuration across all microservices. This page details how to set up and configure each component of the system.
</p>

<div class="space-y-12">
    <!-- Environment Variables Section -->
    <div class="bg-white rounded-xl shadow p-6">
        <div class="flex items-center gap-3 mb-4">
            <HammerIcon class="w-8 h-8 text-gray-700" />
            <h2 class="text-2xl font-semibold text-gray-800">Environment Variables</h2>
        </div>

        <p class="text-gray-700 mb-4">
            The system uses a centralized <code class="bg-gray-100 px-1 rounded text-sm">.env</code> file at the root directory that sets configuration for all services. Create this file before starting the system.
        </p>

        <div class="bg-gray-50 p-4 rounded-lg mb-6">
            <h3 class="text-lg font-medium mb-3 text-gray-700">Common Environment Variables</h3>
            <div class="overflow-x-auto">
                <table class="min-w-full text-sm">
                    <thead>
                    <tr class="bg-gray-100">
                        <th class="px-4 py-2 text-left font-medium text-gray-700">Variable</th>
                        <th class="px-4 py-2 text-left font-medium text-gray-700">Description</th>
                        <th class="px-4 py-2 text-left font-medium text-gray-700">Default</th>
                    </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-200">
                    <tr>
                        <td class="px-4 py-3 font-mono text-blue-600">DB_HOST</td>
                        <td class="px-4 py-3">PostgreSQL database host</td>
                        <td class="px-4 py-3 font-mono">postgres</td>
                    </tr>
                    <tr>
                        <td class="px-4 py-3 font-mono text-blue-600">DB_PORT</td>
                        <td class="px-4 py-3">PostgreSQL database port</td>
                        <td class="px-4 py-3 font-mono">5432</td>
                    </tr>
                    <tr>
                        <td class="px-4 py-3 font-mono text-blue-600">DB_USER</td>
                        <td class="px-4 py-3">PostgreSQL username</td>
                        <td class="px-4 py-3 font-mono">postgres</td>
                    </tr>
                    <tr>
                        <td class="px-4 py-3 font-mono text-blue-600">DB_PASSWORD</td>
                        <td class="px-4 py-3">PostgreSQL password</td>
                        <td class="px-4 py-3 font-mono">postgres</td>
                    </tr>
                    <tr>
                        <td class="px-4 py-3 font-mono text-blue-600">DB_NAME</td>
                        <td class="px-4 py-3">PostgreSQL database name</td>
                        <td class="px-4 py-3 font-mono">mvpdiploma</td>
                    </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div class="mt-6">
            <h3 class="text-lg font-medium mb-3 text-gray-700">Service-Specific Variables</h3>
            <div class="accordion space-y-3">
                <!-- API Gateway -->
                <div class="border border-gray-200 rounded-lg overflow-hidden">
                    <div class="bg-gray-50 px-4 py-3 flex justify-between items-center cursor-pointer">
                        <h4 class="text-md font-medium text-gray-700">API Gateway</h4>
                        <ChevronRight class="w-5 h-5 text-gray-500" />
                    </div>
                    <div class="p-4 bg-white">
                        <div class="overflow-x-auto">
                            <table class="min-w-full text-sm">
                                <thead>
                                <tr class="bg-gray-100">
                                    <th class="px-4 py-2 text-left font-medium text-gray-700">Variable</th>
                                    <th class="px-4 py-2 text-left font-medium text-gray-700">Description</th>
                                    <th class="px-4 py-2 text-left font-medium text-gray-700">Default</th>
                                </tr>
                                </thead>
                                <tbody class="divide-y divide-gray-200">
                                <tr>
                                    <td class="px-4 py-3 font-mono text-blue-600">API_GATEWAY_PORT</td>
                                    <td class="px-4 py-3">Gateway service port</td>
                                    <td class="px-4 py-3 font-mono">5001</td>
                                </tr>
                                <tr>
                                    <td class="px-4 py-3 font-mono text-blue-600">AUTH_SERVICE_URL</td>
                                    <td class="px-4 py-3">Auth service endpoint</td>
                                    <td class="px-4 py-3 font-mono">http://auth-service:5003</td>
                                </tr>
                                <tr>
                                    <td class="px-4 py-3 font-mono text-blue-600">DB_SERVICE_URL</td>
                                    <td class="px-4 py-3">Database service endpoint</td>
                                    <td class="px-4 py-3 font-mono">http://database-service:5002</td>
                                </tr>
                                <tr>
                                    <td class="px-4 py-3 font-mono text-blue-600">HISTORY_SERVICE_URL</td>
                                    <td class="px-4 py-3">History service endpoint</td>
                                    <td class="px-4 py-3 font-mono">http://history-service:5008</td>
                                </tr>
                                <tr>
                                    <td class="px-4 py-3 font-mono text-blue-600">TEXT_TO_SQL_SERVICE_URL</td>
                                    <td class="px-4 py-3">Text-to-SQL service endpoint</td>
                                    <td class="px-4 py-3 font-mono">http://text-to-sql-service:5006</td>
                                </tr>
                                <tr>
                                    <td class="px-4 py-3 font-mono text-blue-600">CORS_ALLOWED_ORIGINS</td>
                                    <td class="px-4 py-3">Origins allowed for CORS</td>
                                    <td class="px-4 py-3 font-mono">http://localhost:4173</td>
                                </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>

                <!-- Auth Service -->
                <div class="border border-gray-200 rounded-lg overflow-hidden">
                    <div class="bg-gray-50 px-4 py-3 flex justify-between items-center cursor-pointer">
                        <h4 class="text-md font-medium text-gray-700">Auth Service</h4>
                        <ChevronRight class="w-5 h-5 text-gray-500" />
                    </div>
                    <div class="p-4 bg-white">
                        <div class="overflow-x-auto">
                            <table class="min-w-full text-sm">
                                <thead>
                                <tr class="bg-gray-100">
                                    <th class="px-4 py-2 text-left font-medium text-gray-700">Variable</th>
                                    <th class="px-4 py-2 text-left font-medium text-gray-700">Description</th>
                                    <th class="px-4 py-2 text-left font-medium text-gray-700">Default</th>
                                </tr>
                                </thead>
                                <tbody class="divide-y divide-gray-200">
                                <tr>
                                    <td class="px-4 py-3 font-mono text-blue-600">AUTH_PORT</td>
                                    <td class="px-4 py-3">Auth service port</td>
                                    <td class="px-4 py-3 font-mono">5003</td>
                                </tr>
                                <tr>
                                    <td class="px-4 py-3 font-mono text-blue-600">JWT_SECRET</td>
                                    <td class="px-4 py-3">Secret for JWT token generation</td>
                                    <td class="px-4 py-3 font-mono">[randomly generated]</td>
                                </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>

                <!-- Text-to-SQL Service -->
                <div class="border border-gray-200 rounded-lg overflow-hidden">
                    <div class="bg-gray-50 px-4 py-3 flex justify-between items-center cursor-pointer">
                        <h4 class="text-md font-medium text-gray-700">Text-to-SQL Service</h4>
                        <ChevronRight class="w-5 h-5 text-gray-500" />
                    </div>
                    <div class="p-4 bg-white">
                        <div class="overflow-x-auto">
                            <table class="min-w-full text-sm">
                                <thead>
                                <tr class="bg-gray-100">
                                    <th class="px-4 py-2 text-left font-medium text-gray-700">Variable</th>
                                    <th class="px-4 py-2 text-left font-medium text-gray-700">Description</th>
                                    <th class="px-4 py-2 text-left font-medium text-gray-700">Default</th>
                                </tr>
                                </thead>
                                <tbody class="divide-y divide-gray-200">
                                <tr>
                                    <td class="px-4 py-3 font-mono text-blue-600">TEXT_TO_SQL_PORT</td>
                                    <td class="px-4 py-3">Text-to-SQL service port</td>
                                    <td class="px-4 py-3 font-mono">5006</td>
                                </tr>
                                <tr>
                                    <td class="px-4 py-3 font-mono text-blue-600">GROQ_API_KEY</td>
                                    <td class="px-4 py-3">API key for Groq Cloud</td>
                                    <td class="px-4 py-3 font-mono">[required]</td>
                                </tr>
                                <tr>
                                    <td class="px-4 py-3 font-mono text-blue-600">LLAMA_MODEL</td>
                                    <td class="px-4 py-3">LLM model to use for text-to-SQL</td>
                                    <td class="px-4 py-3 font-mono">llama-3</td>
                                </tr>
                                </tbody>
                            </table>
                        </div>
                    </div>
                </div>

                <!-- Other Services -->
                <div class="border border-gray-200 rounded-lg overflow-hidden">
                    <div class="bg-gray-50 px-4 py-3 flex justify-between items-center cursor-pointer">
                        <h4 class="text-md font-medium text-gray-700">Other Services</h4>
                        <ChevronRight class="w-5 h-5 text-gray-500" />
                    </div>
                    <div class="p-4 bg-white">
                        <p class="text-gray-700 mb-3">
                            The following services have their own configuration variables. Check the example <code class="bg-gray-100 px-1 rounded text-sm">.env.example</code> file for the complete list.
                        </p>
                        <ul class="list-disc list-inside space-y-1 text-gray-600">
                            <li>Database Service (PORT, SQLite file locations)</li>
                            <li>History Service (retention period, batch size)</li>
                            <li>Visualization Service (SVG settings, image size limits)</li>
                        </ul>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Docker Compose Section -->
    <div class="bg-white rounded-xl shadow p-6">
        <div class="flex items-center gap-3 mb-4">
            <Server class="w-8 h-8 text-blue-600" />
            <h2 class="text-2xl font-semibold text-gray-800">Docker Compose Setup</h2>
        </div>

        <p class="text-gray-600 mb-4">
            The system uses Docker Compose to orchestrate all microservices. The configuration is split into development and production files.
        </p>

        <div class="bg-blue-50 p-4 rounded-lg mb-6">
            <h3 class="text-lg font-medium mb-3 text-blue-800">docker-compose.yml</h3>
            <pre class="bg-gray-800 text-green-400 p-3 rounded-md overflow-x-auto text-sm"><code>version: '3.8'

services:
  postgres:
    image: postgres:14
    environment:
      - POSTGRES_USER=${POSTGRES_USER}
                - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
                - POSTGRES_DB=${POSTGRES_DB}
                volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  api-gateway:
    build: ./api-gateway
    ports:
      - "${API_GATEWAY_PORT}:${API_GATEWAY_PORT}"
    environment:
      - PORT=${API_GATEWAY_PORT}
                - AUTH_SERVICE_URL=AUTH_SERVICE_URL
                - DB_SERVICE_URL=DB_SERVICE_URL
                - HISTORY_SERVICE_URL=HISTORY_SERVICE_URL
                - TEXT_TO_SQL_SERVICE_URL=TEXT_TO_SQL_SERVICE_URL
                depends_on:
      - auth-service
      - database-service
      - history-service
      - text-to-sql-service

  auth-service:
    build: ./auth-service
    environment:
      - PORT=
                - JWT_SECRET=
                - POSTGRES_HOST=
                - POSTGRES_PORT=
                - POSTGRES_USER=
                - POSTGRES_PASSWORD=
                - POSTGRES_DB=
                depends_on:
      postgres:
        condition: service_healthy

  # Additional services defined similarly...

volumes:
  postgres_data:</code></pre>
        </div>

        <div class="mt-6">
            <h3 class="text-lg font-medium mb-3 text-gray-700">Service Dependencies</h3>
            <p class="text-gray-600 mb-4">
                The Docker Compose file establishes the following dependency relationships:
            </p>
            <ul class="list-disc list-inside space-y-1 text-gray-600">
                <li>API Gateway depends on all other services</li>
                <li>Database Service depends on PostgreSQL</li>
                <li>Auth Service depends on PostgreSQL</li>
                <li>History Service depends on PostgreSQL</li>
                <li>Visualization Service is standalone</li>
                <li>Text-to-SQL Service is standalone but requires GROQ API key</li>
            </ul>
        </div>
    </div>

    <!-- Network Configuration -->
    <div class="bg-white rounded-xl shadow p-6">
        <div class="flex items-center gap-3 mb-4">
            <Globe class="w-8 h-8 text-green-600" />
            <h2 class="text-2xl font-semibold text-gray-800">Network Configuration</h2>
        </div>

        <p class="text-gray-600 mb-4">
            The system uses Docker's default bridge network for inter-service communication. Services can communicate using their service names as hostnames.
        </p>

        <div class="bg-gray-50 p-4 rounded-lg mb-6">
            <h3 class="text-lg font-medium mb-3 text-gray-700">Port Mapping</h3>
            <div class="overflow-x-auto">
                <table class="min-w-full text-sm">
                    <thead>
                    <tr class="bg-gray-100">
                        <th class="px-4 py-2 text-left font-medium text-gray-700">Service</th>
                        <th class="px-4 py-2 text-left font-medium text-gray-700">Internal Port</th>
                        <th class="px-4 py-2 text-left font-medium text-gray-700">External Port</th>
                    </tr>
                    </thead>
                    <tbody class="divide-y divide-gray-200">
                    <tr>
                        <td class="px-4 py-3">API Gateway</td>
                        <td class="px-4 py-3 font-mono">8000</td>
                        <td class="px-4 py-3 font-mono">8000</td>
                    </tr>
                    <tr>
                        <td class="px-4 py-3">Auth Service</td>
                        <td class="px-4 py-3 font-mono">8001</td>
                        <td class="px-4 py-3 font-mono">Not exposed</td>
                    </tr>
                    <tr>
                        <td class="px-4 py-3">Database Service</td>
                        <td class="px-4 py-3 font-mono">8002</td>
                        <td class="px-4 py-3 font-mono">Not exposed</td>
                    </tr>
                    <tr>
                        <td class="px-4 py-3">History Service</td>
                        <td class="px-4 py-3 font-mono">8003</td>
                        <td class="px-4 py-3 font-mono">Not exposed</td>
                    </tr>
                    <tr>
                        <td class="px-4 py-3">Text-to-SQL Service</td>
                        <td class="px-4 py-3 font-mono">8004</td>
                        <td class="px-4 py-3 font-mono">Not exposed</td>
                    </tr>
                    <tr>
                        <td class="px-4 py-3">Visualization Service (gRPC)</td>
                        <td class="px-4 py-3 font-mono">50051</td>
                        <td class="px-4 py-3 font-mono">Not exposed</td>
                    </tr>
                    <tr>
                        <td class="px-4 py-3">PostgreSQL</td>
                        <td class="px-4 py-3 font-mono">5432</td>
                        <td class="px-4 py-3 font-mono">5432</td>
                    </tr>
                    <tr>
                        <td class="px-4 py-3">Frontend (SvelteKit)</td>
                        <td class="px-4 py-3 font-mono">3000</td>
                        <td class="px-4 py-3 font-mono">3000</td>
                    </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div class="mt-6">
            <h3 class="text-lg font-medium mb-3 text-gray-700">CORS Configuration</h3>
            <p class="text-gray-600 mb-4">
                Cross-Origin Resource Sharing (CORS) is configured in the API Gateway. Set the <code class="bg-gray-100 px-1 rounded text-sm">CORS_ALLOWED_ORIGINS</code> environment variable to allow requests from your frontend.
            </p>
        </div>
    </div>

    <!-- Security Configuration -->
    <div class="bg-white rounded-xl shadow p-6">
        <div class="flex items-center gap-3 mb-4">
            <ShieldAlert class="w-8 h-8 text-red-600" />
            <h2 class="text-2xl font-semibold text-gray-800">Security Configuration</h2>
        </div>

        <p class="text-gray-600 mb-4">
            Security is implemented across multiple layers of the system.
        </p>

        <div class="space-y-4">
            <div class="bg-gray-50 p-4 rounded-lg">
                <h3 class="text-lg font-medium mb-2 text-gray-700">JWT Authentication</h3>
                <p class="text-gray-600">
                    The system uses JWT tokens for authentication. Configure the following environment variables:
                </p>
                <ul class="list-disc list-inside space-y-1 mt-2 text-gray-600">
                    <li><code class="bg-gray-100 px-1 rounded text-sm">JWT_SECRET</code>: A strong secret key for token signing</li>
                    <li><code class="bg-gray-100 px-1 rounded text-sm">JWT_EXPIRATION</code>: Token expiration time in hours</li>
                </ul>
            </div>

            <div class="bg-gray-50 p-4 rounded-lg">
                <h3 class="text-lg font-medium mb-2 text-gray-700">Database Security</h3>
                <p class="text-gray-600">
                    Database access is controlled through:
                </p>
                <ul class="list-disc list-inside space-y-1 mt-2 text-gray-600">
                    <li>Strong PostgreSQL credentials in production</li>
                    <li>SQLite databases created with restricted permissions</li>
                    <li>Query validation to prevent SQL injection</li>
                </ul>
            </div>

            <div class="bg-gray-50 p-4 rounded-lg">
                <h3 class="text-lg font-medium mb-2 text-gray-700">API Security</h3>
                <p class="text-gray-600">
                    API endpoints are protected using:
                </p>
                <ul class="list-disc list-inside space-y-1 mt-2 text-gray-600">
                    <li>JWT token validation middleware</li>
                    <li>Rate limiting configured in the API Gateway</li>
                    <li>Input validation on all endpoints</li>
                </ul>
            </div>
        </div>
    </div>
</div>

<div class="bg-blue-50 p-6 rounded-xl mt-8">
    <h2 class="text-xl font-semibold mb-4 text-blue-800">Configuration Files</h2>
    <p class="text-gray-700 mb-4">
        The repository includes example configuration files that you should copy and modify:
    </p>

    <ul class="list-disc list-inside space-y-2 text-gray-700">
        <li><code class="bg-gray-100 px-1 rounded text-sm">.env.example</code> → <code class="bg-gray-100 px-1 rounded text-sm">.env</code></li>
        <li><code class="bg-gray-100 px-1 rounded text-sm">docker-compose.yml</code> (for development)</li>
        <li><code class="bg-gray-100 px-1 rounded text-sm">docker-compose.prod.yml</code> (for production)</li>
    </ul>

    <div class="mt-6 flex items-center gap-2 text-blue-600 font-medium">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
        </svg>
        <span>Never commit sensitive values like API keys or passwords to version control.</span>
    </div>
</div>

<div class="border-t border-gray-200 mt-12 pt-8 text-gray-500 text-sm">
    <p>Last updated: April 20, 2025</p>
</div>