<script lang="ts">
    import { GetDailyTotalMinutes } from "@/services/service";
    import {
        eachDayOfInterval,
        format,
        endOfMonth,
        startOfMonth,
    } from "date-fns";
    import * as d3 from "d3";
    import { onMount } from "svelte";

    let { loading, data, error } = $state({
        loading: true,
        data: null,
        error: "",
    });

    const year = new Date().getFullYear();
    const month = new Date().getMonth();

    const daysInMonth = eachDayOfInterval({
        start: startOfMonth(new Date(year, month)),
        end: endOfMonth(new Date(year, month)),
    });

    const yTicks = [0, 30, 60, 90, 120];
    const padding = { top: 20, right: 15, bottom: 20, left: 25 };

    let width = $state(300);
    let height = 300;

    let xScale = $derived(
        d3
            .scaleLinear()
            .domain([0, data.length])
            .range([padding.left, width - padding.right]),
    );

    let yScale = d3
        .scaleLinear()
        .domain([0, Math.max.apply(null, yTicks)])
        .range([height - padding.bottom, padding.top]);

    let innerWidth = $derived(width - (padding.left + padding.right));
    let barWidth = $derived(innerWidth / data.length);

    onMount(async () => {
        try {
            const apiData = await GetDailyTotalMinutes();

            data = daysInMonth.map((day: any) => {
                const compareDateFmt = format(day, "yyyy-MM-dd");
                const apiDataAt = apiData.find(
                    (item) => item.date === compareDateFmt,
                );

                return {
                    date: format(day, "MM-dd"),
                    value: apiDataAt?.total_minutes?.Float64 ?? 0,
                };
            });
        } catch (ex) {
            console.log(ex);
            error = ex;
        } finally {
            loading = false;
        }
    });
</script>

<div class="flex flex-col p-4 w-full">
    <h3 class="h3">Progress</h3>
    <div class="flex h-full justify-center items-center">
        <div class="w-full h-full">
            <div class="chart" bind:clientWidth={width}>
                <svg {width} {height}>
                    <g class="bars">
                        {#each data as point, i}
                            <rect
                                x={xScale(i) + 2}
                                y={yScale(point.value)}
                                width={barWidth * 0.9}
                                height={yScale(0) - yScale(point.value)}
                            />

                            <!-- Circle showing the start of each Bar -->
                            <circle
                                cx={xScale(i) + 2}
                                cy={yScale(point.value)}
                                fill="black"
                                r="5"
                            />
                        {/each}
                    </g>

                    <g class="axis y-axis">
                        {#each yTicks as tick}
                            <g
                                class="tick tick-{tick}"
                                transform="translate(0, {yScale(tick)})"
                            >
                                <line x2="100%" />
                                <text y="-4">{tick} </text></g
                            >
                        {/each}
                    </g>

                    <!-- Design x axis -->
                    <g class="axis x-axis">
                        {#each data as point, i}
                            <g
                                class="tick"
                                transform="translate({xScale(i)}, {height})"
                            >
                                <text x={barWidth / 2} y="-4" font-size="10">
                                    {point.date}
                                </text></g
                            >
                        {/each}
                    </g>
                </svg>
            </div>
        </div>
    </div>
</div>

<style>
    .x-axis .tick text {
        text-anchor: middle;
        color: white;
    }

    .bars rect {
        fill: rgba(var(--color-secondary-500) / 1);
        stroke: none;
    }

    .tick {
        font-family: Poppins, sans-serif;
        font-size: 0.725em;
        font-weight: 200;
        color: black;
    }

    .tick text {
        fill: black;
        text-anchor: start;
        color: black;
    }

    .tick line {
        stroke: rgba(var(--color-secondary-500) / 1);
        stroke-dasharray: 2;
        opacity: 1;
    }

    .tick.tick-0 line {
        display: inline-block;
        stroke-dasharray: 0;
    }
</style>
