<script lang="ts">
    import {
        GetDailyTotalMinutes,
        GetTotalExpForUser,
        GetSkillsByUserID,
        GetUserSkillProgress,
    } from "@/services/service";
    import {
        faArrowLeft,
        faArrowRight,
    } from "@fortawesome/free-solid-svg-icons";
    import {
        eachDayOfInterval,
        format,
        endOfMonth,
        startOfMonth,
        addMonths, // Add this
        subMonths, // Add this
    } from "date-fns";
    import { onMount } from "svelte";
    import Fa from "svelte-fa";
    import { fade } from "svelte/transition";
    import {
        skills,
        userSkillProgress,
        loading,
        error,
    } from "@/stores/skillStore";
    import SkillProgressDisplay from "@/components/SkillProgressDisplay.svelte";

    let totalExp = $state(0);
    let apiData = $state([]); // apiData declared here
    const today = new Date();
    let currentDate = $state(today); // currentDate declared here

    const daysInMonth = $derived(
        eachDayOfInterval({
            start: startOfMonth(currentDate),
            end: endOfMonth(currentDate),
        }),
    );

    const verticalCount = 7;
    const rectSize = 50;
    const rectPadding = 15;

    let width = 500;
    let height = 400;

    const data = $derived(
        daysInMonth.map((day: any) => {
            const compareDateFmt = format(day, "yyyy-MM-dd'T'00:00:00'Z'");
            console.log(compareDateFmt);
            // Ensure apiData is accessed correctly; it might be null initially or if API fails
            const apiDataAt = (apiData || []).find(
                (item) => item.date === compareDateFmt,
            );

            return {
                date: format(day, "dd"),
                value: apiDataAt?.total_minutes?.Float64 ?? 0,
            };
        }),
    );

    onMount(async () => {
        loading.set(true);
        error.set(null);
        try {
            apiData = await GetDailyTotalMinutes();
            totalExp = await GetTotalExpForUser();

            const fetchedSkills = await GetSkillsByUserID(1); // Assuming user ID 1
            skills.set(fetchedSkills);

            const fetchedProgress: any[] = [];
            for (const skill of fetchedSkills) {
                const progress = await GetUserSkillProgress(1, skill.id); // Assuming user ID 1
                console.log(progress);
                if (progress) {
                    fetchedProgress.push(progress);
                }
            }
            userSkillProgress.set(fetchedProgress);
        } catch (ex: any) {
            console.error(ex);
            error.set(ex.message);
        } finally {
            loading.set(false);
        }
    });

    function goToPreviousMonth() {
        currentDate = subMonths(currentDate, 1);
    }

    function goToNextMonth() {
        currentDate = addMonths(currentDate, 1);
    }

    function showTooltip(value, x, y) {
        const tooltip = document.getElementById("tooltip");
        tooltip.style.display = "block";
        tooltip.style.left = `${x}px`;
        tooltip.style.top = `${y + 2 * rectSize}px`;
        tooltip.innerHTML = `Minutes: ${value}`;
    }

    function hideTooltip() {
        const tooltip = document.getElementById("tooltip");
        tooltip.style.display = "none";
    }
</script>

<div class="flex flex-col p-4 w-full">
    <div class="flex flex-row align-center justify-between">
        <h3>Total Exp: {totalExp}</h3>
    </div>
    <div class="flex flex-row gap-8 mt-4">
        <!-- Left Column: Graph -->
        <div class="flex-1">
            <div
                id="tooltip"
                style="position: absolute; display: none; background: white; border: 1px solid black; padding: 5px; border-radius: 5px;"
            ></div>
            <div class="flex h-full items-center">
                <div>
                    <div class="flex flex-row align-center justify-around">
                        <button class="btn btn-icon" onclick={goToPreviousMonth}>
                            <Fa icon={faArrowLeft} />
                        </button>
                        <h2>{format(currentDate, "MMMM yyyy")}</h2>
                        <button class="btn btn-icon" onclick={goToNextMonth}>
                            <Fa icon={faArrowRight} />
                        </button>
                    </div>
                    <svg {width} {height} style="max-width: 500px; min-width: 300px;">
                        <g class="boxes">
                            {#each data as point, i}
                                {@const row: number = Math.floor(i / verticalCount)}
                                {@const col: number = Math.floor(i % verticalCount)}
                                {@const x: number =
                                    col * (rectSize + rectPadding) + rectPadding}
                                {@const y: number =
                                    row * (rectSize + rectPadding) + rectPadding}

                                <rect
                                    {x}
                                    y={y + 5}
                                    width={rectSize}
                                    height={rectSize}
                                    fill={point.value > 20 ? "green" : "grey"}
                                    rx="10"
                                    ry="10"
                                    in:fade={{ delay: i * 20, duration: 50 }}
                                    onmouseover={() => showTooltip(point.value, x, y)}
                                    onfocus={() => showTooltip(point.value, x, y)}
                                    onmouseout={hideTooltip}
                                    onblur={hideTooltip}
                                    role="img"
                                    aria-label={`Day ${point.date}, Minutes: ${point.value}`}
                                    class:progress-box-today={format(
                                        daysInMonth[i],
                                        "yyyy-MM-dd",
                                    ) === format(today, "yyyy-MM-dd")}
                                />
                                <rect
                                    x={
                                        col * (rectSize + rectPadding) +
                                        rectPadding +
                                        rectSize / 2 -
                                        15
                                    }
                                    y={
                                        row * (rectSize + rectPadding) +
                                        rectPadding +
                                        rectSize -
                                        5
                                    }
                                    width="30"
                                    height="15"
                                    fill="white"
                                    rx="5"
                                    ry="5"
                                />
                                <text
                                    x={
                                        col * (rectSize + rectPadding) +
                                        rectPadding +
                                        rectSize / 2
                                    }
                                    y={
                                        row * (rectSize + rectPadding) +
                                        rectPadding +
                                        rectSize +
                                        5
                                    }
                                    text-anchor="middle"
                                    alignment-baseline="middle"
                                    font-size="10"
                                    fill="black"
                                >
                                    {point.date}
                                </text>
                            {/each}
                        </g>
                    </svg>
                </div>
            </div>
        </div>

        <!-- Right Column: Skills List -->
        <div class="flex-1">
            <h3 class="text-xl font-semibold mb-4">Skill Progress</h3>
            {#if $loading}
                <p class="text-blue-400">Loading skills progress...</p>
            {:else if $error}
                <p class="text-red-400">Error: {$error}</p>
            {:else if $skills.length === 0}
                <p class="text-gray-400">
                    No skills defined yet. Go to Settings to add some!
                </p>
            {:else}
                <div class="flex flex-col gap-4">
                    {#each $skills as skill (skill.id)}
                        <SkillProgressDisplay
                            {skill}
                            progress={$userSkillProgress.find(
                                (p) => p.skill_id === skill.id,
                            )}
                        />
                    {/each}
                </div>
            {/if}
        </div>
    </div>
</div>
