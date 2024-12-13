<script lang="ts">
    import { AddCard } from "@/services/service";

    import AddCardPopup from "./AddCardPopup.svelte";
    import { emitEvent } from "@/stores/store";
    import { EVENTS } from "@/constants";

    const onAddCardSubmit = async (data: any) => {
        const { title, projectId = 1, estimatedMin = 0 } = data;
        try {
            await AddCard(projectId, title, estimatedMin);
            emitEvent(EVENTS.CARD_ADDED, { title });
        } catch (err) {
            console.error(err);
        }
    };
</script>

<AddCardPopup submit={onAddCardSubmit} />
