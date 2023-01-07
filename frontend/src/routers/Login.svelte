<script>
    import {params, push} from "svelte-spa-router";
    import {setToken} from "../lib/Token.js";
</script>

{#if params.token }
    <script>
        import {setToken} from "../lib/Token.js";
        import {params, push} from "svelte-spa-router";

        setToken(params.token)
        push("/")
    </script>
{/if}

https://t.me/BeneburgBot