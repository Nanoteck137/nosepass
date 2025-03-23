<script lang="ts">
  import VideoPlayer from "$lib/components/VideoPlayer.svelte";
  import { Button, DropdownMenu, Select } from "@nanoteck137/nano-ui";

  const { data } = $props();

  let error = $state<string>();
  let videoUrl = $state<string>();

  const languages = $derived(getAvailableLanguages());
  let languageSelected = $state(getAvailableLanguages()[0]);

  function getAvailableLanguages() {
    const l = data.collection.media.map((e) =>
      e.variants.map((v) => v.language),
    );

    const all = l.reduce((p, n) => {
      return p.concat(n);
    });

    return Array.from(new Set(all));
  }
</script>

<p>{data.collection.name}</p>

{#if videoUrl !== undefined}
  <VideoPlayer bind:videoUrl {error} />
{/if}

<Select.Root type="single" bind:value={languageSelected}>
  <Select.Trigger class="w-[180px]">
    {languageSelected}
  </Select.Trigger>
  <Select.Content>
    <Select.Group>
      <Select.GroupHeading>Languages</Select.GroupHeading>
      {#each languages as language}
        <Select.Item value={language} label={language} />
      {/each}
    </Select.Group>
  </Select.Content>
</Select.Root>

<div class="flex">
  {#each languages as language}
    <Button
      onclick={() => {
        type Episode = {
          name: string;
          url: string | null;
        };

        const episodes = data.collection.media.map((e) => {
          const variant = e.variants.find((v) => v.language === language);

          return {
            name: e.path,
            url: variant
              ? `${data.apiAddress}/api/stream/${variant.id}/index.m3u8`
              : null,
          } as Episode;
        });

        console.log(episodes);
      }}
    >
      {language}
    </Button>
  {/each}
</div>

<div class="flex flex-col gap-2">
  {#each data.collection.media as media}
    <Button
      onclick={() => {
        // const variant = episode.variants[0];
        const variant = media.variants.find(
          (v) => v.language === languageSelected,
        );

        if (!variant) {
          error = "No video with language available";
          videoUrl = "";
          return;
        }

        videoUrl = `${data.apiAddress}/api/stream/${variant.id}/index.m3u8`;
      }}
      disabled={!media.variants.filter((v) => v.language === languageSelected)}
    >
      {media.path}
    </Button>
  {/each}
</div>
