<script lang="ts">
  import { onMount } from "svelte";
  import type { UserSkill } from "@bindings_database";
  import {
    CreateSkill,
    GetSkillsByUserID,
    UpdateSkill,
    DeleteSkill,
  } from "@/services/service";
  import Fa from "svelte-fa";
  import { faPen, faTrash } from "@fortawesome/free-solid-svg-icons";

  let skills = $state<UserSkill[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);

  let skillName = $state("");
  let skillDescription = $state("");
  let editingSkillId: number | null = $state(null);

  onMount(async () => {
    await loadSkills();
  });

  async function loadSkills() {
    loading = true;
    error = null;
    try {
      skills = await GetSkillsByUserID(1); // Assuming user ID 1
    } catch (e: any) {
      error = e.message;
      console.error("Failed to load skills:", e);
    } finally {
      loading = false;
    }
  }

  async function handleSubmit() {
    if (!skillName.trim()) {
      alert("Skill name cannot be empty.");
      return;
    }

    // No loading state change here for faster UI feedback
    try {
      if (editingSkillId) {
        const updatedSkill = await UpdateSkill(
          editingSkillId,
          skillName,
          skillDescription
        );
        skills = skills.map((s) => (s.id === updatedSkill.id ? updatedSkill : s));
      } else {
        const newSkill = await CreateSkill(1, skillName, skillDescription); // Assuming user ID 1
        skills = [...skills, newSkill];
      }
      resetForm();
    } catch (e: any) {
      error = e.message;
      console.error("Failed to create/update skill:", e);
    }
  }

  function resetForm() {
    skillName = "";
    skillDescription = "";
    editingSkillId = null;
  }

  function handleEdit(skill: UserSkill) {
    editingSkillId = skill.id;
    skillName = skill.name;
    skillDescription =
      skill.description && skill.description.Valid ? skill.description.String : "";
  }

  async function handleDelete(id: number) {
    if (confirm("Are you sure you want to delete this skill?")) {
      try {
        await DeleteSkill(id);
        skills = skills.filter((skill) => skill.id !== id);
      } catch (e: any) {
        error = e.message;
        console.error("Failed to delete skill:", e);
      }
    }
  }
</script>

<div class="space-y-8">
  <h2 class="h2 mb-4">Skill Management</h2>

  <form onsubmit={handleSubmit} class="mb-6 space-y-4">
    <label class="label">
      <span>Skill Name</span>
      <input
        class="input"
        type="text"
        bind:value={skillName}
        placeholder="e.g., Go Programming"
      />
    </label>
    <label class="label">
      <span>Description (Optional)</span>
      <textarea
        class="textarea"
        bind:value={skillDescription}
        rows="3"
        placeholder="Describe this skill..."
      ></textarea>
    </label>
    <div class="flex space-x-2">
      <button type="submit" class="btn variant-filled-primary">
        {editingSkillId ? "Update Skill" : "Add Skill"}
      </button>
      {#if editingSkillId}
        <button type="button" onclick={resetForm} class="btn variant-outline-secondary">
          Cancel
        </button>
      {/if}
    </div>
  </form>

  {#if loading}
    <div class="space-y-2">
        <div class="skeleton h-20 w-full"></div>
        <div class="skeleton h-20 w-full"></div>
        <div class="skeleton h-20 w-full"></div>
    </div>
  {:else if error}
    <div class="alert variant-filled-error">{error}</div>
  {:else if skills.length === 0}
    <div class="alert variant-filled-warning">No skills added yet. Add one above!</div>
  {:else}
    <div class="space-y-4">
      <h3 class="h3">Your Skills:</h3>
      {#each skills as skill (skill.id)}
        <div class="card p-3 w-full flex justify-between items-center">
          <div>
            <p class="font-bold">{skill.name}</p>
            {#if skill.description && skill.description.Valid}
              <p class="text-sm text-surface-500">{skill.description.String}</p>
            {/if}
          </div>
          <div class="flex space-x-2">
            <button onclick={() => handleEdit(skill)} class="btn btn-icon variant-ghost-secondary w-10 h-10">
              <Fa icon={faPen} />
            </button>
            <button onclick={() => handleDelete(skill.id)} class="btn btn-icon variant-ghost-error w-10 h-10">
              <Fa icon={faTrash} />
            </button>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>
