<script lang="ts">
  import type { UserSkill, UserSkillProgress } from "@bindings_database";

  const { skill, progress } = $props<{
    skill: UserSkill;
    progress?: UserSkillProgress;
  }>();

  console.log("Progress prop", progress)
  function formatMinutesToHours(minutes: number): string {
    console.log("Progress", minutes)
    if (minutes < 60) {
      return `${minutes} min`;
    } else {
      const hours = Math.floor(minutes / 60);
      const remainingMinutes = minutes % 60;
      return `${hours} hr ${remainingMinutes} min`;
    }
  }
</script>

<div class="card p-4 flex justify-between items-center">
  <div
    title={skill.description && skill.description.Valid
      ? skill.description.String
      : ""}
  >
    <p class="h3">{skill.name}</p>
    {#if skill.description && skill.description.Valid}
      <p class="text-surface-500 text-sm">
        {skill.description.String.length > 200
          ? `${skill.description.String.substring(0, 200)}...`
          : skill.description.String}
      </p>
    {/if}
  </div>
  <div class="text-right">
    <p class="h2 text-bold">
      {formatMinutesToHours(progress?.total_minutes_tracked.Valid ? progress.total_minutes_tracked.Int64 : 0)}
    </p>
    <p class="text-surface-500 text-sm">tracked</p>
  </div>
</div>
