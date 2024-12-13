import { writable } from "svelte/store";

export const eventBus = writable({});

// Function to emit events
export type EventType = string;
export type EventPayload<T> = T;

export function emitEvent<T>(eventType: EventType, payload: EventPayload<T>) {
  eventBus.update((events) => ({
    ...events,
    ...{ type: eventType, payload },
  }));
}
