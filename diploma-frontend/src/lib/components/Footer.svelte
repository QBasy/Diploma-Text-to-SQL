<script lang="ts">
    import { Github, Mail, FileText, Heart, Database, LineChart } from 'lucide-svelte';
</script>

<footer class="bg-white border-t border-gray-200 mt-12">
    <div class="container mx-auto px-4 py-8">
        <!-- Logo and tagline -->
        <div class="flex flex-col sm:flex-row justify-between items-center mb-6">
            <div class="flex items-center gap-2 mb-4 sm:mb-0">
                <LineChart class="w-5 h-5 text-blue-600" />
                <span class="text-lg font-bold text-blue-600">SQL Translation</span>
            </div>
            <div class="text-sm text-gray-500">
                Transform natural language into SQL with data visualization
            </div>
        </div>

        <!-- Navigation links -->
        <div class="flex flex-wrap justify-center md:justify-between gap-y-4 border-t border-gray-100 pt-6">
            <div class="flex flex-wrap justify-center gap-x-8 gap-y-2">
                <a href="/documentation" class="flex items-center gap-2 text-gray-600 hover:text-blue-600 transition">
                    <FileText class="w-4 h-4" />
                    <span>Documentation</span>
                </a>
                <a href="/contacts" class="flex items-center gap-2 text-gray-600 hover:text-blue-600 transition">
                    <Mail class="w-4 h-4" />
                    <span>Contact</span>
                </a>
                <a href="https://github.com/QBasy/Diploma-Text-to-SQL" target="_blank" rel="noopener noreferrer"
                   class="flex items-center gap-2 text-gray-600 hover:text-blue-600 transition">
                    <Github class="w-4 h-4" />
                    <span>GitHub</span>
                </a>
            </div>

            <div class="w-full md:w-auto text-center md:text-right">
                <p class="text-xs text-gray-400 flex items-center justify-center md:justify-end gap-1">
                    © 2024 SQL Generator. Built with <Heart class="w-3 h-3 text-red-500 fill-red-500" /> using
                    <span class="font-medium">Go, SvelteKit, GroqCloud & Llama Scout 4</span>
                </p>
            </div>
        </div>
    </div>
</footer>