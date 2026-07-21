<script>
    import { createEventDispatcher, onMount } from 'svelte';
    import { api } from '../api.js';

    const dispatch = createEventDispatcher();

    let events = [];
    let eventsLoading = true;
    let eventsError = '';

    let title = '';
    let description = '';
    let date = '';
    let seatsCount = 20;

    let error = '';
    let success = '';
    let loading = false;

    // Validation Errors
    let titleError = '';
    let descriptionError = '';
    let dateError = '';
    let seatsError = '';

    onMount(async () => {
        await fetchEvents();
    });

    async function fetchEvents() {
        eventsLoading = true;
        try {
            events = await api.getEvents();
        } catch (err) {
            eventsError = err.message;
        } finally {
            eventsLoading = false;
        }
    }

    async function handleDelete(eventId) {
        if (!confirm('Are you sure you want to delete this event? This will also delete all associated seats and bookings.')) {
            return;
        }
        
        try {
            await api.deleteEvent(eventId);
            events = events.filter(e => e.ID !== eventId);
        } catch (err) {
            alert(err.message || 'Failed to delete event');
        }
    }

    function validate() {
        titleError = '';
        descriptionError = '';
        dateError = '';
        seatsError = '';
        let isValid = true;

        if (title.trim() === '') {
            titleError = 'Title is required.';
            isValid = false;
        }

        if (description.trim() === '') {
            descriptionError = 'Description is required.';
            isValid = false;
        }

        if (!date) {
            dateError = 'Date and time are required.';
            isValid = false;
        } else {
            const selectedDate = new Date(date);
            if (selectedDate <= new Date()) {
                dateError = 'Event date must be in the future.';
                isValid = false;
            }
        }

        if (seatsCount < 1 || seatsCount > 500) {
            seatsError = 'Seats count must be between 1 and 500.';
            isValid = false;
        }

        return isValid;
    }

    async function handleSubmit() {
        if (!validate()) return;

        loading = true;
        error = '';
        success = '';

        try {
            // Convert local datetime-local to RFC3339 for Go
            const isoDate = new Date(date).toISOString();
            
            const eventData = {
                title,
                description,
                date: isoDate,
                seats_count: parseInt(seatsCount)
            };

            await api.createEvent(eventData);
            
            success = `Successfully created "${title}" with ${seatsCount} seats!`;
            
            // Refresh events list
            await fetchEvents();
            
            // Reset form
            title = '';
            description = '';
            date = '';
            seatsCount = 20;

        } catch (err) {
            error = err.message || 'Failed to create event.';
        } finally {
            loading = false;
        }
    }
</script>

<div class="w-full max-w-2xl mx-auto">
    <div class="mb-8 border-b border-gray-200 pb-4 flex justify-between items-end">
        <div>
            <h2 class="text-3xl font-bold text-gray-900">Admin Dashboard</h2>
            <p class="text-gray-500 mt-1">Create and manage events</p>
        </div>
        <div class="bg-black text-white text-xs font-bold px-3 py-1 rounded-full uppercase tracking-widest">
            Admin Access
        </div>
    </div>

    <div class="bg-white border border-gray-200 rounded-xl p-8 shadow-sm">
        <h3 class="text-xl font-bold mb-6">Create New Event</h3>

        {#if success}
            <div class="mb-6 p-4 rounded bg-green-50 border border-green-200 text-green-700 flex items-center gap-3">
                <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
                {success}
            </div>
        {/if}

        {#if error}
            <div class="mb-6 p-4 rounded bg-red-50 border border-red-200 text-red-700">
                {error}
            </div>
        {/if}

        <form on:submit|preventDefault={handleSubmit} class="space-y-6">
            <div>
                <label for="title" class="block text-sm font-medium text-gray-700 mb-1">Event Title</label>
                <input 
                    id="title" 
                    type="text" 
                    bind:value={title} 
                    class="block w-full rounded-md border border-gray-300 px-4 py-2 text-gray-900 placeholder-gray-400 focus:border-black focus:outline-none focus:ring-1 focus:ring-black sm:text-sm"
                    placeholder="e.g., Tech Conference 2026"
                />
                {#if titleError}<p class="text-red-500 text-xs mt-1">{titleError}</p>{/if}
            </div>

            <div>
                <label for="description" class="block text-sm font-medium text-gray-700 mb-1">Description</label>
                <textarea 
                    id="description" 
                    rows="3"
                    bind:value={description} 
                    class="block w-full rounded-md border border-gray-300 px-4 py-2 text-gray-900 placeholder-gray-400 focus:border-black focus:outline-none focus:ring-1 focus:ring-black sm:text-sm"
                    placeholder="Brief details about the event..."
                ></textarea>
                {#if descriptionError}<p class="text-red-500 text-xs mt-1">{descriptionError}</p>{/if}
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div>
                    <label for="date" class="block text-sm font-medium text-gray-700 mb-1">Date & Time</label>
                    <input 
                        id="date" 
                        type="datetime-local" 
                        bind:value={date} 
                        class="block w-full rounded-md border border-gray-300 px-4 py-2 text-gray-900 focus:border-black focus:outline-none focus:ring-1 focus:ring-black sm:text-sm"
                    />
                    {#if dateError}<p class="text-red-500 text-xs mt-1">{dateError}</p>{/if}
                </div>

                <div>
                    <label for="seats" class="block text-sm font-medium text-gray-700 mb-1">Number of Seats (1-500)</label>
                    <input 
                        id="seats" 
                        type="number" 
                        bind:value={seatsCount} 
                        min="1"
                        max="500"
                        class="block w-full rounded-md border border-gray-300 px-4 py-2 text-gray-900 focus:border-black focus:outline-none focus:ring-1 focus:ring-black sm:text-sm"
                    />
                    {#if seatsError}<p class="text-red-500 text-xs mt-1">{seatsError}</p>{/if}
                </div>
            </div>

            <div class="pt-4 border-t border-gray-100">
                <button 
                    type="submit" 
                    disabled={loading}
                    class="w-full flex justify-center items-center rounded-md bg-black px-4 py-3 text-sm font-semibold text-white hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-black focus:ring-offset-2 transition-colors disabled:opacity-50"
                >
                    {#if loading}
                        <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        Creating Event & Generating Seats...
                    {:else}
                        Publish Event
                    {/if}
                </button>
            </div>
        </form>
    </div>

    <!-- Manage Existing Events -->
    <div class="mt-8 bg-white border border-gray-200 rounded-xl p-8 shadow-sm">
        <h3 class="text-xl font-bold mb-6">Manage Events</h3>
        
        {#if eventsLoading}
            <div class="animate-pulse space-y-4">
                {#each Array(3) as _}
                    <div class="h-16 bg-gray-100 rounded-md"></div>
                {/each}
            </div>
        {:else if eventsError}
            <div class="p-4 rounded bg-red-50 text-red-700">{eventsError}</div>
        {:else if events.length === 0}
            <div class="text-center py-8 text-gray-500 border border-dashed rounded-lg">No events found.</div>
        {:else}
            <div class="space-y-4">
                {#each events as event}
                    <div class="flex items-center justify-between p-4 border border-gray-200 rounded-md hover:border-gray-300 transition-colors">
                        <div>
                            <h4 class="font-bold text-gray-900">{event.Title}</h4>
                            <p class="text-xs text-gray-500">{new Date(event.Date).toLocaleDateString()} &middot; ID: <span class="font-mono">{event.ID.split('-')[0]}</span></p>
                        </div>
                        <button 
                            on:click={() => handleDelete(event.ID)}
                            class="text-sm font-medium text-red-600 hover:text-red-800 hover:bg-red-50 px-3 py-1.5 rounded transition-colors"
                        >
                            Delete
                        </button>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
</div>
