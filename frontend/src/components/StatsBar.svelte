<script lang="ts">
    import { onMount } from "svelte";
    import { GetStats } from "../services/service";

    onMount(async () => {
        const data = await GetStats();
        console.log(data);
    });
</script>
