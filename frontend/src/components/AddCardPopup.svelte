<script lang="ts">
    import { EventsOff, EventsOn } from "@wailsjs/runtime";
    import { onDestroy, onMount } from "svelte";

    let { submit } = $props();

    let dialog = $state<HTMLDialogElement>();
    let cardTitle = $state("");

    onMount(() => {
        EventsOn("globalMenu:AddCard", () => {
            showDialog();
        });
    });

    onDestroy(() => {
        EventsOff("globalMenu:AddCard");
    });

    const showDialog = () => {
        dialog.showModal();
    };

    const resetDialogContent = () => {
        cardTitle = "";
    };

    const closeDialog = () => {
        dialog.close();
        resetDialogContent();
    };

    const handleSubmit = (event: any) => {
        event.preventDefault(); // Prevent the default form submission
        submit({ title: cardTitle }); // Emit the submit event with card title
        closeDialog(); // Close the dialog after submission
    };
</script>

<button class="btn btn-sm variant-filled-primary" onclick={showDialog}>
    Add Card (Ctrl + C)
</button>

<dialog bind:this={dialog} class="centered-dialog">
    <form method="dialog" onsubmit={handleSubmit} class="flex flex-col">
        <h2 class="text-xl mb-4">Add Card Title</h2>
        <input
            type="text"
            placeholder="Enter card title"
            bind:value={cardTitle}
            class="border border-gray-300 p-2 rounded mb-4"
            required
        />
        <div class="flex justify-end">
            <button
                type="button"
                class="bg-gray-300 text-black p-2 rounded mr-2"
                onclick={closeDialog}>Cancel</button
            >
            <button type="submit" class="bg-green-500 text-white p-2 rounded"
                >Submit</button
            >
        </div>
    </form>
</dialog>

<style>
    dialog {
        backdrop-filter: blur(5px);
        background-color: rgba(255, 255, 255, 0.9);
    }
</style>
