const API_BASE = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';

// Helper function to handle standard fetch requests with automatic JSON parsing and error handling
async function fetchAPI(endpoint, options = {}) {
    const token = localStorage.getItem('jwt_token');
    
    const headers = {
        'Content-Type': 'application/json',
        ...options.headers,
    };

    // Attach JWT if it exists
    if (token) {
        headers['Authorization'] = `Bearer ${token}`;
    }

    const response = await fetch(`${API_BASE}${endpoint}`, {
        ...options,
        headers,
    });

    const data = await response.json();

    if (!response.ok) {
        throw new Error(data.error || 'Something went wrong');
    }

    return data;
}

export const api = {
    login: (email, password) => fetchAPI('/auth/login', {
        method: 'POST',
        body: JSON.stringify({ email, password })
    }),
    
    register: (email, password) => fetchAPI('/auth/register', {
        method: 'POST',
        body: JSON.stringify({ email, password })
    }),

    getEvents: () => fetchAPI('/events'),
    
    getEventSeats: (eventId) => fetchAPI(`/events/${eventId}/seats`),
    
    bookSeat: (eventId, seatId) => fetchAPI('/bookings', {
        method: 'POST',
        body: JSON.stringify({ event_id: eventId, seat_id: seatId })
    }),
    getMyBookings: () => fetchAPI('/bookings'),
    cancelBooking: (bookingId) => fetchAPI(`/bookings/${bookingId}`, { method: 'DELETE' }),
    createEvent: (eventData) => fetchAPI('/events', { 
        method: 'POST', 
        body: JSON.stringify(eventData) 
    }),
    deleteEvent: (eventId) => fetchAPI(`/events/${eventId}`, { method: 'DELETE' })
};
