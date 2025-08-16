import { writable } from 'svelte/store';
import type { Profile } from '@bindings/github.com/sriram15/progressor-todo-app/internal/profile/models';

type ProfileState = 'loading' | 'no-profiles' | 'needs-selection' | 'profile-active';

interface ProfileStore {
    state: ProfileState;
    activeProfile: Profile | null;
}

const { subscribe, update } = writable<ProfileStore>({
    state: 'loading',
    activeProfile: null,
});

export const profileStore = {
    subscribe,
    setProfileState: (state: ProfileState) => update(store => ({ ...store, state })),
    setActiveProfile: (profile: Profile) => update(store => ({ ...store, activeProfile: profile, state: 'profile-active' }))
};
