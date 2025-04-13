<script lang="ts"></script>

<h1 class="text-2xl font-bold mb-4">Visualisation Service</h1>

<p class="mb-4 text-gray-700">
    The Visualisation Service is responsible for generating SVG diagrams based on SQL result sets. It uses Go and gRPC to communicate with other services.
</p>

<h2 class="text-xl font-semibold mt-6 mb-2">Endpoints</h2>
<ul class="list-disc ml-6 text-gray-700 space-y-2">
    <li><code>/generate-graph</code>: Accepts SQL results and returns an SVG graph.</li>
</ul>

<h2 class="text-xl font-semibold mt-6 mb-2">Tech Stack</h2>
<ul class="list-disc ml-6 text-gray-700">
    <li>Go</li>
    <li>gRPC</li>
    <li>SVG rendering</li>
</ul>

<h2 class="text-xl font-semibold mt-6 mb-2">Usage Example</h2>
<pre class="bg-gray-100 p-4 rounded-lg text-sm overflow-x-auto">
GenerateChart

Request:
    &#123;
        "columns": ["user_id", "order_id"],
        "rows": [
            [1, 101],
            [2, 102],
            [3, 103]
        ]
    &#125;

Response:
&lt;svg&gt;...&lt;/svg&gt;
</pre>
