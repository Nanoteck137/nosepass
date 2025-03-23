<script lang="ts">
  import VideoPlayer from "$lib/components/VideoPlayer.svelte";
  import { Button } from "@nanoteck137/nano-ui";

  const { data } = $props();

  let variant = $state<(typeof data.media.variants)[number]>();
</script>

<p>{data.media.path}</p>

{#each data.media.variants as v}
  <Button
    onclick={() => {
      variant = v;
    }}
  >
    {v.name}
  </Button>
{/each}

{#if variant}
  <VideoPlayer
    videoUrl="{data.apiAddress}/api/stream/{variant.id}/index.m3u8"
  />
{/if}
