<script lang="ts">
    import { AppBar } from "@skeletonlabs/skeleton";
    import "./app.css";
    import { Fa } from "svelte-fa";
    import {
        faCheckCircle,
        faUserCog,
        faTools,
        faArrowUpAZ,
        faHome,
        faArrowUpLong,
    } from "@fortawesome/free-solid-svg-icons";

    import { Router, links, Route, Link } from "svelte-routing";
    import SearchBar from "./components/SearchBar.svelte";
    import SettingsPage from "./pages/SettingsPage.svelte";
    import ProjectViewPage from "./pages/ProjectViewPage.svelte";
    import ProgressPage from "./pages/ProgressPage.svelte";
    import OnboardingPage from './pages/OnboardingPage.svelte'; // New import
    import { GetProfiles, SwitchProfile } from '@bindings/github.com/sriram15/progressor-todo-app/progressorapp'; // New import
    import { onMount } from "svelte";

    const { url = "" } = $props();

    let hasProfiles = $state(false); // New state variable

    const getProfiles = async () => {
        try {
            const profiles = await GetProfiles();
            console.log("Fetched profiles:", profiles);
            if (profiles && profiles.length > 0) {
                // Switch to the first profile to establish an active session
                await SwitchProfile(profiles[0].id);
                hasProfiles = true;
            } else {
                hasProfiles = false;
            }
        } catch (error) {
            console.error("Error fetching profiles:", error);
            hasProfiles = false; // Assume no profiles or error, show onboarding
        }
    };

    onMount(() => {
        getProfiles();
    });
</script>

{#if !hasProfiles}
    <OnboardingPage />
{:else}
    <div id="app-shell" class="grid h-screen">
        <Router {url}>
            <header>
                <div
                    class="border-b border-gray-300 flex flex-row items-center h-16 p-2"
                >
                    <h1 class="h1 flex-1">Progressor</h1>
                    <nav class="flex-1 flex flex-row items-center gap-4">
                        <Link to="/" let:active>
                            <span
                                class={`${active ? "navlink-active" : ""}
                                p-2 flex items-center gap-2`}
                                ><Fa icon={faHome} />Home</span
                            >
                        </Link>
                        <Link to="/progress" let:active>
                            <span
                                class={`${active ? "navlink-active" : ""}
                                p-2 flex items-center gap-2`}
                                ><Fa icon={faArrowUpLong} />Progress</span
                            >
                        </Link>
                        <Link to="/settings" let:active>
                            <span
                                class={`${active ? "navlink-active" : ""}
                                p-2 flex items-center gap-2`}
                                ><Fa icon={faUserCog} />Settings</span
                            >
                        </Link>
                    </nav>
                    <SearchBar></SearchBar>
                </div>
            </header>

            <main class="p-4 flex w-full overflow-hidden">
                <Route path="/"><ProjectViewPage projectId={1} /></Route>
                <Route path="/progress"><ProgressPage /></Route>
                <Route path="/settings"><SettingsPage /></Route>
            </main>
        </Router>
    </div>
{/if}

<style>
    #app-shell {
        display: grid;
        grid-template-rows: auto 1fr;
        grid-template-columns: 1fr;
        grid-template-areas:
            "header"
            "main";

        min-block-size: 100vh;
        min-block-size: 100dvh;
    }

    header {
        grid-area: header;
    }

    aside {
        grid-area: aside;
    }

    main {
        grid-area: main;
    }
</style>
