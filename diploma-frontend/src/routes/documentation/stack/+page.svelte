<script lang="ts">
    import { BookOpen } from 'lucide-svelte';
</script>

<h1 class="text-3xl font-bold mb-4 text-gray-800">🛠️ Tech Stack</h1>
<p class="text-gray-600 mb-8">
    This automatic data visualization system is built with a modern microservices architecture, utilizing a variety of technologies for different components. Below is a detailed breakdown of each part of the stack.
</p>

<div class="space-y-12">
    <!-- Architecture Overview -->
    <div class="bg-white rounded-xl shadow p-6">
        <h2 class="text-2xl font-semibold mb-4 text-gray-800">System Architecture</h2>
        <div class="bg-gray-100 rounded-lg p-6">
            <svg viewBox="0 0 900 400" xmlns="http://www.w3.org/2000/svg">
                <!-- Background -->
                <rect width="900" height="400" fill="#f8f9fa" rx="10" ry="10" />

                <!-- Frontend Box -->
                <rect x="340" y="30" width="120" height="60" fill="#4361ee" rx="8" ry="8" />
                <text x="400" y="65" text-anchor="middle" font-size="16" fill="white" font-weight="bold">Frontend</text>
                <text x="400" y="85" text-anchor="middle" font-size="12" fill="white">SvelteKit</text>

                <!-- API Gateway Box -->
                <rect x="340" y="150" width="120" height="60" fill="#3a0ca3" rx="8" ry="8" />
                <text x="400" y="185" text-anchor="middle" font-size="16" fill="white" font-weight="bold">API Gateway</text>
                <text x="400" y="205" text-anchor="middle" font-size="12" fill="white">Go + Gin</text>

                <!-- Connection Lines -->
                <line x1="400" y1="90" x2="400" y2="150" stroke="#6c757d" stroke-width="2" stroke-dasharray="5,5" />

                <!-- Microservices Boxes -->
                <rect x="100" y="270" width="120" height="60" fill="#7209b7" rx="8" ry="8" />
                <text x="160" y="305" text-anchor="middle" font-size="14" fill="white" font-weight="bold">Auth Service</text>
                <text x="160" y="325" text-anchor="middle" font-size="12" fill="white">Go + Gin + GORM</text>

                <rect x="250" y="270" width="120" height="60" fill="#f72585" rx="8" ry="8" />
                <text x="310" y="295" text-anchor="middle" font-size="14" fill="white" font-weight="bold">Database Service</text>
                <text x="310" y="310" text-anchor="middle" font-size="12" fill="white">Go + Gin + GORM</text>
                <text x="310" y="325" text-anchor="middle" font-size="12" fill="white">+ SQLite (User DBs)</text>

                <rect x="400" y="270" width="120" height="60" fill="#4cc9f0" rx="8" ry="8" />
                <text x="460" y="305" text-anchor="middle" font-size="14" fill="white" font-weight="bold">History Service</text>
                <text x="460" y="325" text-anchor="middle" font-size="12" fill="white">Go + Gin + GORM</text>

                <rect x="550" y="270" width="120" height="60" fill="#06d6a0" rx="8" ry="8" />
                <text x="610" y="305" text-anchor="middle" font-size="14" fill="white" font-weight="bold">Text-to-SQL</text>
                <text x="610" y="325" text-anchor="middle" font-size="12" fill="white">Go + Llama/GroqCloud</text>

                <rect x="700" y="270" width="120" height="60" fill="#22577a" rx="8" ry="8" />
                <text x="760" y="305" text-anchor="middle" font-size="14" fill="white" font-weight="bold">Visualization</text>
                <text x="760" y="325" text-anchor="middle" font-size="12" fill="white">Go + gRPC + SVGo</text>

                <!-- Connection Lines to Microservices -->
                <line x1="400" y1="210" x2="160" y2="270" stroke="#6c757d" stroke-width="2" />
                <line x1="400" y1="210" x2="310" y2="270" stroke="#6c757d" stroke-width="2" />
                <line x1="400" y1="210" x2="460" y2="270" stroke="#6c757d" stroke-width="2" />
                <line x1="400" y1="210" x2="610" y2="270" stroke="#6c757d" stroke-width="2" />
                <line x1="400" y1="210" x2="760" y2="270" stroke="#6c757d" stroke-width="2" />

                <!-- gRPC Connection Between Database and Visualization -->
                <path d="M 370 290 C 500 220 650 220 750 270" stroke="#ffd166" stroke-width="3" fill="none" stroke-dasharray="10,5" />
                <text x="560" y="240" text-anchor="middle" font-size="12" fill="#ff6b6b" font-weight="bold">gRPC</text>

                <!-- PostgreSQL -->
                <rect x="150" y="370" width="140" height="30" fill="#0b525b" rx="5" ry="5" />
                <text x="220" y="390" text-anchor="middle" font-size="12" fill="white">PostgreSQL (Main DB)</text>

                <!-- DB Connections -->
                <line x1="310" y1="330" x2="220" y2="370" stroke="#6c757d" stroke-width="2" stroke-dasharray="5,5" />
                <line x1="160" y1="330" x2="220" y2="370" stroke="#6c757d" stroke-width="2" stroke-dasharray="5,5" />
                <line x1="460" y1="330" x2="220" y2="370" stroke="#6c757d" stroke-width="2" stroke-dasharray="5,5" />
            </svg>
        </div>
        <p class="text-gray-600 mt-4">
            The system follows a microservices architecture with an API Gateway handling client requests and routing them to appropriate specialized services.
        </p>
    </div>

    <!-- Backend Technologies -->
    <div class="bg-white rounded-xl shadow p-6">
        <h2 class="text-2xl font-semibold mb-4 text-gray-800">Backend Technologies</h2>

        <div class="space-y-6">
            <!-- Go Section -->
            <div class="border-l-4 border-blue-600 pl-4">
                <h3 class="text-xl font-medium text-gray-800">Go (Golang)</h3>
                <p class="text-gray-600 mb-2">
                    The core language used for all backend microservices, providing excellent performance, concurrency, and simplicity.
                </p>
                <ul class="list-disc list-inside text-gray-600 space-y-1">
                    <li>API Gateway: Handles routing, authentication verification, and load balancing</li>
                    <li>Microservices: Powers all five specialized services with consistent patterns</li>
                    <li>Version: 1.22.x (with modules support)</li>
                </ul>
            </div>

            <!-- Framework Section -->
            <div class="border-l-4 border-green-600 pl-4">
                <h3 class="text-xl font-medium text-gray-800">Gin Web Framework</h3>
                <p class="text-gray-600 mb-2">
                    A high-performance HTTP web framework used across most services for REST API endpoints.
                </p>
                <ul class="list-disc list-inside text-gray-600 space-y-1">
                    <li>Middleware support for logging, recovery, and CORS</li>
                    <li>Route grouping for API versioning and organization</li>
                    <li>Request validation and binding</li>
                </ul>
            </div>

            <!-- GORM Section -->
            <div class="border-l-4 border-purple-600 pl-4">
                <h3 class="text-xl font-medium text-gray-800">GORM</h3>
                <p class="text-gray-600 mb-2">
                    Object-Relational Mapping library for Go, simplifying database operations.
                </p>
                <ul class="list-disc list-inside text-gray-600 space-y-1">
                    <li>Used in auth-service, database-service, and history-service</li>
                    <li>Supports both PostgreSQL and SQLite connections</li>
                    <li>Provides migrations, relationships, and transaction management</li>
                </ul>
            </div>

            <!-- gRPC Section -->
            <div class="border-l-4 border-yellow-600 pl-4">
                <h3 class="text-xl font-medium text-gray-800">gRPC</h3>
                <p class="text-gray-600 mb-2">
                    High-performance RPC framework used for internal service communication.
                </p>
                <ul class="list-disc list-inside text-gray-600 space-y-1">
                    <li>Communication between database-service and visualization-service</li>
                    <li>Strongly typed contracts with Protocol Buffers</li>
                    <li>Efficient binary serialization for performance</li>
                </ul>
            </div>

            <!-- SVGo Section -->
            <div class="border-l-4 border-red-600 pl-4">
                <h3 class="text-xl font-medium text-gray-800">SVGo</h3>
                <p class="text-gray-600 mb-2">
                    Go library for SVG generation used in the visualization-service.
                </p>
                <ul class="list-disc list-inside text-gray-600 space-y-1">
                    <li>Package: "github.com/ajstarks/svgo"</li>
                    <li>Server-side generation of all chart types</li>
                    <li>No JavaScript dependencies for rendering charts</li>
                </ul>
            </div>

            <!-- Llama/GroqCloud Section -->
            <div class="border-l-4 border-indigo-600 pl-4">
                <h3 class="text-xl font-medium text-gray-800">Llama via GroqCloud</h3>
                <p class="text-gray-600 mb-2">
                    LLM integration for natural language to SQL translation.
                </p>
                <ul class="list-disc list-inside text-gray-600 space-y-1">
                    <li>Used in the text-to-sql-service</li>
                    <li>Translates natural language queries to valid SQL</li>
                    <li>API integration with GroqCloud for inference</li>
                </ul>
            </div>
        </div>
    </div>

    <!-- Frontend Technologies -->
    <div class="bg-white rounded-xl shadow p-6">
        <h2 class="text-2xl font-semibold mb-4 text-gray-800">Frontend Technology</h2>

        <div class="border-l-4 border-blue-600 pl-4">
            <h3 class="text-xl font-medium text-gray-800">SvelteKit</h3>
            <p class="text-gray-600 mb-2">
                A modern framework for building web applications with optimal performance.
            </p>
            <ul class="list-disc list-inside text-gray-600 space-y-1">
                <li>Component-based architecture for reusability</li>
                <li>SSR (Server-Side Rendering) support for SEO and performance</li>
                <li>TypeScript integration for type safety</li>
                <li>Efficient reactivity system with minimal boilerplate</li>
            </ul>
        </div>

        <div class="mt-6 border-l-4 border-teal-600 pl-4">
            <h3 class="text-xl font-medium text-gray-800">Styling & UI Components</h3>
            <p class="text-gray-600 mb-2">
                Modern styling solutions for responsive, beautiful interfaces.
            </p>
            <ul class="list-disc list-inside text-gray-600 space-y-1">
                <li>TailwindCSS for utility-first styling</li>
                <li>Lucide icons for consistent iconography</li>
                <li>Custom components for visualization controls and data display</li>
            </ul>
        </div>
    </div>

    <!-- Database Technologies -->
    <div class="bg-white rounded-xl shadow p-6">
        <h2 class="text-2xl font-semibold mb-4 text-gray-800">Database Systems</h2>

        <div class="grid md:grid-cols-2 gap-6">
            <div class="border-l-4 border-blue-600 pl-4">
                <h3 class="text-xl font-medium text-gray-800">PostgreSQL</h3>
                <p class="text-gray-600 mb-2">
                    Primary database for the system with advanced features.
                </p>
                <ul class="list-disc list-inside text-gray-600 space-y-1">
                    <li>Used for user accounts and system metadata</li>
                    <li>Stores query history and translation records</li>
                    <li>Authentication data and permissions</li>
                    <li>Advanced SQL features and performance</li>
                </ul>
            </div>

            <div class="border-l-4 border-green-600 pl-4">
                <h3 class="text-xl font-medium text-gray-800">SQLite</h3>
                <p class="text-gray-600 mb-2">
                    Lightweight database for user-specific data.
                </p>
                <ul class="list-disc list-inside text-gray-600 space-y-1">
                    <li>Created per-user for testing and development</li>
                    <li>Self-contained file-based database</li>
                    <li>Used to test NLQ to SQL translations</li>
                    <li>Perfect for small datasets and quick testing</li>
                </ul>
            </div>
        </div>
    </div>

    <!-- Infrastructure -->
    <div class="bg-white rounded-xl shadow p-6">
        <h2 class="text-2xl font-semibold mb-4 text-gray-800">Development & Deployment</h2>

        <div class="border-l-4 border-blue-600 pl-4 mb-6">
            <h3 class="text-xl font-medium text-gray-800">Docker & Containers</h3>
            <p class="text-gray-600 mb-2">
                Containerization for consistent development and deployment.
            </p>
            <ul class="list-disc list-inside text-gray-600 space-y-1">
                <li>Individual Dockerfiles for each service</li>
                <li>Docker Compose for local development</li>
                <li>Environment variable management</li>
                <li>Volume mapping for persistent data</li>
            </ul>
        </div>

        <div class="border-l-4 border-yellow-600 pl-4">
            <h3 class="text-xl font-medium text-gray-800">Development Workflow</h3>
            <p class="text-gray-600 mb-2">
                Tools and practices for efficient development.
            </p>
            <ul class="list-disc list-inside text-gray-600 space-y-1">
                <li>Git for version control</li>
                <li>Environment-specific configuration</li>
                <li>Hot-reloading for development</li>
            </ul>
        </div>
    </div>
</div>

<div class="bg-blue-50 p-6 rounded-xl mt-8">
    <h2 class="text-xl font-semibold mb-4 text-blue-800">Technology Selection Rationale</h2>
    <p class="text-gray-700 mb-4">
        The technologies in this stack were chosen with these principles in mind:
    </p>

    <ul class="list-disc list-inside space-y-2 text-gray-700">
        <li><strong>Performance</strong>: Go provides excellent runtime performance with low resource usage</li>
        <li><strong>Scalability</strong>: Microservice architecture allows independent scaling of components</li>
        <li><strong>Developer Experience</strong>: SvelteKit and Go provide excellent developer ergonomics</li>
        <li><strong>Separation of Concerns</strong>: Each service has a clear, focused responsibility</li>
        <li><strong>Maintainability</strong>: Consistent patterns and technologies across services</li>
    </ul>
</div>

<div class="border-t border-gray-200 mt-12 pt-8 text-gray-500 text-sm">
    <p>Last updated: April 20, 2025</p>
</div>
