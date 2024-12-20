<script lang="ts">
    import EditorJS from "@editorjs/editorjs";
    import Header from "@editorjs/header";
    import List from "@editorjs/list";

    const { key, initialValue, name } = $props();

    let holderId = $state("");

    $effect(() => {
        holderId = `${name}_editor_${key}`;

        if (editor.save) {
            // editor.clear();

            if (Object.keys(initialValue).length) {
                editor.render(initialValue);
            } else {
                editor.clear();
            }
        }
    });
    let editor = new EditorJS({
        holder: `${name}_editor_${key}`,

        tools: {
            header: Header,
            list: List,
        },
        placeholder: "Start writing your notes here",
        data: initialValue,
    });

    export async function save() {
        return editor.save();
    }
</script>

<div class="bg-white" id={`${name}_editor_${key}`}></div>
