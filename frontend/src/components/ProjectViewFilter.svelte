<script lang="ts">
    type ProjectViewFilterItem = {
        status: { label: string; value: number };
    };
    type ProjectViewFilterProps = {
        filterStatus: ProjectViewFilterItem;
        setFilterStatus: (newFilterStatus: ProjectViewFilterItem) => void;
    };

    const { filterStatus, setFilterStatus }: ProjectViewFilterProps = $props();

    const onStatusChange = (type: string, value: string | number) => {
        const currentFilters: ProjectViewFilterItem =
            Object.assign(filterStatus);

        switch (type) {
            case "status":
                const numValue = value as number;
                currentFilters.status.value = numValue;
                break;
        }

        setFilterStatus(currentFilters);
    };
</script>

<div class="pt-4 pb-4 flex flex-row gap-8">
    <div>
        <label for="projectSelect" class="text-sm text-black text-opacity-50"
            >Project</label
        >
        <select id="projectSelect">
            <option value={`inbox`}>Inbox</option>
        </select>
    </div>
    <div>
        <label for="statusSelect" class="text-sm text-black text-opacity-50"
            >Status</label
        >
        <button
            onclick={() => onStatusChange("status", 0)}
            class={`btn btn-sm rounded-full border-2 ${filterStatus.status.value == 0 ? " variant-filled-primary" : ""}`}
            >Open</button
        >
        <button
            onclick={() => onStatusChange("status", 1)}
            class={`btn btn-sm rounded-full border-2 ${filterStatus.status.value == 1 ? " variant-filled-primary" : ""}`}
            >Closed</button
        >
    </div>
</div>
