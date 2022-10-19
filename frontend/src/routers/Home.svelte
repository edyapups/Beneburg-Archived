<script lang="ts">
    import {getUsers} from "../lib/API";
    import User from "../lib/User.svelte";

    const users_promise: Promise<User[]> = getUsers();
</script>

{#await users_promise}
    <div class="container">
        <p>loading...</p>
    </div>
{:then users}
    <div class="container">
        <div class="row row-cols-5">
            {#each users as user}
                <div class="btn hover col shadow rounded text-center pt-4">
                    <User {user} />
                </div>
            {/each}
        </div>
    </div>
{:catch error}
    <div class="container"><p>{error.message}</p></div>
{/await}

<style>
    .hover:hover {
        transform:scale(1.05);
        -webkit-filter: brightness(70%);
        transition: all 0.1s;
    }
</style>
