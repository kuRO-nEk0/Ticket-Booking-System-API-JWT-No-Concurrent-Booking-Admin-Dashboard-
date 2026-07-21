<script>
  import { onMount } from 'svelte';
  import LoginForm from './lib/components/LoginForm.svelte';
  import EventList from './lib/components/EventList.svelte';
  import SeatMap from './lib/components/SeatMap.svelte';
  import TicketHistory from './lib/components/TicketHistory.svelte';
  import VirtualTicket from './lib/components/VirtualTicket.svelte';
  import AdminPanel from './lib/components/AdminPanel.svelte';
  
  let currentRoute = 'login'; // 'login', 'events', 'booking', 'history', 'ticket', 'admin'
  let currentEventId = null;
  let currentBooking = null;
  let token = null;
  let userEmail = null;

  onMount(() => {
    token = localStorage.getItem('jwt_token');
    if (token) {
      try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        userEmail = payload.email;
      } catch (e) {
        console.error("Failed to parse token");
      }
      currentRoute = 'events';
    }
  });

  function navigate(route, eventId = null) {
    currentRoute = route;
    currentEventId = eventId;
    if (route !== 'ticket') {
        currentBooking = null;
    }
  }

  function handleLogout() {
    localStorage.removeItem('jwt_token');
    token = null;
    userEmail = null;
    currentRoute = 'login';
  }

  function onBookingSuccess(e) {
      currentBooking = e.detail.booking;
      currentRoute = 'ticket';
  }
</script>

<main class="min-h-screen bg-gray-50 flex flex-col items-center py-12 px-4 font-sans text-gray-900">
  <div class="w-full max-w-5xl">
    <!-- Navbar -->
    <header class="flex justify-between items-center mb-8 pb-4 border-b border-gray-200">
      <h1><button class="text-2xl font-bold tracking-tight hover:text-gray-600 transition-colors focus:outline-none" on:click={() => token ? navigate('events') : navigate('login')}>BookTickets</button></h1>
      <nav>
        {#if !token}
          {#if currentRoute !== 'login'}
            <button class="text-sm font-medium text-gray-600 hover:text-gray-900" on:click={() => navigate('login')}>Sign In</button>
          {/if}
        {:else}
          <div class="flex items-center gap-6">
            <button class="text-sm font-medium {currentRoute === 'events' ? 'text-black font-bold' : 'text-gray-600'} hover:text-gray-900" on:click={() => navigate('events')}>Events</button>
            <button class="text-sm font-medium {currentRoute === 'history' ? 'text-black font-bold' : 'text-gray-600'} hover:text-gray-900" on:click={() => navigate('history')}>My Tickets</button>
            {#if userEmail === 'tmarked4l@gmail.com'}
              <button class="text-sm font-medium {currentRoute === 'admin' ? 'text-indigo-600 font-bold' : 'text-indigo-600'} hover:text-indigo-800" on:click={() => navigate('admin')}>Admin Dashboard</button>
            {/if}
            <button class="text-sm font-medium text-red-600 hover:text-red-700" on:click={handleLogout}>Sign Out</button>
          </div>
        {/if}
      </nav>
    </header>

    <!-- Main Content Area -->
    <div class="bg-white border border-gray-200 rounded-xl p-8 min-h-[500px] flex flex-col items-center shadow-sm relative">
      {#if currentRoute === 'booking'}
        <button on:click={() => navigate('events')} class="absolute top-6 left-6 text-sm text-gray-500 hover:text-gray-900 transition-colors">&larr; Back to Events</button>
      {/if}
      
      {#if currentRoute === 'ticket'}
        <button on:click={() => navigate('history')} class="absolute top-6 left-6 text-sm text-gray-500 hover:text-gray-900 transition-colors">&larr; Back to My Tickets</button>
      {/if}

      {#if currentRoute === 'login'}
        <div class="text-center w-full max-w-md mx-auto">
          <h2 class="text-3xl font-bold mb-2">Welcome</h2>
          <p class="text-gray-500 mb-8">Sign in to book your tickets.</p>
          <LoginForm on:authSuccess={() => { 
            token = localStorage.getItem('jwt_token'); 
            try {
              const payload = JSON.parse(atob(token.split('.')[1]));
              userEmail = payload.email;
            } catch (e) {
              console.error("Failed to parse token");
            }
            navigate('events'); 
          }} />
        </div>
      {:else if currentRoute === 'events'}
        <div class="w-full">
          <h2 class="text-2xl font-bold mb-6">Upcoming Events</h2>
          <EventList on:eventSelected={(e) => navigate('booking', e.detail.eventId)} />
        </div>
      {:else if currentRoute === 'booking'}
        <div class="w-full">
          <h2 class="text-2xl font-bold mb-2 text-center">Select Your Seat</h2>
          <p class="text-gray-500 text-sm text-center mb-8">Click on an available seat to book your ticket.</p>
          <SeatMap eventId={currentEventId} on:bookingSuccess={onBookingSuccess} />
        </div>
      {:else if currentRoute === 'history'}
        <div class="w-full">
          <TicketHistory />
        </div>
      {:else if currentRoute === 'ticket'}
        <div class="w-full py-8">
            <h2 class="text-2xl font-bold mb-8 text-center text-green-600 flex items-center justify-center gap-2">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path></svg>
                Booking Confirmed!
            </h2>
            <VirtualTicket booking={currentBooking} />
        </div>
      {:else if currentRoute === 'admin'}
        <div class="w-full">
          <AdminPanel />
        </div>
      {/if}
    </div>
  </div>
</main>
