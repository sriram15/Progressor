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
            estimatedMins: 0,
        },
        onSubmit: async (values) => {
            const descriptionData = await descriptionEditorRef.save();

            // TODO: Make this easier. Current parsing to string and back to number to make sure value from input is converted to number
            const estimatedMinsStr: string = values.estimatedMins.toString();
            const updateCardParam = {
                title: values.title,
                estimatedMins: parseInt(estimatedMinsStr),
                description: JSON.stringify(descriptionData),
            };
            console.log(updateCardParam);
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
            $form.estimatedMins = data.estimatedmins || 0;
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
        <div
            class="sticky p-4 top-0 bg-surface-100-800-token flex flex-row items-center"
        >
            <h3 class="h3 flex-1">Details</h3>
            <button class="btn variant-ringed-secondary" onclick={handleSubmit}
                >Save</button
            >
            <button
                class="btn btn-icon"
                onclick={() => deselectCard(selectedCardId)}
                ><Fa icon={faClose} /></button
            >
        </div>

        <form class="p-4 flex flex-col gap-4">
            <label class="label" for="title">
                Title
                <input
                    class="border-2 p-2 w-full"
                    name="title"
                    type="text"
                    onchange={handleChange}
                    bind:value={$form.title}
                />
            </label>

            <label class="label" for="description">
                Notes
                <div class="p-2 bg-white">
                    <Editor
                        key={`description-${selectedCardId}`}
                        name="description"
                        initialValue={descriptionData}
                        bind:this={descriptionEditorRef}
                    />
                </div>
            </label>

            <label class="label" for="estimatedMins">
                Estimate (in Mins)
                <input
                    class="border-2 p-2 w-full"
                    name="estimatedMins"
                    type="number"
                    onchange={handleChange}
                    bind:value={$form.estimatedMins}
                />
            </label>
        </form>
    {/if}
{/if}
