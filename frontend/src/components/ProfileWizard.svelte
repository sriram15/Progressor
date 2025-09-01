<script lang="ts">
    import {
        CreateProfile,
        SwitchProfile,
    } from "@bindings/github.com/sriram15/progressor-todo-app/progressorapp";
    import {
        Profile,
        DBType,
    } from "@bindings/github.com/sriram15/progressor-todo-app/internal/profile/models";

    let { onComplete } = $props<{ onComplete: () => void }>();

    let currentStep = $state(1);
    let profileName = $state("");
    let dbType = $state<DBType>(DBType.DBTypeSQLite);
    let tursoUrl = $state("");
    let tursoToken = $state("");
    let encryptionKey = $state("");

    let errorMessage = $state("");
    let isLoading = $state(false);

    function nextStep() {
        errorMessage = ""; // Clear previous errors
        if (currentStep === 1 && profileName.trim() === "") {
            errorMessage = "Profile name cannot be empty.";
            return;
        }
        currentStep++;
    }

    function prevStep() {
        errorMessage = ""; // Clear previous errors
        currentStep--;
    }

    async function createAndSwitchProfile() {
        isLoading = true;
        errorMessage = "";
        try {
            const newProfile = new Profile({
                name: profileName,
                dbType: dbType,
                dbUrl: dbType === DBType.DBTypeTurso ? tursoUrl : undefined,
            });

            const createdProfile = await CreateProfile(
                newProfile,
                tursoToken,
                encryptionKey,
            );

            if (createdProfile && createdProfile.id) {
                await SwitchProfile(createdProfile.id);
                onComplete();
            } else {
                errorMessage =
                    "Failed to create profile: No profile ID returned.";
            }
        } catch (error: any) {
            console.error("Error creating or switching profile:", error);
            errorMessage = `Error: ${error.message || "An unknown error occurred."}`;
        } finally {
            isLoading = false;
        }
    }
</script>

<div class="card p-8 w-full max-w-md shadow-xl">
    <h2 class="h2 text-center mb-6">Create a New Profile</h2>

    {#if errorMessage}
        <div class="alert alert-error mb-4">
            {errorMessage}
        </div>
    {/if}

    {#if currentStep === 1}
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
            <button class="btn variant-filled-primary" on:click={nextStep}
                >Next</button
            >
        </div>
    {:else if currentStep === 2}
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
            <div class="form-group mt-4">
                <label for="encryptionKey" class="label">Encryption Key</label>
                <input
                    type="password"
                    id="encryptionKey"
                    class="input"
                    bind:value={encryptionKey}
                    placeholder="Create a new Encryption key to encrypt your data of replica stored in local"
                />
                <p class="text-sm text-surface-500 mt-1">
                    NOTE: This key is used to encrypt your data for your local
                    replica. People cannot open your local replica db without
                    this key.
                </p>
            </div>
        {/if}

        <div class="flex justify-between mt-6">
            <button class="btn variant-ghost-primary" on:click={prevStep}
                >Back</button
            >
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
