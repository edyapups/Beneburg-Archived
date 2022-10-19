<script lang="ts">
    import type {User} from "./API";

    export let user: User
</script>

<h1 class="fw-normal">{user.name}
    {#if user.age}
        <span class="text-muted fs-5 text-center align-middle">{user.age}</span>
    {/if}
</h1>

<p class="text-muted">
    {#if user.username }
        @{user.username}
    {:else}
        {user.telegram_id}
    {/if}
</p>
