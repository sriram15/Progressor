<script lang="ts">
    import { onMount } from 'svelte';
    import { GetProfiles, SwitchProfile } from '@bindings/github.com/sriram15/progressor-todo-app/progressorapp';
    import type { Profile } from '@bindings/github.com/sriram15/progressor-todo-app/internal/profile/models';
    import { Fa } from 'svelte-fa';
    import { faUserCircle } from '@fortawesome/free-solid-svg-icons';
    import { profileStore } from '@/stores/profileStore';

    let profiles = $state<Profile[]>([]);
    let loading = $state(true);
    let error = $state<string | null>(null);

    onMount(async () => {
        try {
            const fetchedProfiles = await GetProfiles();
            if (fetchedProfiles && fetchedProfiles.length > 0) {
                profiles = fetchedProfiles;
            }
        } catch (e: any) {
            error = e.message;
            console.error("Failed to load profiles:", e);
        } finally {
            loading = false;
        }
    });

    async function handleProfileSelect(profile: Profile) {
        try {
            await SwitchProfile(profile.id);
            profileStore.setActiveProfile(profile);
        } catch (e: any) {
            error = `Failed to switch to profile: ${e.message}`;
            console.error(error);
        }
    }
</script>

<div class="container h-full mx-auto flex justify-center items-center">
    <div class="card p-8 w-full max-w-md shadow-xl text-center">
        <h2 class="h2 mb-6">Select a Profile</h2>

        {#if loading}
            <div class="skeleton h-48 w-full"></div>
        {:else if error}
            <div class="alert variant-filled-error">{error}</div>
        {:else if profiles.length > 0}
            <ul class="space-y-4">
                {#each profiles as profile (profile.id)}
                    <li>
                        <button
                            class="btn btn-lg w-full variant-soft-primary"
                            on:click={() => handleProfileSelect(profile)}
                        >
                            <Fa icon={faUserCircle} class="mr-2" />
                            <span>{profile.name}</span>
                        </button>
                    </li>
                {/each}
            </ul>
        {:else}
             <p>No profiles found.</p>
        {/if}
    </div>
</div>
