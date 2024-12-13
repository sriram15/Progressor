<script lang="ts">
    import { EVENTS } from "@/constants";
    import { emitEvent, eventBus } from "@/stores/store";
    import { onMount } from "svelte";
    import { GetCardById, UpdateCard } from "@/services/service";
    import { database } from "@wailsjs/go/models";
    import Fa from "svelte-fa";
    import { faClose } from "@fortawesome/free-solid-svg-icons";
    import { createForm } from "svelte-forms-lib";
    import Editor from "./Editor.svelte";

    type GetCardData = {
        isLoading: boolean;
        error: string | null;
        data: database.GetCardRow | null;
    };

    let descriptionEditorRef: any = $state();
    let descriptionData: any = $state({});

    let cardData: GetCardData = $state({
        isLoading: true,
        error: null,
        data: null,
    });
    let selectedCardId: number | null = $state(null);
    let projectId: number | null = $state(null);

    const { form, handleChange, handleSubmit } = createForm({
        initialValues: {
            title: "",
        },
        onSubmit: async (values) => {
            const descriptionData = await descriptionEditorRef.save();

            const updateCardParam = {
                title: values.title,
                description: JSON.stringify(descriptionData),
            };
            await UpdateCard(projectId, selectedCardId, updateCardParam);
            emitEvent(EVENTS.CARD_UPDATED, { id: selectedCardId });
        },
    });

    onMount(() => {
        const unsubscribe = eventBus.subscribe(
            (event: { type: string; payload: any }) => {
                if (event.type === EVENTS.CARD_SELECTED) {
                    const { projectId: pId, cardId } = event.payload;
                    selectedCardId = cardId;
                    projectId = pId;

                    getCardDetails(projectId, selectedCardId);
                }

                if (event.type === EVENTS.CARD_UNSELECTED) {
                    selectedCardId = null;
                    projectId = null;
                }
            },
        );

        return () => unsubscribe();
    });

    const getCardDetails = async (projectId: number, cardId: number) => {
        try {
            const data = await GetCardById(projectId, cardId);

            $form.title = data.title;
            descriptionData =
                data.description.Valid &&
                data.description.String &&
                data.description.String.length > 0
                    ? JSON.parse(data.description.String)
                    : {};

            cardData = {
                isLoading: false,
                error: null,
                data: data,
            };
        } catch (err) {
            cardData = {
                isLoading: false,
                error: err,
                data: null,
            };
        }
    };

    const deselectCard = (cardId: number) => {
        if (!cardId) {
            return;
        }

        emitEvent(EVENTS.CARD_UNSELECTED, cardId);
    };
</script>

{#if selectedCardId != null}
    {#if cardData.isLoading}
        <b>Loading...</b>
    {:else if cardData.error}
        <b>Error while loading the data - {cardData.error}</b>
    {:else}
        <div class="flex flex-row items-center mb-4">
            <h5 class="h5 flex-1">Details</h5>
            <button class="btn variant-filled-secondary" onclick={handleSubmit}
                >Save</button
            >
            <button
                class="btn btn-icon"
                onclick={() => deselectCard(selectedCardId)}
                ><Fa icon={faClose} /></button
            >
        </div>

        <form class="flex flex-col gap-4">
            <label class="label" for="title">
                <input
                    class="border-2 p-2 w-full"
                    name="title"
                    type="text"
                    onchange={handleChange}
                    bind:value={$form.title}
                />
            </label>

            <div class="flex-1 border-2 p-2">
                <Editor
                    class="bg-white h-full"
                    key={`description-${selectedCardId}`}
                    name="description"
                    initialValue={descriptionData}
                    bind:this={descriptionEditorRef}
                />
            </div>
        </form>
    {/if}
{/if}
