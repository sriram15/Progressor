<script lang="ts">
  import SkillManagement from "@/components/SkillManagement.svelte";
  import ProjectSkillAssociation from "@/components/ProjectSkillAssociation.svelte";
  
  import { GetAllSettings } from "@/services/service";
  import type { SettingsItem } from "@bindings/github.com/sriram15/progressor-todo-app/internal/service/models";

  type SettingView = "General" | "Skills" | "Projects";

  let currentView: SettingView = $state("General");
  let settings = $state<SettingsItem[]>([])
  let loading = $state(true);
  let error = $state<string | null>(null);

  $effect(() => {
    async function loadSettings() {
        if (currentView === 'General') {
            loading = true;
            error = null;
            try {
                const settingsMap = await GetAllSettings();
                settings = settingsMap
                console.log("Settings loaded:", settings);
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

    
  </main>
</div>
