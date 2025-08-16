<script lang="ts">
    import ProfileWizard from '@/components/ProfileWizard.svelte';
    import { profileStore } from '@/stores/profileStore';

    function handleOnboardingComplete() {
        profileStore.setProfileState('profile-active');
    }
</script>

<div class="container h-full mx-auto flex justify-center items-center">
    <ProfileWizard onComplete={handleOnboardingComplete} />
</div>

