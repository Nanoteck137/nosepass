<script lang="ts">
  import VideoPlayer from "$lib/components/VideoPlayer.svelte";
  import { Button, DropdownMenu, Select } from "@nanoteck137/nano-ui";

  const { data } = $props();

  let error = $state<string>();
  let media = $state<(typeof data.collection.media)[number]>();

  let videoUrl = $derived(getVideoUrl());

  // TODO(patrik): Get default audio track from data
  let audioIndex = $state(0);
  // TODO(patrik): Get default subtitle from data
  let subtitleIndex = $state(0);

  function getVideoUrl() {
    if (!media) {
      return null;
    }

    return `${data.apiAddress}/api/stream2/${media.id}/index.m3u8?audio=${audioIndex}&subtitle=${subtitleIndex}`;
  }
</script>

<p>{data.collection.name}</p>

{#if !!media}
  <VideoPlayer
    videoUrl={videoUrl!}
    {error}
    audioTracks={media.audioTracks.map((t) => ({
      index: t.index,
      name: t.language,
    }))}
    subtitles={media.subtitles.map((s) => ({
      index: s.index,
      title: s.title,
    }))}
    onAudioTrackSelected={(index) => {
      audioIndex = index;
    }}
    onSubtitleSelected={(index) => {
      subtitleIndex = index;
    }}
  />
{/if}

<div class="flex flex-col gap-2">
  {#each data.collection.media as m}
    <Button
      onclick={() => {
        media = m;
      }}
    >
      {m.path}
    </Button>
  {/each}
</div>
