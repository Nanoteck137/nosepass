<script lang="ts">
  const { data } = $props();
</script>

<p class="p-4 text-xl">Home Page</p>

{#each data.series as serie}
  <a href="/series/{serie.id}">{serie.name}</a>
{/each}
