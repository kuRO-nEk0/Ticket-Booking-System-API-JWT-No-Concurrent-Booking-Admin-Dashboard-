<script>
    import { onMount } from 'svelte';
    import { api } from '../api.js';
    import VirtualTicket from './VirtualTicket.svelte';

    let bookings = [];
    let loading = true;
    let error = '';

    onMount(async () => {
        try {
            bookings = await api.getMyBookings();
        } catch (err) {
            error = err.message;
        } finally {
            loading = false;
        }
    });

    function handleCanceled(event) {
        const canceledId = event.detail;
        bookings = bookings.filter(b => b.ID !== canceledId);
    }
</script>

<div class="w-full">
    <h2 class="text-2xl font-bold mb-8">My Booking History</h2>
    
    {#if loading}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            {#each Array(2) as _}
                <div class="max-w-sm w-full mx-auto bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden animate-pulse">
                    <div class="bg-gray-300 h-24 w-full"></div>
                    <div class="p-6 h-48 bg-gray-50 flex flex-col justify-between">
                        <div class="h-4 bg-gray-200 w-1/3 mx-auto rounded"></div>
                        <div class="flex justify-between">
                            <div class="h-8 bg-gray-200 w-1/4 rounded"></div>
                            <div class="h-8 bg-gray-200 w-1/4 rounded"></div>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
    {:else if error}
        <div class="mb-4 p-3 rounded bg-red-50 border border-red-200 text-red-700 text-center text-sm">
            {error}
        </div>
    {:else if bookings.length === 0}
        <div class="text-center py-12 border border-dashed border-gray-300 rounded-lg bg-gray-50">
            <p class="text-gray-500 mb-2">You haven't booked any tickets yet.</p>
        </div>
    {:else}
        <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
            {#each bookings as booking}
                <VirtualTicket {booking} allowCancel={true} on:canceled={handleCanceled} />
            {/each}
        </div>
    {/if}
</div>
