<script>
	import '../app.css';
    import Navbar from "../lib/Navbar.svelte";
    import Footer from "../lib/Footer.svelte";
    import {onMount} from "svelte";
	let { children } = $props();

    function getAuthToken() {
        const cookies = document.cookie.split('; ');
        const tokenCookie = cookies.find(row => row.startsWith('token='));
        return tokenCookie ? tokenCookie.split('=')[1] : null;
    }

    onMount(() => {
        if (getAuthToken) {
            
        }
    })
</script>

<Navbar/>
{@render children()}
<Footer/>