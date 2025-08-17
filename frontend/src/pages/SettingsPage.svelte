<script lang="ts">
  import SkillManagement from "@/components/SkillManagement.svelte";
  import ProjectSkillAssociation from "@/components/ProjectSkillAssociation.svelte";
  import ProfileWizard from "@/components/ProfileWizard.svelte";
  import { profileStore } from "@/stores/profileStore";
  import { GetAllSettings } from "@/services/service";
    import type { SettingsItem } from "@bindings/github.com/sriram15/progressor-todo-app/internal/service/models";

  type SettingView = "General" | "Skills" | "Projects" | "Profiles";

  let currentView: SettingView = $state("General");
  let settings = $state<SettingsItem[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);
  let showProfileWizard = $state(false);

  $effect(() => {
    async function loadSettings() {
        if (currentView === 'General') {
            loading = true;
            error = null;
            try {
                const settingsMap = await GetAllSettings();
                settings = settingsMap
            } catch (e: any) {
                error = e.message;
            } finally {
                loading = false;
            }
        }
    }
    loadSettings();
  });

  function setView(view: SettingView) {
    currentView = view;
  }

  function handleProfileCreation() {
    showProfileWizard = false;
    // Potentially refresh profile list if displaying a list
  }

  const activeClass = "bg-primary-500 text-on-primary";
  const inactiveClass = "hover:bg-primary-200";
</script>

<div class="flex h-full text-color-base">
  <!-- Left Sidebar for Navigation -->
  <aside class="min-w-48 max-w-64 flex-shrink-0 p-4 border-r border-surface-300">
    <h2 class="h2 mb-6">Settings</h2>
    <nav class="space-y-2">
      <button
        onclick={() => setView("General")}
        class="w-full text-left px-4 py-2 rounded-md font-semibold {currentView === 'General' ? activeClass : inactiveClass}"
      >
        General
      </button>
      <button
        onclick={() => setView("Skills")}
        class="w-full text-left px-4 py-2 rounded-md font-semibold {currentView === 'Skills' ? activeClass : inactiveClass}"
      >
        Skills
      </button>
      <button
        onclick={() => setView("Projects")}
        class="w-full text-left px-4 py-2 rounded-md font-semibold {currentView === 'Projects' ? activeClass : inactiveClass}"
      >
        Projects
      </button>
      <button
        onclick={() => setView("Profiles")}
        class="w-full text-left px-4 py-2 rounded-md font-semibold {currentView === 'Profiles' ? activeClass : inactiveClass}"
      >
        Profiles
      </button>
    </nav>
  </aside>

  <!-- Main Content Area -->
  <main class="flex-1 p-8 overflow-y-auto">
    {#if currentView === "General"}
      <div>
        <h3 class="h3 mb-4">General Settings</h3>
        {#if loading}
            <div class="space-y-2">
                <div class="skeleton h-12 w-full"></div>
                <div class="skeleton h-12 w-full"></div>
                <div class="skeleton h-12 w-full"></div>
            </div>
        {:else if error}
          <p class="text-error-500">Error: {error}</p>
        {:else}
          <div class="card p-4">
            <table class="table table-compact w-full">
              <thead>
                <tr>
                  <th>Setting</th>
                  <th>Value</th>
                </tr>
              </thead>
              <tbody>
                {#each settings as setting}
                  <tr>
                    <td>{setting.display}</td>
                    <td>{setting.value}</td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        {/if}
      </div>
    {/if}

    {#if currentView === "Skills"}
      <SkillManagement />
    {/if}

    {#if currentView === "Projects"}
      <ProjectSkillAssociation />
    {/if}

    {#if currentView === "Profiles"}
      <div>
        <h3 class="h3 mb-4">Profile Management</h3>
        <div class="card p-4 space-y-4">
            <p><strong>Current Profile:</strong> {$profileStore.activeProfile?.name}</p>
            <button class="btn variant-filled-primary" onclick={() => showProfileWizard = true}>
                Create New Profile
            </button>
        </div>
      </div>
    {/if}
  </main>
</div>

{#if showProfileWizard}
  <div class="fixed inset-0 bg-black bg-opacity-50 flex justify-center items-center z-50">
    <div class="relative">
        <button class="btn btn-icon absolute top-2 right-2" onclick={() => showProfileWizard = false}>X</button>
        <ProfileWizard onComplete={handleProfileCreation} />
    </div>
  </div>
{/if}
