<script lang="ts">
    import {
        Mail,
        Github,
        Linkedin,
        GraduationCap,
        Code,
        Database,
        BarChart3
    } from 'lucide-svelte';
    import NavbarDocumentation from '$lib/components/NavbarDocumentation.svelte';
    import Footer from '$lib/components/Footer.svelte';

    const team = [
        {
            name: 'Aruzhan Yeleukeshova',
            role: 'Project Manager, Designer',
            description: 'Led the development of the user interface, implemented responsive design and user authentication.',
            image: '/profile-aruzhan.jpg',
            icon: Code,
            iconColor: 'text-blue-500',
            links: {
                github: " https://github.com/Aaizya",
                linkedin: " https://www.linkedin.com/in/aruzhan-yeleukeshova-a8b7912b1?utm_source=share&utm_campaign=share_via&utm_content=profile&utm_medium=ios_app",
                email: "aruzhan1471@gmail.com",
            }
        },
        {
            name: 'Sayat Adilkhanov',
            role: 'Fullstack Developer',
            description: 'Architected the Go microservices, implemented the SQL generation engine and database connectivity layers.',
            image: '/profile-sayat.jpg',
            icon: Database,
            iconColor: 'text-green-500',
            links: {
                github: "",
                linkedin: "",
                email: "",
            }
        },
        {
            name: 'Guzal Mazhitova',
            role: 'ML Engineer, Researcher',
            description: 'Made research of the NLP & LLM models for text-to-SQL conversion, trained visualization recommendation algorithms.',
            image: '/profile-guzal.jpg',
            icon: BarChart3,
            iconColor: 'text-purple-500',
            links: {
                github: " https://github.com/Guzalm",
                linkedin: " https://www.linkedin.com/in/%D0%B3%D1%83%D0%B7%D0%B0%D0%BB%D1%8C-%D0%BC%D0%B0%D0%B6%D0%B8%D1%82%D0%BE%D0%B2%D0%B0-b36b29360?utm_source=share&utm_campaign=share_via&utm_content=profile&utm_medium=ios_app ",
                email: "guzql09@gmail.com",
            }
        }
    ];

    const supervisor = {
        name: 'Almas Myrzakhmetov',
        title: 'Associate Professor',
        department: 'Department of Computer Engineering',
        institution: 'Astana IT University',
        email: 'a.myrzakhmetov@astanait.edu.kz'
    };
</script>

<NavbarDocumentation />

<div class="container mx-auto px-6 py-16 w-11/12">
    <h1 class="text-4xl font-bold text-gray-800 mb-4">Meet the Team</h1>
    <p class="text-gray-600 max-w-2xl mb-12">The minds behind the text-to-SQL visualization system, developed as part of our diploma project.</p>

    <!-- Team -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-10 mb-20">
        {#each team as member}
            <div class="bg-white shadow rounded-2xl overflow-hidden hover:shadow-lg transition">
                <div class="relative h-36 bg-gradient-to-r from-blue-50 to-indigo-50 flex items-center justify-center">
                    <svelte:component this={member.icon} class="w-12 h-12 {member.iconColor}" />
                </div>
                <div class="p-6">
                    <div class="flex items-center gap-4 mb-4">
                        <img src={member.image} alt={member.name} class="w-12 h-12 rounded-full object-cover border" />
                        <div>
                            <h3 class="text-lg font-semibold text-gray-800">{member.name}</h3>
                            <p class="text-sm text-blue-600 font-medium">{member.role}</p>
                        </div>
                    </div>
                    <p class="text-sm text-gray-600 mb-4">{member.description}</p>
                    <div class="flex gap-4 text-gray-400">
                        <a href="#" aria-label="Email" class="hover:text-blue-600"><Mail class="w-5 h-5" /></a>
                        <a href="#" aria-label="GitHub" class="hover:text-gray-800"><Github class="w-5 h-5" /></a>
                        <a href="#" aria-label="LinkedIn" class="hover:text-blue-700"><Linkedin class="w-5 h-5" /></a>
                    </div>
                </div>
            </div>
        {/each}
    </div>

    <!-- Supervisor -->
    <h2 class="text-3xl font-bold text-gray-800 mb-6">Project Supervisor</h2>
    <div class="bg-white rounded-xl shadow p-6 flex flex-col sm:flex-row gap-6 items-center sm:items-start max-w-3xl mb-20">
        <div class="w-20 h-20 rounded-full bg-gradient-to-r from-amber-100 to-amber-200 flex items-center justify-center flex-shrink-0">
            <GraduationCap class="w-10 h-10 text-amber-600" />
        </div>
        <div>
            <h3 class="text-xl font-semibold text-gray-800 mb-1">{supervisor.name}</h3>
            <p class="text-amber-600 font-medium mb-1">{supervisor.title}</p>
            <p class="text-gray-600 mb-4">{supervisor.department}, {supervisor.institution}</p>
            <a href="mailto:{supervisor.email}" class="inline-flex items-center gap-2 text-gray-600 hover:text-amber-600 transition">
                <Mail class="w-4 h-4" /> {supervisor.email}
            </a>
        </div>
    </div>

    <!-- Contact -->
    <h2 class="text-3xl font-bold text-gray-800 mb-6">Get in Touch</h2>
    <form class="bg-white rounded-xl shadow p-6 max-w-3xl space-y-4">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
                <label for="name" class="block text-sm font-medium text-gray-700">Name</label>
                <input id="name" type="text" class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500" />
            </div>
            <div>
                <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
                <input id="email" type="email" class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500" />
            </div>
        </div>
        <div>
            <label for="subject" class="block text-sm font-medium text-gray-700">Subject</label>
            <input id="subject" type="text" class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500" />
        </div>
        <div>
            <label for="message" class="block text-sm font-medium text-gray-700">Message</label>
            <textarea id="message" rows="4" class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:ring-blue-500 focus:border-blue-500"></textarea>
        </div>
        <button type="submit" class="inline-flex items-center gap-2 bg-blue-600 text-white px-4 py-2 rounded-md hover:bg-blue-700 transition">
            <Mail class="w-4 h-4" />
            Send Message
        </button>
    </form>
</div>

<Footer />
