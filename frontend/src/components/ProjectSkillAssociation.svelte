<script lang="ts">
    import { onMount } from "svelte";
    import type { UserSkill, Project } from "@bindings_database";
    import {
        GetProjects,
        GetSkillsForProject,
        AddProjectSkill,
        RemoveProjectSkill,
        GetSkillsByUserID,
    } from "@/services/service";

    let allSkills = $state<UserSkill[]>([]);
    let projects = $state<Project[]>([]);
    let selectedProjectId = $state<number | null>(null);
    let associatedSkillIds = $state<number[]>([]);
    let loading = $state(true);
    let error = $state<string | null>(null);

    let availableSkills = $derived(
        allSkills.filter((s) => !associatedSkillIds.includes(s.id)),
    );

    onMount(async () => {
        loading = true;
        error = null;
        try {
            // Fetch projects and skills in parallel
            const [fetchedProjects, fetchedSkills] = await Promise.all([
                GetProjects(),
                GetSkillsByUserID(1), // Assuming user ID 1
            ]);

            projects = fetchedProjects;
            allSkills = fetchedSkills;

            if (projects.length > 0) {
                selectedProjectId = projects[0].id;
            }
        } catch (e: any) {
            error = e.message;
            console.error("Failed to load projects or skills:", e);
        } finally {
            loading = false;
        }
    });

    $effect(() => {
        if (selectedProjectId !== null) {
            loadAssociatedSkills(selectedProjectId);
        }
    });

    async function loadAssociatedSkills(projectId: number) {
        try {
            const linkedSkills = await GetSkillsForProject(projectId);
            associatedSkillIds = linkedSkills.map((s) => s.id);
        } catch (e: any) {
            error = e.message;
            console.error(`Failed to load skills for project ${projectId}:`, e);
        }
    }

    async function handleAssociateSkill(skillId: number) {
        if (selectedProjectId === null) return;
        try {
            await AddProjectSkill(selectedProjectId, skillId);
            associatedSkillIds = [...associatedSkillIds, skillId];
        } catch (e: any) {
            console.error("Failed to link skill:", e);
            alert(`Error: ${e.message}`);
        }
    }

    async function handleUnlinkSkill(skillId: number) {
        if (selectedProjectId === null) return;
        try {
            await RemoveProjectSkill(selectedProjectId, skillId);
            associatedSkillIds = associatedSkillIds.filter(
                (id) => id !== skillId,
            );
        } catch (e: any) {
            console.error("Failed to unlink skill:", e);
            alert(`Error: ${e.message}`);
        }
    }
</script>

<div class="card p-4">
    <h2 class="h2 mb-4">Associate Skills with Projects</h2>

    {#if loading}
        <div class="space-y-4">
            <div class="skeleton h-10 w-full"></div>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mt-4">
                <div class="skeleton h-48 w-full"></div>
                <div class="skeleton h-48 w-full"></div>
            </div>
        </div>
    {:else if error}
        <div class="alert variant-filled-error">{error}</div>
    {:else if projects.length === 0}
        <div class="alert variant-filled-warning">
            No projects found. Please create a project first.
        </div>
    {:else}
        <div class="space-y-4">
            <label class="label">
                <span>Select Project</span>
                <select class="select" bind:value={selectedProjectId}>
                    {#each projects as project (project.id)}
                        <option value={project.id}>{project.name}</option>
                    {/each}
                </select>
            </label>

            {#if selectedProjectId !== null}
                <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mt-4">
                    <!-- Available Skills -->
                    <div>
                        <h3 class="h3 mb-2">Available Skills</h3>
                        <div
                            class="card p-2 space-y-2 max-h-60 overflow-y-auto"
                        >
                            {#if availableSkills.length === 0}
                                <p class="p-2 text-surface-500">
                                    No more skills to add.
                                </p>
                            {/if}
                            {#each availableSkills as skill (skill.id)}
                                <div
                                    class="card variant-soft p-2 flex justify-between items-center"
                                >
                                    <span>{skill.name}</span>
                                    <button
                                        onclick={() =>
                                            handleAssociateSkill(skill.id)}
                                        class="btn btn-sm variant-filled-success"
                                    >
                                        Add
                                    </button>
                                </div>
                            {/each}
                        </div>
                    </div>

                    <!-- Associated Skills -->
                    <div>
                        <h3 class="h3 mb-2">Associated Skills</h3>
                        <div
                            class="card p-2 space-y-2 max-h-60 overflow-y-auto"
                        >
                            {#if associatedSkillIds.length === 0}
                                <p class="p-2 text-surface-500">
                                    No skills associated yet.
                                </p>
                            {/if}
                            {#each allSkills.filter( (s) => associatedSkillIds.includes(s.id), ) as skill (skill.id)}
                                <div
                                    class="card variant-soft p-2 flex justify-between items-center"
                                >
                                    <span>{skill.name}</span>
                                    <button
                                        onclick={() =>
                                            handleUnlinkSkill(skill.id)}
                                        class="btn btn-sm variant-filled-error"
                                    >
                                        Remove
                                    </button>
                                </div>
                            {/each}
                        </div>
                    </div>
                </div>
            {/if}
        </div>
    {/if}
</div>
