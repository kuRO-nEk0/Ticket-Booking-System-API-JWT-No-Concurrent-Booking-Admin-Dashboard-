<script>
    import { createEventDispatcher } from 'svelte';
    import { api } from '../api.js';

    export let booking;
    export let allowCancel = false;

    const dispatch = createEventDispatcher();
    let isCancelling = false;
    let confirmCancel = false;
    let error = '';

    async function handleCancel() {
        if (!confirmCancel) {
            confirmCancel = true;
            setTimeout(() => confirmCancel = false, 3000);
            return;
        }

        isCancelling = true;
        error = '';
        try {
            await api.cancelBooking(booking.ID);
            dispatch('canceled', booking.ID);
        } catch (err) {
            error = err.message;
            isCancelling = false;
            confirmCancel = false;
        }
    }
</script>

<div class="max-w-sm w-full mx-auto bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
    <!-- Ticket Header -->
    <div class="bg-black p-6 text-center text-white">
        <h3 class="text-sm font-semibold tracking-widest uppercase opacity-80">Virtual Ticket</h3>
        <h2 class="text-2xl font-bold mt-2 truncate">{booking.Event.Title}</h2>
    </div>

    <!-- Ticket Body -->
    <div class="p-6 relative">
        <!-- Tear off effect -->
        <div class="absolute -left-3 top-1/2 w-6 h-6 bg-gray-50 rounded-full border-r border-gray-200 transform -translate-y-1/2"></div>
        <div class="absolute -right-3 top-1/2 w-6 h-6 bg-gray-50 rounded-full border-l border-gray-200 transform -translate-y-1/2"></div>
        <div class="absolute left-0 right-0 top-1/2 border-t-2 border-dashed border-gray-200 transform -translate-y-1/2"></div>

        <div class="text-center mb-6">
            <p class="text-xs text-gray-400 uppercase tracking-widest mb-1">Purchased On</p>
            <p class="text-sm font-medium text-gray-700">
                {new Date(booking.CreatedAt).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })} at {new Date(booking.CreatedAt).toLocaleTimeString(undefined, { hour: '2-digit', minute: '2-digit' })}
            </p>
        </div>

        <div class="mb-10 text-center">
            <p class="text-sm text-gray-500 uppercase tracking-widest mb-1">Event Date</p>
            <p class="font-semibold text-gray-900">{new Date(booking.Event.Date).toLocaleDateString(undefined, { weekday: 'long', year: 'numeric', month: 'long', day: 'numeric' })}</p>
        </div>

        <div class="mt-10 flex justify-between items-end text-center">
            <div class="flex-1">
                <p class="text-xs text-gray-500 uppercase tracking-widest mb-1">Seat</p>
                <p class="text-3xl font-bold text-gray-900">{booking.Seat.SeatNumber}</p>
            </div>
            <div class="flex-1 border-l border-gray-200">
                <p class="text-xs text-gray-500 uppercase tracking-widest mb-1">Status</p>
                <p class="text-lg font-bold text-green-600">CONFIRMED</p>
            </div>
        </div>
    </div>

    <!-- Barcode -->
    <div class="bg-gray-50 p-6 border-t border-gray-200 text-center">
        <!-- CSS Fake Barcode -->
        <div class="h-12 w-full flex justify-center gap-1 mx-auto mb-2 opacity-80 mix-blend-multiply">
            {#each Array(30) as _}
                <div class="bg-black h-full" style="width: {Math.random() * 4 + 1}px;"></div>
            {/each}
        </div>
        <p class="text-xs text-gray-400 font-mono tracking-widest uppercase">{booking.ID.split('-')[0]}-{booking.SeatID.split('-')[0]}</p>
        
        {#if allowCancel}
            <div class="mt-6 pt-4 border-t border-gray-200">
                {#if error}
                    <p class="text-red-500 text-xs mb-2">{error}</p>
                {/if}
                <button 
                    class="w-full py-2 text-xs font-bold rounded transition-colors {confirmCancel ? 'bg-red-600 text-white hover:bg-red-700' : 'bg-white border border-gray-300 text-gray-500 hover:text-red-600 hover:border-red-600'}"
                    on:click={handleCancel}
                    disabled={isCancelling}
                >
                    {isCancelling ? 'Cancelling...' : (confirmCancel ? 'Click again to confirm cancel' : 'Cancel Ticket')}
                </button>
            </div>
        {/if}
    </div>
</div>
