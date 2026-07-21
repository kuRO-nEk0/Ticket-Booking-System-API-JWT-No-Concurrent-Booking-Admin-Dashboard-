<script>
    import { createEventDispatcher, onMount } from 'svelte';
    import { api } from '../api.js';

    const dispatch = createEventDispatcher();
    
    let events = [];
    let loading = true;
    let error = '';

    onMount(async () => {
        try {
            const data = await api.getEvents();
            events = data;

        } catch (err) {
            error = err.message;
        } finally {
            loading = false;
        }
    });
</script>

<div class="w-full">
    {#if loading}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            {#each Array(4) as _}
                <div class="p-5 border border-gray-200 rounded-lg animate-pulse bg-white">
                    <div class="h-6 bg-gray-200 rounded w-3/4 mb-3"></div>
                    <div class="h-4 bg-gray-200 rounded w-full mb-2"></div>
                    <div class="h-4 bg-gray-200 rounded w-5/6 mb-6"></div>
                    <div class="flex justify-between items-center">
                        <div class="h-3 bg-gray-200 rounded w-1/4"></div>
                        <div class="h-3 bg-gray-200 rounded w-1/6"></div>
                    </div>
                </div>
            {/each}
        </div>
    {:else if error}
        <p class="text-red-600 text-sm text-center">{error}</p>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            {#each events as event}
                <button 
                    on:click={() => dispatch('eventSelected', { eventId: event.ID })}
                    class="group p-5 border border-gray-200 rounded-lg text-left hover:border-black hover:bg-gray-50 transition-all focus:outline-none"
                >
                    <h3 class="text-lg font-bold text-gray-900 mb-1">{event.Title}</h3>
                    <p class="text-gray-600 text-sm mb-4 line-clamp-2">{event.Description}</p>
                    <div class="flex justify-between items-center text-xs font-medium text-gray-500 group-hover:text-black">
                        <span>{new Date(event.Date).toLocaleDateString()}</span>
                        <span>Book &rarr;</span>
                    </div>
                </button>
            {/each}
        </div>
    {/if}
</div>
