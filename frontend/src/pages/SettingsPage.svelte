<script lang="ts">
    import { GetAllSettings } from "@wailsjs/go/service/settingService";

    let settingsPromise = $state(GetAllSettings());
</script>

<div class="flex flex-col p-4">
    <h3 class="h3">Settings</h3>
    {#await settingsPromise}
        <p>Loading settings...</p>
    {:then settingsList}
        <table class="p-4 table-auto w-full">
            <tbody>
                {#each settingsList as setting}
                    <tr>
                        <td class="border px-4 py-2">{setting.display}</td>
                        <td class="border px-4 py-2">{setting.value}</td>
                    </tr>
                {/each}
            </tbody>
        </table>
    {:catch error}
        <p style="color: red">{error.message}</p>
    {/await}
</div>
