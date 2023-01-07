<script>
    import Router from 'svelte-spa-router';
    import {wrap} from 'svelte-spa-router/wrap'
    import Home from "./routers/Home.svelte";
    import NotFound from "./routers/NotFound.svelte";
    import {location} from 'svelte-spa-router'
    import Header from "./lib/Header.svelte";

    const routes = {
        '/': Home,
        '/login/:token?': wrap({
            asyncComponent: () => import('./routers/Login.svelte'),
        }),
        '*': NotFound,
    }
</script>

<body>
    <Header />
    <Router {routes}/>
    <div class="container mt-4"><p>The current page is: {$location}</p></div>
</body>


