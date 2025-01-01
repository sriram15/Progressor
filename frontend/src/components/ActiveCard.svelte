<script lang="ts">
    import type { database } from "@wailsjs/go/models";
    import { GetActiveTimeEntry } from "@/services/service";
    import { onMount } from "svelte";
    import { faClock } from "@fortawesome/free-solid-svg-icons";
    import Fa from "svelte-fa";

    type ActiveCardProps = {
        projectId: number;
        activeCard: database.ListOpenOrCTCardsRow;
    };

    let { projectId, activeCard }: ActiveCardProps = $props();

    let { loading, error, data } = $state({
        loading: false,
        error: "",
        data: null,
    });
    let elapsedTimeInSecondSinceStart = $state(0);

    const FetchActiveTimeEntryData = async () => {
        loading = true;

        try {
            const activeEntryData = await GetActiveTimeEntry(
                projectId,
                activeCard.id,
            );
            data = activeEntryData;
            if (data && elapsedTimeInSecondSinceStart == 0) {
                // Initial first time setting the value
                const timeDiffinMilli =
                    Date.now() - new Date(activeEntryData.starttime).getTime();

                const elapsedSeconds = Math.ceil(timeDiffinMilli / 1000);
                elapsedTimeInSecondSinceStart = elapsedSeconds;
            }
        } catch (ex) {
            error = ex;
        } finally {
            loading = false;
        }
    };
    onMount(() => {
        FetchActiveTimeEntryData();

        const interval = setInterval(() => {
            if (elapsedTimeInSecondSinceStart != 0) {
                elapsedTimeInSecondSinceStart += 1;
            }
        }, 1000);

        return () => clearInterval(interval);
    });
</script>

<div
    class="flex flex-row items-center bg-primary-500 p-4 rounded-lg space-between mb-4"
>
    <span class="flex-1">{activeCard.title}</span>
    {#if elapsedTimeInSecondSinceStart > 0}
        <div class="flex items-center gap-4">
            <span
                >{Math.floor(elapsedTimeInSecondSinceStart / 60)} Mins {elapsedTimeInSecondSinceStart %
                    60} Seconds
            </span>
            <Fa icon={faClock} />
        </div>
    {/if}
</div>
