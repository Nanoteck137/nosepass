<script lang="ts">
  const { data } = $props();
</script>

<div class="flex flex-col">
  {#each data.media as media}
    <a href="/media/{media.id}">{media.path}</a>
  {/each}
</div>
