<script>
    import { createEventDispatcher } from 'svelte';
    import { api } from '../api.js';

    const dispatch = createEventDispatcher();

    let email = '';
    let password = '';
    let isRegistering = false;
    let error = '';
    let emailError = '';
    let passwordError = '';
    let loading = false;
    let showSuccessPopup = false;
    let showPassword = false;

    function validate() {
        emailError = '';
        passwordError = '';
        let isValid = true;

        if (!email.includes('@') || !email.includes('.')) {
            emailError = 'Please enter a valid email address.';
            isValid = false;
        }

        if (password.length < 6) {
            passwordError = 'Password must be at least 6 characters.';
            isValid = false;
        }

        return isValid;
    }

    async function handleSubmit() {
        if (!validate()) return;
        
        error = '';
        loading = true;
        try {
            if (isRegistering) {
                await api.register(email, password);
                
                // Show custom popup of affirmation
                showSuccessPopup = true;
                
                // Auto login after 1.5s delay to let them read the popup
                setTimeout(async () => {
                    try {
                        const data = await api.login(email, password);
                        localStorage.setItem('jwt_token', data.token);
                        showSuccessPopup = false;
                        dispatch('authSuccess');
                    } catch (e) {
                        error = e.message;
                        showSuccessPopup = false;
                    }
                }, 1500);

            } else {
                const data = await api.login(email, password);
                localStorage.setItem('jwt_token', data.token);
                dispatch('authSuccess');
            }
        } catch (err) {
            error = err.message;
            loading = false;
        }
    }
</script>

{#if showSuccessPopup}
    <!-- Custom Affirmation Popup overlay -->
    <div class="fixed inset-0 bg-gray-900/50 flex items-center justify-center z-50">
        <div class="bg-white rounded-lg p-6 shadow-xl max-w-sm w-full text-center">
            <div class="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-green-100 mb-4">
                <svg class="h-6 w-6 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                </svg>
            </div>
            <h3 class="text-lg font-bold text-gray-900 mb-1">Account Created!</h3>
            <p class="text-sm text-gray-500">Signing you in automatically...</p>
        </div>
    </div>
{/if}

<div class="w-full">
    <form on:submit|preventDefault={handleSubmit} class="space-y-5 text-left">
        <div>
            <label for="email" class="block text-sm font-medium text-gray-700">Email Address</label>
            <input 
                id="email" 
                type="email" 
                bind:value={email} 
                required 
                class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 text-gray-900 placeholder-gray-400 focus:border-black focus:outline-none focus:ring-1 focus:ring-black sm:text-sm"
                placeholder="you@example.com"
            />
            {#if emailError}
                <p class="text-red-500 text-xs mt-1">{emailError}</p>
            {/if}
        </div>

        <div>
            <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
            <div class="relative mt-1">
                <input 
                    id="password" 
                    type={showPassword ? "text" : "password"}
                    value={password} 
                    on:input={(e) => password = e.target.value}
                    required 
                    class="block w-full rounded-md border border-gray-300 px-3 py-2 pr-10 text-gray-900 placeholder-gray-400 focus:border-black focus:outline-none focus:ring-1 focus:ring-black sm:text-sm"
                    placeholder="••••••••"
                />
                <button 
                    type="button" 
                    class="absolute inset-y-0 right-0 flex items-center pr-3 text-gray-400 hover:text-gray-600 focus:outline-none"
                    on:click={() => showPassword = !showPassword}
                >
                    {#if showPassword}
                        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                        </svg>
                    {:else}
                        <svg class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.543 7 0 4.057-3.732 7-8.268 7-4.478 0-8.268-2.943-9.543-7z" />
                        </svg>
                    {/if}
                </button>
            </div>
            {#if passwordError}
                <p class="text-red-500 text-xs mt-1">{passwordError}</p>
            {/if}
        </div>

        {#if error}
            <p class="text-red-600 text-sm text-center">{error}</p>
        {/if}

        <button 
            type="submit" 
            disabled={loading}
            class="w-full rounded-md bg-black px-4 py-2 text-sm font-semibold text-white hover:bg-gray-800 focus:outline-none focus:ring-2 focus:ring-black focus:ring-offset-2 transition-colors disabled:opacity-50"
        >
            {loading && !showSuccessPopup ? 'Processing...' : (isRegistering ? 'Create Account' : 'Sign In')}
        </button>
    </form>

    <p class="mt-6 text-center text-sm text-gray-600">
        {isRegistering ? 'Already have an account?' : 'Need an account?'}
        <button 
            class="font-medium text-black hover:underline focus:outline-none"
            on:click={() => { isRegistering = !isRegistering; error = ''; }}
        >
            {isRegistering ? 'Sign in instead' : 'Create one now'}
        </button>
    </p>
</div>
