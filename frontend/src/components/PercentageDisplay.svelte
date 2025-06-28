<script lang="ts">
    export let current: number;
    export let previous: number;
    export let suffix: string = "";

    function getPercentDiff(current: number, previous: number): string {
        if (previous === 0) {
            if (current === 0) return "0";
            return current > 0 ? "100" : "-100";
        }
        const diff = ((current - previous) / Math.abs(previous)) * 100;
        return diff.toFixed(2);
    }

    const percentage = getPercentDiff(current, previous);
    const isPositive = parseFloat(percentage) >= 0;
</script>

<span class="text-sm text-white text-opacity-50">
    <span class:text-positive={isPositive} class:text-negative={!isPositive}>
        {percentage}%
    </span>
    {suffix}
</span>