<script lang="ts">
    import {
        GetDailyTotalMinutes,
        GetTotalExpForUser,
    } from "@/services/service";
    import {
        eachDayOfInterval,
        format,
        endOfMonth,
        startOfMonth,
        addMonths, // Add this
        subMonths  // Add this
    } from "date-fns";
    import { onMount } from "svelte";
    import { fade } from "svelte/transition";

    let { loading, error } = $state({ // data removed from here
        loading: true,
        error: "",
    });

    let totalExp = $state(0);
    let apiData = $state([]); // apiData declared here
    let currentDate = $state(new Date()); // currentDate declared here

    // year and month constants removed

    const daysInMonth = $derived(eachDayOfInterval({
        start: startOfMonth(currentDate),
        end: endOfMonth(currentDate),
    }));

    const verticalCount = 7;
    const rectSize = 50;
    const rectPadding = 15;

    let width = 500;
    let height = 400;

    const data = $derived(daysInMonth.map((day: any) => {
        const compareDateFmt = format(day, "yyyy-MM-dd");
        // Ensure apiData is accessed correctly; it might be null initially or if API fails
        const apiDataAt = (apiData || []).find( 
            (item) => item.date === compareDateFmt,
        );

        return {
            date: format(day, "dd"),
            value: apiDataAt?.total_minutes?.Float64 ?? 0,
        };
    }));

    onMount(async () => {
        try {
            // apiData is now declared outside and assigned here
            apiData = await GetDailyTotalMinutes(); 
            totalExp = await GetTotalExpForUser();
        } catch (ex)
            console.log(ex);
            error = ex;
        } finally {
            loading = false;
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
        <button onclick={goToPreviousMonth} class="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300">- Prev</button>
        <h3>{format(currentDate, "MMMM yyyy")}</h3>
        <button onclick={goToNextMonth} class="px-4 py-2 bg-gray-200 rounded hover:bg-gray-300">Next +</button>
        <h3>Total Exp: {totalExp}</h3>
    </div>
    <div
        id="tooltip"
        style="position: absolute; display: none; background: white; border: 1px solid black; padding: 5px; border-radius: 5px;"
    ></div>
    <div class="flex h-full justify-center items-center">
        <div class="w-full h-full">
            <!-- profile chart</script> -->
            <svg {width} {height} style="max-width: 500px; min-width: 300px;">
                <g class="boxes">
                    {#each data as point, i}
                        {@const row: number = Math.floor(i / verticalCount)}
                        {@const col: number = Math.floor(i % verticalCount)}
                        {@const x: number = col * (rectSize + rectPadding) + rectPadding}
                        {@const y: number = row * (rectSize + rectPadding) + rectPadding}

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
                        />
                        <rect
                            x={col * (rectSize + rectPadding) +
                                rectPadding +
                                rectSize / 2 -
                                15}
                            y={row * (rectSize + rectPadding) +
                                rectPadding +
                                rectSize -
                                5}
                            width="30"
                            height="15"
                            fill="white"
                            rx="5"
                            ry="5"
                        />
                        <text
                            x={col * (rectSize + rectPadding) +
                                rectPadding +
                                rectSize / 2}
                            y={row * (rectSize + rectPadding) +
                                rectPadding +
                                rectSize +
                                5}
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
