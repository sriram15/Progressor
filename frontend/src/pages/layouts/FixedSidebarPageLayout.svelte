<script>
    // @ts-nocheck

    let {
        children,
        sidebar,
        bottombar = null,
        showBottomBar = false,
    } = $props();
</script>

<div class="grid grid-cols-3 w-full">
    <div class="col-span-2 overflow-y-auto">
        {#if children}
            {@render children?.()}
        {:else}
            <p>Page content</p>
        {/if}
    </div>
    <div
        class="border-2 rounded-sm h-full overflow-y-auto bg-surface-100-800-token"
    >
        {@render sidebar?.()}
    </div>
    <div
        class={`w-1/2 left-[25%] fixed bottom-0 bg-secondary-400 text-white p-4 rounded-t-lg transition-transform ${showBottomBar ? "translate-y-0" : "translate-y-full"}`}
    >
        {@render bottombar?.()}
    </div>
</div>
