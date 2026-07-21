<script>
    import { onMount, createEventDispatcher } from 'svelte';
    import { api } from '../api.js';

    export let eventId;
    
    const dispatch = createEventDispatcher();

    let seats = [];
    let loading = true;
    let error = '';
    let isBooking = false;
    
    // Modal state
    let showConfirmModal = false;
    let selectedSeat = null;

    onMount(async () => {
        await loadSeats();
    });

    async function loadSeats() {
        try {
            const data = await api.getEventSeats(eventId);
            seats = data;
        } catch (err) {
            error = err.message;
        } finally {
            loading = false;
        }
    }

    function promptConfirm(seat) {
        if (seat.Status === 'booked') return;
        selectedSeat = seat;
        showConfirmModal = true;
    }

    function cancelConfirm() {
        showConfirmModal = false;
        selectedSeat = null;
    }

    async function confirmBook() {
        if (!selectedSeat) return;

        isBooking = true;
        error = '';
        const seatToBook = selectedSeat;
        showConfirmModal = false;

        try {
            const res = await api.bookSeat(eventId, seatToBook.ID);
            const seatIndex = seats.findIndex(s => s.ID === seatToBook.ID);
            if (seatIndex !== -1) {
                seats[seatIndex].Status = 'booked';
            }
            // Dispatch the successful booking object to the parent so it can display the virtual ticket
            dispatch('bookingSuccess', { booking: res.booking });
        } catch (err) {
            error = err.message;
            await loadSeats();
        } finally {
            isBooking = false;
            selectedSeat = null;
        }
    }
</script>

{#if showConfirmModal}
    <!-- Confirmation Modal Overlay -->
    <div class="fixed inset-0 bg-gray-900/50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 shadow-xl max-w-sm w-full">
            <h3 class="text-lg font-bold text-gray-900 mb-2">Confirm Booking</h3>
            <p class="text-gray-600 mb-6">Are you sure you want to book seat <span class="font-bold text-black">{selectedSeat?.SeatNumber}</span>?</p>
            <div class="flex justify-end gap-3">
                <button 
                    class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded hover:bg-gray-50 focus:outline-none"
                    on:click={cancelConfirm}
                >
                    Cancel
                </button>
                <button 
                    class="px-4 py-2 text-sm font-medium text-white bg-black rounded hover:bg-gray-800 focus:outline-none"
                    on:click={confirmBook}
                >
                    Confirm Booking
                </button>
            </div>
        </div>
    </div>
{/if}

<div class="w-full">
    {#if loading}
        <div class="flex items-center justify-center gap-6 mb-8">
            {#each Array(3) as _}
                <div class="w-16 h-4 bg-gray-200 rounded animate-pulse"></div>
            {/each}
        </div>
        <div class="w-full max-w-md mx-auto h-8 border-t-2 border-gray-200 mb-10"></div>
        <div class="grid grid-cols-5 sm:grid-cols-10 gap-2 justify-center max-w-2xl mx-auto">
            {#each Array(40) as _}
                <div class="w-10 h-10 rounded bg-gray-200 animate-pulse"></div>
            {/each}
        </div>
    {:else}
        {#if error}
            <div class="mb-4 p-3 rounded bg-red-50 border border-red-200 text-red-700 text-center text-sm">
                {error}
            </div>
        {/if}

        <div class="flex items-center justify-center gap-6 mb-8 text-sm text-gray-600">
            <div class="flex items-center gap-2">
                <div class="w-4 h-4 rounded border border-gray-300 bg-white"></div> Available
            </div>
            <div class="flex items-center gap-2">
                <div class="w-4 h-4 rounded bg-black"></div> Selected
            </div>
            <div class="flex items-center gap-2">
                <div class="w-4 h-4 rounded bg-gray-200 border border-gray-300"></div> Booked
            </div>
        </div>

        <!-- The Stage -->
        <div class="w-full max-w-md mx-auto h-8 border-t-2 border-gray-300 mb-10 flex items-center justify-center text-xs text-gray-400 uppercase tracking-widest">
            Stage
        </div>

        <!-- Seat Grid -->
        <div class="grid grid-cols-5 sm:grid-cols-10 gap-2 justify-center max-w-2xl mx-auto">
            {#each seats as seat}
                <button 
                    disabled={seat.Status === 'booked' || isBooking}
                    on:click={() => promptConfirm(seat)}
                    class="w-10 h-10 rounded flex items-center justify-center text-xs font-medium transition-colors border
                        {seat.Status === 'booked' 
                            ? 'bg-gray-100 border-gray-200 text-gray-400 cursor-not-allowed' 
                            : 'bg-white border-gray-300 text-gray-700 hover:border-black hover:text-black focus:outline-none focus:ring-2 focus:ring-black'
                        }
                    "
                >
                    {seat.SeatNumber}
                </button>
            {/each}
        </div>
    {/if}
</div>
