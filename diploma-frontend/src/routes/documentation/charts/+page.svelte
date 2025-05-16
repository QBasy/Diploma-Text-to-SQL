<script lang="ts">
    import { BarChart3, LineChart, PieChart, ScatterChart } from 'lucide-svelte';
</script>

<h1 class="text-3xl font-bold mb-4 text-gray-800">📊 Chart Types</h1>
<p class="text-gray-600 mb-8">
    This visualization system supports multiple chart types that are automatically selected based on your data and SQL query patterns. All charts are rendered as SVG on the backend without JavaScript.
</p>

<div class="space-y-12">
    <!-- Bar Chart Section -->
    <div class="bg-white rounded-xl shadow p-6">
        <div class="flex items-center gap-3 mb-4">
            <BarChart3 class="w-8 h-8 text-blue-600" />
            <h2 class="text-2xl font-semibold text-gray-800">Bar Charts</h2>
        </div>

        <div class="grid md:grid-cols-2 gap-6">
            <div>
                <h3 class="text-lg font-medium mb-2 text-gray-700">When It's Used</h3>
                <ul class="list-disc list-inside space-y-1 text-gray-600">
                    <li>SQL queries using <code class="bg-gray-100 px-1 rounded text-sm">COUNT</code> or <code class="bg-gray-100 px-1 rounded text-sm">SUM</code> functions</li>
                    <li>Queries with <code class="bg-gray-100 px-1 rounded text-sm">GROUP BY</code> clauses</li>
                    <li>Default chart type when no specific pattern is detected</li>
                </ul>
                <h3 class="text-lg font-medium mt-4 mb-2 text-gray-700">Features</h3>
                <ul class="list-disc list-inside space-y-1 text-gray-600">
                    <li>Automatic scaling based on data values</li>
                    <li>Color-coded bars with responsive width</li>
                    <li>Value labels on top of each bar</li>
                    <li>Grid lines for better readability</li>
                </ul>
            </div>
            <div class="bg-gray-100 rounded-lg p-4 flex items-center justify-center">
                <svg width="300" height="200" viewBox="0 0 300 200">
                    <rect x="0" y="0" width="300" height="200" fill="#f8f9fa" />
                    <text x="150" y="20" text-anchor="middle" font-size="16" font-weight="bold" fill="#212529">Sample Bar Chart</text>

                    <!-- Grid lines -->
                    <line x1="40" y1="40" x2="280" y2="40" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />
                    <line x1="40" y1="70" x2="280" y2="70" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />
                    <line x1="40" y1="100" x2="280" y2="100" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />
                    <line x1="40" y1="130" x2="280" y2="130" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />
                    <line x1="40" y1="160" x2="280" y2="160" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />

                    <!-- Axes -->
                    <line x1="40" y1="160" x2="280" y2="160" stroke="#343a40" stroke-width="2" />
                    <line x1="40" y1="40" x2="40" y2="160" stroke="#343a40" stroke-width="2" />

                    <!-- Bars -->
                    <rect x="60" y="60" width="30" height="100" fill="#4361ee" rx="3" ry="3" />
                    <rect x="120" y="80" width="30" height="80" fill="#3a0ca3" rx="3" ry="3" />
                    <rect x="180" y="100" width="30" height="60" fill="#7209b7" rx="3" ry="3" />
                    <rect x="240" y="50" width="30" height="110" fill="#f72585" rx="3" ry="3" />

                    <!-- Labels -->
                    <text x="75" y="170" font-size="10" text-anchor="middle" fill="#495057">A</text>
                    <text x="135" y="170" font-size="10" text-anchor="middle" fill="#495057">B</text>
                    <text x="195" y="170" font-size="10" text-anchor="middle" fill="#495057">C</text>
                    <text x="255" y="170" font-size="10" text-anchor="middle" fill="#495057">D</text>
                </svg>
            </div>
        </div>

        <div class="mt-4 bg-blue-50 p-4 rounded-lg">
            <h3 class="text-lg font-medium mb-2 text-blue-800">Example SQL</h3>
            <pre class="bg-gray-800 text-green-400 p-3 rounded-md overflow-x-auto text-sm"><code>SELECT department, COUNT(*) as employee_count
FROM employees
GROUP BY department
ORDER BY employee_count DESC;</code></pre>
        </div>
    </div>

    <!-- Line Chart Section -->
    <div class="bg-white rounded-xl shadow p-6">
        <div class="flex items-center gap-3 mb-4">
            <LineChart class="w-8 h-8 text-green-600" />
            <h2 class="text-2xl font-semibold text-gray-800">Line Charts</h2>
        </div>

        <div class="grid md:grid-cols-2 gap-6">
            <div>
                <h3 class="text-lg font-medium mb-2 text-gray-700">When It's Used</h3>
                <ul class="list-disc list-inside space-y-1 text-gray-600">
                    <li>Queries containing <code class="bg-gray-100 px-1 rounded text-sm">DATE</code> or <code class="bg-gray-100 px-1 rounded text-sm">TIME</code> fields</li>
                    <li>Time series data visualization</li>
                    <li>Sequential or ordered data points</li>
                </ul>
                <h3 class="text-lg font-medium mt-4 mb-2 text-gray-700">Features</h3>
                <ul class="list-disc list-inside space-y-1 text-gray-600">
                    <li>Connected data points with smooth lines</li>
                    <li>Highlighted data points with tooltips</li>
                    <li>Semi-transparent area fill below the line</li>
                    <li>Adaptive X and Y axis scales</li>
                </ul>
            </div>
            <div class="bg-gray-100 rounded-lg p-4 flex items-center justify-center">
                <svg width="300" height="200" viewBox="0 0 300 200">
                    <rect x="0" y="0" width="300" height="200" fill="#f8f9fa" />
                    <text x="150" y="20" text-anchor="middle" font-size="16" font-weight="bold" fill="#212529">Sample Line Chart</text>

                    <!-- Grid lines -->
                    <line x1="40" y1="40" x2="280" y2="40" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />
                    <line x1="40" y1="70" x2="280" y2="70" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />
                    <line x1="40" y1="100" x2="280" y2="100" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />
                    <line x1="40" y1="130" x2="280" y2="130" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />
                    <line x1="40" y1="160" x2="280" y2="160" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />

                    <!-- Axes -->
                    <line x1="40" y1="160" x2="280" y2="160" stroke="#343a40" stroke-width="2" />
                    <line x1="40" y1="40" x2="40" y2="160" stroke="#343a40" stroke-width="2" />

                    <!-- Area fill -->
                    <path d="M 40 120 L 80 100 L 120 140 L 160 60 L 200 80 L 240 110 L 280 70 L 280 160 L 40 160 Z" fill="#4895ef" fill-opacity="0.1" />

                    <!-- Line -->
                    <path d="M 40 120 L 80 100 L 120 140 L 160 60 L 200 80 L 240 110 L 280 70" fill="none" stroke="#4895ef" stroke-width="3" stroke-linecap="round" stroke-linejoin="round" />

                    <!-- Points -->
                    <circle cx="40" cy="120" r="5" fill="#4895ef" stroke="#fff" stroke-width="2" />
                    <circle cx="80" cy="100" r="5" fill="#4895ef" stroke="#fff" stroke-width="2" />
                    <circle cx="120" cy="140" r="5" fill="#4895ef" stroke="#fff" stroke-width="2" />
                    <circle cx="160" cy="60" r="5" fill="#4895ef" stroke="#fff" stroke-width="2" />
                    <circle cx="200" cy="80" r="5" fill="#4895ef" stroke="#fff" stroke-width="2" />
                    <circle cx="240" cy="110" r="5" fill="#4895ef" stroke="#fff" stroke-width="2" />
                    <circle cx="280" cy="70" r="5" fill="#4895ef" stroke="#fff" stroke-width="2" />

                    <!-- Labels -->
                    <text x="40" y="175" font-size="10" text-anchor="middle" fill="#495057">Jan</text>
                    <text x="80" y="175" font-size="10" text-anchor="middle" fill="#495057">Feb</text>
                    <text x="120" y="175" font-size="10" text-anchor="middle" fill="#495057">Mar</text>
                    <text x="160" y="175" font-size="10" text-anchor="middle" fill="#495057">Apr</text>
                    <text x="200" y="175" font-size="10" text-anchor="middle" fill="#495057">May</text>
                    <text x="240" y="175" font-size="10" text-anchor="middle" fill="#495057">Jun</text>
                    <text x="280" y="175" font-size="10" text-anchor="middle" fill="#495057">Jul</text>
                </svg>
            </div>
        </div>

        <div class="mt-4 bg-green-50 p-4 rounded-lg">
            <h3 class="text-lg font-medium mb-2 text-green-800">Example SQL</h3>
            <pre class="bg-gray-800 text-green-400 p-3 rounded-md overflow-x-auto text-sm"><code>SELECT date_trunc('month', order_date) as month,
       SUM(total_amount) as revenue
FROM orders
WHERE order_date >= '2023-01-01'
GROUP BY month
ORDER BY month;</code></pre>
        </div>
    </div>

    <!-- Pie Chart Section -->
    <div class="bg-white rounded-xl shadow p-6">
        <div class="flex items-center gap-3 mb-4">
            <PieChart class="w-8 h-8 text-purple-600" />
            <h2 class="text-2xl font-semibold text-gray-800">Pie Charts</h2>
        </div>

        <div class="grid md:grid-cols-2 gap-6">
            <div>
                <h3 class="text-lg font-medium mb-2 text-gray-700">When It's Used</h3>
                <ul class="list-disc list-inside space-y-1 text-gray-600">
                    <li>Queries containing <code class="bg-gray-100 px-1 rounded text-sm">PERCENT</code> or <code class="bg-gray-100 px-1 rounded text-sm">RATIO</code> terms</li>
                    <li>Proportion or distribution analysis</li>
                    <li>Categorical data with parts-to-whole relationships</li>
                </ul>
                <h3 class="text-lg font-medium mt-4 mb-2 text-gray-700">Features</h3>
                <ul class="list-disc list-inside space-y-1 text-gray-600">
                    <li>Auto-colored segments with percentage labels</li>
                    <li>Legend with category names and values</li>
                    <li>Optimized segment spacing and sizing</li>
                    <li>Small segments handling more than 3%) </li>
                </ul>
            </div>
            <div class="bg-gray-100 rounded-lg p-4 flex items-center justify-center">
                <svg width="300" height="200" viewBox="0 0 300 200">
                    <rect x="0" y="0" width="300" height="200" fill="#f8f9fa" />
                    <text x="150" y="20" text-anchor="middle" font-size="16" font-weight="bold" fill="#212529">Sample Pie Chart</text>

                    <!-- Pie -->
                    <g transform="translate(100, 100)">
                        <!-- Segments -->
                        <path d="M 0 0 L 60 0 A 60 60 0 0 1 30 52 Z" fill="#4361ee" />
                        <path d="M 0 0 L 30 52 A 60 60 0 0 1 -30 52 Z" fill="#3a0ca3" />
                        <path d="M 0 0 L -30 52 A 60 60 0 0 1 -60 0 Z" fill="#7209b7" />
                        <path d="M 0 0 L -60 0 A 60 60 0 0 1 0 -60 Z" fill="#f72585" />
                        <path d="M 0 0 L 0 -60 A 60 60 0 0 1 60 0 Z" fill="#4cc9f0" />

                        <!-- Labels -->
                        <text x="35" y="20" text-anchor="middle" font-size="12" font-weight="bold" fill="white">30%</text>
                        <text x="0" y="40" text-anchor="middle" font-size="12" font-weight="bold" fill="white">20%</text>
                        <text x="-35" y="20" text-anchor="middle" font-size="12" font-weight="bold" fill="white">20%</text>
                        <text x="-20" y="-20" text-anchor="middle" font-size="12" font-weight="bold" fill="white">15%</text>
                        <text x="20" y="-20" text-anchor="middle" font-size="12" font-weight="bold" fill="white">15%</text>
                    </g>

                    <!-- Legend -->
                    <rect x="180" y="60" width="12" height="12" fill="#4361ee" />
                    <text x="198" y="70" font-size="10" fill="#212529">Category A: 30%</text>

                    <rect x="180" y="80" width="12" height="12" fill="#3a0ca3" />
                    <text x="198" y="90" font-size="10" fill="#212529">Category B: 20%</text>

                    <rect x="180" y="100" width="12" height="12" fill="#7209b7" />
                    <text x="198" y="110" font-size="10" fill="#212529">Category C: 20%</text>

                    <rect x="180" y="120" width="12" height="12" fill="#f72585" />
                    <text x="198" y="130" font-size="10" fill="#212529">Category D: 15%</text>

                    <rect x="180" y="140" width="12" height="12" fill="#4cc9f0" />
                    <text x="198" y="150" font-size="10" fill="#212529">Category E: 15%</text>
                </svg>
            </div>
        </div>

        <div class="mt-4 bg-purple-50 p-4 rounded-lg">
            <h3 class="text-lg font-medium mb-2 text-purple-800">Example SQL</h3>
            <pre class="bg-gray-800 text-green-400 p-3 rounded-md overflow-x-auto text-sm"><code>SELECT product_category,
       SUM(sales_amount) / (SELECT SUM(sales_amount) FROM sales) * 100 as percentage
FROM sales
GROUP BY product_category
ORDER BY percentage DESC;</code></pre>
        </div>
    </div>

    <!-- Scatter Chart Section -->
    <div class="bg-white rounded-xl shadow p-6">
        <div class="flex items-center gap-3 mb-4">
            <ScatterChart class="w-8 h-8 text-yellow-600" />
            <h2 class="text-2xl font-semibold text-gray-800">Scatter Charts</h2>
        </div>

        <div class="grid md:grid-cols-2 gap-6">
            <div>
                <h3 class="text-lg font-medium mb-2 text-gray-700">When It's Used</h3>
                <ul class="list-disc list-inside space-y-1 text-gray-600">
                    <li>Queries returning at least two numeric columns</li>
                    <li>Correlation or relationship analysis</li>
                    <li>Data with X-Y coordinate pairs</li>
                </ul>
                <h3 class="text-lg font-medium mt-4 mb-2 text-gray-700">Features</h3>
                <ul class="list-disc list-inside space-y-1 text-gray-600">
                    <li>Individual data points with consistent styling</li>
                    <li>Automatically scaled X and Y axes</li>
                    <li>Grid lines for value approximation</li>
                    <li>Tooltip support with point value details</li>
                </ul>
            </div>
            <div class="bg-gray-100 rounded-lg p-4 flex items-center justify-center">
                <svg width="300" height="200" viewBox="0 0 300 200">
                    <rect x="0" y="0" width="300" height="200" fill="#ffffff" />
                    <text x="150" y="20" text-anchor="middle" font-size="16" font-weight="bold" fill="#212529">Sample Scatter Chart</text>

                    <!-- Grid lines -->
                    <line x1="40" y1="40" x2="280" y2="40" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />
                    <line x1="40" y1="70" x2="280" y2="70" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />
                    <line x1="40" y1="100" x2="280" y2="100" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />
                    <line x1="40" y1="130" x2="280" y2="130" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />
                    <line x1="40" y1="160" x2="280" y2="160" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />

                    <line x1="40" y1="40" x2="40" y2="160" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />
                    <line x1="100" y1="40" x2="100" y2="160" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />
                    <line x1="160" y1="40" x2="160" y2="160" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />
                    <line x1="220" y1="40" x2="220" y2="160" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />
                    <line x1="280" y1="40" x2="280" y2="160" stroke="#dee2e6" stroke-width="1" stroke-dasharray="5,5" />

                    <!-- Axes -->
                    <line x1="40" y1="160" x2="280" y2="160" stroke="#343a40" stroke-width="2" />
                    <line x1="40" y1="40" x2="40" y2="160" stroke="#343a40" stroke-width="2" />

                    <!-- Points -->
                    <circle cx="70" cy="60" r="6" fill="#4895ef" stroke="#ffffff" stroke-width="2" />
                    <circle cx="120" cy="90" r="6" fill="#4895ef" stroke="#ffffff" stroke-width="2" />
                    <circle cx="150" cy="130" r="6" fill="#4895ef" stroke="#ffffff" stroke-width="2" />
                    <circle cx="190" cy="70" r="6" fill="#4895ef" stroke="#ffffff" stroke-width="2" />
                    <circle cx="220" cy="110" r="6" fill="#4895ef" stroke="#ffffff" stroke-width="2" />
                    <circle cx="250" cy="50" r="6" fill="#4895ef" stroke="#ffffff" stroke-width="2" />
                    <circle cx="90" cy="140" r="6" fill="#4895ef" stroke="#ffffff" stroke-width="2" />

                    <!-- Axis Labels -->
                    <text x="160" y="190" text-anchor="middle" font-size="14" font-weight="bold" fill="#343a40">X</text>
                    <text x="10" y="100" text-anchor="middle" font-size="14" font-weight="bold" fill="#343a40" transform="rotate(-90 10,100)">Y</text>
                </svg>
            </div>
        </div>

        <div class="mt-4 bg-yellow-50 p-4 rounded-lg">
            <h3 class="text-lg font-medium mb-2 text-yellow-800">Example SQL</h3>
            <pre class="bg-gray-800 text-green-400 p-3 rounded-md overflow-x-auto text-sm"><code>SELECT product_name,
       price,
       units_sold
FROM products
JOIN sales ON products.id = sales.product_id
WHERE launch_date > '2023-01-01';</code></pre>
        </div>
    </div>
</div>

<div class="bg-blue-50 p-6 rounded-xl mt-8">
    <h2 class="text-xl font-semibold mb-4 text-blue-800">How Chart Type Is Determined</h2>
    <p class="text-gray-700 mb-4">
        The system automatically analyzes your SQL query and data to select the most appropriate chart type:
    </p>

    <ol class="list-decimal list-inside space-y-2 text-gray-700">
        <li>First, it examines your SQL query for specific keywords like <code class="bg-gray-100 px-1 rounded text-sm">COUNT</code>, <code class="bg-gray-100 px-1 rounded text-sm">DATE</code>, or <code class="bg-gray-100 px-1 rounded text-sm">PERCENT</code></li>
        <li>Then, it analyzes the structure of your returned data (number of numeric columns, data types)</li>
        <li>Based on these factors, it selects the optimal chart type for your data visualization</li>
        <li>If no clear pattern is detected, it defaults to a bar chart for simplicity and readability</li>
    </ol>

    <div class="mt-6 flex items-center gap-2 text-blue-600 font-medium">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
        </svg>
        <span>All charts are rendered as SVG on the backend without JavaScript dependencies, ensuring fast loading and consistent appearance.</span>
    </div>
</div>

<div class="border-t border-gray-200 mt-12 pt-8 text-gray-500 text-sm">
    <p>Last updated: April 20, 2025</p>
</div>