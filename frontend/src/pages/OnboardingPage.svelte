<script lang="ts">
    import { navigate } from 'svelte-routing';
    import { CreateProfile, SwitchProfile } from '@bindings/github.com/sriram15/progressor-todo-app/progressorapp';
    import { Profile, DBType } from '@bindings/github.com/sriram15/progressor-todo-app/internal/profile/models';

    let currentStep = $state(1);
    let profileName = $state('');
    let dbType = $state<DBType>(DBType.DBTypeSQLite);
    let tursoUrl = $state('');
    let tursoToken = $state('');

    let errorMessage = $state('');
    let isLoading = $state(false);

    function nextStep() {
        errorMessage = ''; // Clear previous errors
        if (currentStep === 1 && profileName.trim() === '') {
            errorMessage = 'Profile name cannot be empty.';
            return;
        }
        currentStep++;
    }

    function prevStep() {
        errorMessage = ''; // Clear previous errors
        currentStep--;
    }

    async function createAndSwitchProfile() {
        isLoading = true;
        errorMessage = '';
        try {
            const newProfile = new Profile({
                name: profileName,
                dbType: dbType,
                dbUrl: dbType === DBType.DBTypeTurso ? tursoUrl : undefined,
            });

            // Create the profile
            const createdProfile = await CreateProfile(newProfile, tursoToken);

            if (createdProfile && createdProfile.id) {
                // Switch to the newly created profile
                await SwitchProfile(createdProfile.id);
                // Navigate to the main application page (e.g., ProgressPage)
                navigate('/progress'); // Assuming /progress is your main app page
            } else {
                errorMessage = 'Failed to create profile: No profile ID returned.';
            }
        } catch (error: any) {
            console.error('Error creating or switching profile:', error);
            errorMessage = `Error: ${error.message || 'An unknown error occurred.'}`;
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="container h-full mx-auto flex justify-center items-center">
    <div class="card p-8 w-full max-w-md shadow-xl">
        <h2 class="h2 text-center mb-6">Create Your First Profile</h2>

        {#if errorMessage}
            <div class="alert alert-error mb-4">
                {errorMessage}
            </div>
        {/if}

        {#if currentStep === 1}
            <!-- Step 1: Profile Name -->
            <div class="form-group">
                <label for="profileName" class="label">Profile Name</label>
                <input
                    type="text"
                    id="profileName"
                    class="input"
                    bind:value={profileName}
                    placeholder="e.g., Work, Personal, Gaming"
                />
            </div>
            <div class="flex justify-end mt-6">
                <button class="btn variant-filled-primary" on:click={nextStep}>Next</button>
            </div>
        {:else if currentStep === 2}
            <!-- Step 2: Database Type -->
            <div class="form-group">
                <label class="label">Database Type</label>
                <ul class="space-y-2">
                    <li>
                        <label class="flex items-center space-x-2">
                            <input
                                type="radio"
                                name="dbType"
                                value={DBType.DBTypeSQLite}
                                bind:group={dbType}
                                class="radio"
                            />
                            <span>Local (SQLite)</span>
                        </label>
                    </li>
                    <li>
                        <label class="flex items-center space-x-2">
                            <input
                                type="radio"
                                name="dbType"
                                value={DBType.DBTypeTurso}
                                bind:group={dbType}
                                class="radio"
                            />
                            <span>Cloud Sync (Turso)</span>
                        </label>
                    </li>
                </ul>
            </div>

            {#if dbType === DBType.DBTypeTurso}
                <!-- Conditional Step 2.1: Turso Details -->
                <div class="form-group mt-4">
                    <label for="tursoUrl" class="label">Turso Database URL</label>
                    <input
                        type="text"
                        id="tursoUrl"
                        class="input"
                        bind:value={tursoUrl}
                        placeholder="e.g., libsql://your-db-name-username.turso.io"
                    />
                    <p class="text-sm text-surface-500 mt-1">
                        Find this in your Turso dashboard.
                    </p>
                </div>
                <div class="form-group mt-4">
                    <label for="tursoToken" class="label">Turso Auth Token</label>
                    <input
                        type="password"
                        id="tursoToken"
                        class="input"
                        bind:value={tursoToken}
                        placeholder="Your Turso auth token"
                    />
                    <p class="text-sm text-surface-500 mt-1">
                        Generate a new token in your Turso dashboard.
                    </p>
                </div>
            {/if}

            <div class="flex justify-between mt-6">
                <button class="btn variant-ghost-primary" on:click={prevStep}>Back</button>
                <button
                    class="btn variant-filled-primary"
                    on:click={createAndSwitchProfile}
                    disabled={isLoading}
                >
                    {#if isLoading}
                        Creating...
                    {:else}
                        Create Profile
                    {/if}
                </button>
            </div>
        {/if}
    </div>
</div>
