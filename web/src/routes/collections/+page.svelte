<script lang="ts">
  const { data } = $props();
</script>

<div class="flex flex-col gap-2">
  {#each data.collections as col}
    <a href="/collections/{col.id}">{col.name}</a>
  {/each}
</div>
