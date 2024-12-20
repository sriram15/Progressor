<script lang="ts">
    import Fa from "svelte-fa";
    import {
        faRunning,
        faStop,
        faTrash,
    } from "@fortawesome/free-solid-svg-icons";
    import {
        GetOpenCards,
        DeleteCard,
        UpdateCardStatus,
        StartCard,
    } from "@/services/service";
    import { emitEvent, eventBus } from "@/stores/store";
    import { onMount } from "svelte";
    import { EVENTS } from "@/constants";
    import FixedSidebarPageLayout from "./layouts/FixedSidebarPageLayout.svelte";
    import type { database } from "@wailsjs/go/models";
    import EditCard from "@/components/EditCard.svelte";
    import { StopCard } from "@wailsjs/go/service/cardService";
    import ActiveCard from "@/components/ActiveCard.svelte";

    interface Props {
        projectId: number;
    }

    let { projectId }: Props = $props();

    let selectedCardId: number | null = $state(null);
    let activeCard: database.ListOpenOrCTCardsRow | null = $state(null);
    let cards = $state([]);
    let isLoading = $state(false);
    let error = $state("");

    const loadCards = async (projectId: number) => {
        try {
            const data = await GetOpenCards(projectId);
            cards = data;
            activeCard = data.find((row) => row.isactive) ?? null;
        } catch (err) {
            error = err;
        }
    };

    onMount(() => {
        loadCards(projectId);

        // Event subscription to udpate the current state
        const unsubscribe = eventBus.subscribe((events: any) => {
            if (
                events.type === EVENTS.CARD_ADDED ||
                events.type === EVENTS.CARD_REMOVED ||
                events.type === EVENTS.CARD_UPDATED
            ) {
                loadCards(projectId);
            }

            if (events.type === EVENTS.CARD_START) {
                activeCard = events.payload;
            }

            if (events.type === EVENTS.CARD_STOP) {
                activeCard = null;
            }

            if (events.type === EVENTS.CARD_SELECTED) {
                const { cardId } = events.payload;
                selectedCardId = cardId;
            }

            if (events.type === EVENTS.CARD_UNSELECTED) {
                selectedCardId = null;
            }
        });

        return () => {
            unsubscribe();
        };
    });

    const onCardSelect = (card_id: Number) => {
        emitEvent(EVENTS.CARD_SELECTED, {
            projectId: projectId,
            cardId: card_id,
        });
    };

    const onDeleteCard = async (card: database.ListCardsRow) => {
        // TODO: Make this more pretty or use wails MessageDialog to show this message
        const res = confirm("Are you sure you want to delete the card?");

        if (res) {
            try {
                await DeleteCard(projectId, card.card_id);
                emitEvent(EVENTS.CARD_REMOVED, card);
            } catch (err) {
                console.error(err);
            }
        }
    };

    const onStatusChange = async (card: database.ListCardsRow) => {
        try {
            await UpdateCardStatus(
                projectId,
                card.card_id,
                card.status === 0 ? 1 : 0, // TODO: Update this to use the backedn CArdStatus variable
            );
            emitEvent(EVENTS.CARD_UPDATED, card);
            emitEvent(EVENTS.CARD_UNSELECTED, card.card_id);
        } catch (err) {
            console.error(err);
        }
    };

    const onTrackingChange = async (card: database.ListCardsRow) => {
        try {
            if (!card.isactive) {
                await StartCard(projectId, card.id);
                emitEvent(EVENTS.CARD_START, card);
            } else {
                await StopCard(projectId, card.id);
                emitEvent(EVENTS.CARD_STOP, card);
            }

            emitEvent(EVENTS.CARD_UPDATED, card);
            if (selectedCardId != null) {
                emitEvent(EVENTS.CARD_UNSELECTED, card);
            }
        } catch (err) {
            console.error(err);
        }
    };
</script>

<FixedSidebarPageLayout>
    {#if isLoading}
        <p>...waiting</p>
    {:else if error.length > 0}
        <p style="color: red">{error}</p>
    {:else}
        {@const openCards = cards
            ? cards.filter((card) => card.status == 0)
            : []}
        {@const completedToday = cards
            ? cards.filter((card) => card.status == 1)
            : []}

        {#if activeCard != null}
            <!-- <ActiveCard {activeCard} /> -->
        {/if}

        {#if openCards.length > 0}
            <div class="grid grid-cols-1 gap-1 p-4 pl-8 pr-8">
                <h3 class="h3">Cards</h3>
                {#each openCards as card (card.card_id)}
                    {@render displayCard(card, card.status === 1)}
                {/each}
            </div>
        {:else}
            <p class="text-surface-400">
                Start creating new cards to add to todo
            </p>
        {/if}

        {#if completedToday.length > 0}
            <hr />

            <div class="grid grid-cols-1 gap-1 p-4 pl-8 pr-8">
                <h4 class="h4">Completed Recently</h4>
                {#each completedToday as card (card.card_id)}
                    {@render displayCard(card, card.status === 1)}
                {/each}
            </div>
        {/if}
    {/if}

    {#snippet sidebar()}
        <EditCard />
    {/snippet}
</FixedSidebarPageLayout>

{#snippet displayCard(
    card: database.ListOpenOrCTCardsRow,
    isCompleted: boolean = false,
)}
    <div
        class={`rounded p-2 flex items-center group border-2 ${card.id === selectedCardId ? "border-primary-500" : ""} ${card.isactive ? "border-secondary-500" : ""}`}
        onclick={() => onCardSelect(card.card_id)}
        role="button"
        tabindex="0"
        onkeyup={(event) => event.key === "Enter" && onCardSelect(card.card_id)}
    >
        <input
            type="checkbox"
            id="status-checkbox-{card.card_id}"
            class="mr-2"
            checked={isCompleted}
            onchange={() => onStatusChange(card)}
        />
        <p
            class={`flex flex-1 items-center justify-start flex-grow ${isCompleted ? "line-through" : ""}`}
            style="flex-basis: 65%"
        >
            {card.title}
        </p>
        <p style="flex-basis: 10%" class="text-surface-400">
            {card.trackedmins} mins
        </p>
        <div
            class={`opacity-0 group-hover:opacity-100 ${card.id === selectedCardId ? "opacity-100" : ""} flex flex-1 items-center justify-end flex-grow`}
            style="flex-basis: 25%"
        >
            <button class="btn btn-icon" onclick={() => onTrackingChange(card)}
                ><Fa icon={!card.isactive ? faRunning : faStop} /></button
            >
            <button class="btn btn-icon" onclick={() => onDeleteCard(card)}
                ><Fa icon={faTrash} /></button
            >
        </div>
    </div>
{/snippet}
