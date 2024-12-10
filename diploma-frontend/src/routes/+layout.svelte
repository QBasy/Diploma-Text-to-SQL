<script>
	import '../app.css';
    import { getAuthToken } from '$lib/api';
    import Navbar from "../lib/Navbar.svelte";
    import Footer from "../lib/Footer.svelte";
    import {onMount} from "svelte";
	let { children } = $props();

    onMount(() => {
        const token = getAuthToken();
        if (token) {
            console.log('User is authenticated');
        } else {
            console.log('User is not authenticated');
        }
    });

</script>

<Navbar/>
{@render children()}
<Footer/>