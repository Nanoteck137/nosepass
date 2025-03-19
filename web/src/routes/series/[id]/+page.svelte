<script lang="ts">
  const { data } = $props();
</script>

<p>{data.serie.name}</p>

{#each data.serie.seasons as season}
  <a href="/series/{data.serie.id}/seasons/{season.id}">{season.name}</a>
{/each}
