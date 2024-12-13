<script lang="ts">
    import type { database } from "@wailsjs/go/models";

    let activeCard: database.ListOpenOrCTCardsRow = $props();
</script>

<div class="flex flex-row items-center space-between">
    <span>{activeCard.title}</span>
</div>
