<script lang="ts">
    import { onMount } from "svelte";
    import { GetStats } from "../services/service";
    import PercentageDisplay from "./PercentageDisplay.svelte";
    import { getDaysInMonth } from "date-fns";

    let data = $state(null);
    let loading = $state(true);
    let error = $state("");

    onMount(async () => {
        try {
            data = await GetStats();
        } catch (ex) {
            error = ex;
        } finally {
            loading = false;
        }
    });

    const totalDaysInCurrentMonth = getDaysInMonth(new Date());

    function getPercentDiff(current: number, previous: number): string {
        if (previous === 0) {
            if (current === 0) return "0";
            // If previous is 0 and current is not, treat as 100% increase or -100% decrease
            return current > 0 ? "100" : "-100";
        }
        const diff = ((current - previous) / Math.abs(previous)) * 100;
        return diff.toFixed(2);
    }
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
                        <h4 class="h4">{data.weekHrs.value} Hrs</h4>
                    </div>
                    <PercentageDisplay
                        current={data.weekHrs.value}
                        previous={data.weekHrs.prevValue}
                        suffix=" vs Previous week"
                    />
                </div>
                <div class="rounded-lg flex-1 p-4 bg-black text-white">
                    <div class="flex flex-row justify-between">
                        <p>Progress</p>
                        <h4 class="h4">{data.weekProgress.value}/7 Days</h4>
                    </div>
                    <PercentageDisplay
                        current={data.weekProgress.value}
                        previous={data.weekProgress.prevValue}
                        suffix=" vs Previous week"
                    />
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
                        <h4 class="h4">{data.monthHrs.value} Hrs</h4>
                    </div>
                    <PercentageDisplay
                        current={data.monthHrs.value}
                        previous={data.monthHrs.prevValue}
                        suffix=" vs Previous month"
                    />
                </div>
                <div class="rounded-lg flex-1 p-4 bg-black text-white">
                    <div class="flex flex-row justify-between">
                        <p>Progress</p>
                        <h4 class="h4">
                            {data.monthProgress.value}/{totalDaysInCurrentMonth}
                            Days
                        </h4>
                    </div>
                    <PercentageDisplay
                        current={data.monthProgress.value}
                        previous={data.monthProgress.prevValue}
                        suffix=" vs Previous month"
                    />
                </div>
            </div>
        </div>
    </div>
{/if}
