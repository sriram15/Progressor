<script lang="ts">
    import { onMount } from "svelte";
    import { GetStats } from "../services/service";

    let data = $state(null);
    let loading = $state(true);
    let error = $state("");

    onMount(async () => {
        try {
            data = await GetStats();
            console.log(data);
        } catch (ex) {
            error = ex;
            console.log(ex);
        } finally {
            console.log("finally");
            loading = false;
        }
    });
</script>

{#if loading}
    <p>Loading...</p>
{:else if error.length > 0}
    <p>Error: {error}</p>
{:else}
    <div class="flex flex-col p-4">
        <div class="flex flex-row justify-between">
            <h2 class="h2">Progress</h2>
            <p>Level 1</p>
        </div>

        <div class="pt-8 pb-8">
            <h3 class="h3">This Week</h3>
            <p>Amount of hours worked on all project this week</p>

            <div class="flex justify-between gap-8 mt-4 mb-4">
                <div class="rounded-lg flex-1 p-4 bg-black text-white">
                    <div class="flex flex-row justify-between">
                        <p>Hours</p>
                        <h4 class="h4">{data.weekHrs} Hrs</h4>
                    </div>
                    <p class="text-sm text-white text-opacity-50">
                        2.5% vs Previous week
                    </p>
                </div>
                <div class="rounded-lg flex-1 p-4 bg-black text-white">
                    <div class="flex flex-row justify-between">
                        <p>Progress</p>
                        <h4 class="h4">2 Days</h4>
                    </div>
                    <p class="text-sm text-white text-opacity-50">
                        2.5% vs Previous week
                    </p>
                </div>
            </div>
        </div>
        <div class="pt-8 pb-8">
            <h3 class="h3">This Month</h3>
            <p>Amount of hours worked on all project this month</p>

            <div class="flex justify-between gap-8 mt-4 mb-4">
                <div class="rounded-lg flex-1 p-4 bg-black text-white">
                    <div class="flex flex-row justify-between">
                        <p>Hours</p>
                        <h4 class="h4">{data.monthHrs} Hrs</h4>
                    </div>
                    <p class="text-sm text-white text-opacity-50">
                        2.5% vs Previous week
                    </p>
                </div>
                <div class="rounded-lg flex-1 p-4 bg-black text-white">
                    <div class="flex flex-row justify-between">
                        <p>Progress</p>
                        <h4 class="h4">10 Days</h4>
                    </div>
                    <p class="text-sm text-white text-opacity-50">
                        2.5% vs Previous week
                    </p>
                </div>
            </div>
        </div>
        <div class="pt-8 pb-8">
            <h3 class="h3">This Year</h3>
            <p>Amount of hours worked on all project this year</p>

            <div class="flex justify-between gap-8 mt-4 mb-4">
                <div class="rounded-lg flex-1 p-4 bg-black text-white">
                    <div class="flex flex-row justify-between">
                        <p>Hours</p>
                        <h4 class="h4">{data.yearHrs} Hrs</h4>
                    </div>
                    <p class="text-sm text-white text-opacity-50">
                        2.5% vs Previous week
                    </p>
                </div>
                <div class="rounded-lg flex-1 p-4 bg-black text-white">
                    <div class="flex flex-row justify-between">
                        <p>Progress</p>
                        <h4 class="h4">42 Days</h4>
                    </div>
                    <p class="text-sm text-white text-opacity-50">
                        2.5% vs Previous week
                    </p>
                </div>
            </div>
        </div>
    </div>
{/if}
