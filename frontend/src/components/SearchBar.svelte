<script lang="ts">
    import { AddCard } from "@/services/service";

    import AddCardPopup from "./AddCardPopup.svelte";
    import { emitEvent } from "@/stores/store";
    import { EVENTS } from "@/constants";
    import { faSearch } from "@fortawesome/free-solid-svg-icons";
    import Fa from "svelte-fa";
    // import { EventsOff, EventsOn } from "@wailsjs/runtime";
    import { onDestroy, onMount } from "svelte";

    onMount(() => {
        // EventsOn("globalMenu:CommandPrompt", () => {
        //     showAddCardPopup = true;
        // });
    });

    onDestroy(() => {
        // EventsOff("globalMenu:CommandPrompt");
    });

    const onAddCardSubmit = async (data: any) => {
        const { title, projectId = 1, estimatedMin = 0 } = data;
        try {
            await AddCard(projectId, title, estimatedMin);
            emitEvent(EVENTS.CARD_ADDED, { title });
        } catch (err) {
            console.error(err);
        }
    };

    let showAddCardPopup = $state(false);
</script>

<button
    class="btn variant-filled rounded"
    onclick={() => (showAddCardPopup = true)}
>
    <span><Fa icon={faSearch} /></span>
    <span>Add/Search...</span>

    <span class="text-sm pl-4">Ctrl+K</span>
</button>
<AddCardPopup
    submit={onAddCardSubmit}
    show={showAddCardPopup}
    onClose={() => (showAddCardPopup = false)}
/>
