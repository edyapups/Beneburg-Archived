<script>
    import Router from 'svelte-spa-router';
    import Home from "./routers/Home.svelte";
    import NotFound from "./routers/NotFound.svelte";
    import {location} from 'svelte-spa-router'
    import Header from "./lib/Header.svelte";

    const routes = {
        '/': Home,
        '*': NotFound,
    }
</script>

<body>
    <Header />
    <Router {routes}/>
    <div class="container mt-4"><p>The current page is: {$location}</p></div>
</body>


