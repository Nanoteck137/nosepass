<script lang="ts">
  import SeekSlider from "$lib/components/SeekSlider.svelte";
  import Spinner from "$lib/components/Spinner.svelte";
  import { cn, formatTime, shouldDisplayHour } from "$lib/utils";
  import { Button, buttonVariants, DropdownMenu } from "@nanoteck137/nano-ui";
  import Hls from "hls.js";
  import {
    EllipsisVertical,
    Fullscreen,
    Merge,
    Pause,
    Play,
  } from "lucide-svelte";
  import { onDestroy, onMount } from "svelte";
  import debounce from "just-debounce-it";
  import screenfull from "screenfull";
  import { fade } from "svelte/transition";
  import { page } from "$app/stores";

  const { data } = $props();

  type VideoState = "playing" | "paused";
  let videoState = $state<VideoState>("paused");

  let audioIndex = $state(0);
  let subtitleIndex = $state(-1);

  let loading = $state(false);

  let time = $state(0);
  let buffered = $state(0);
  let duration = $state(0);

  let volume = $state(0);

  let insideControls = $state(false);
  let isIdle = $state(false);
  let showControls = $derived(
    videoState == "paused" || !isIdle || insideControls,
  );

  let video: HTMLMediaElement;
  let hls = $state<Hls>();

  onMount(() => {
    volume = video.volume;

    if (Hls.isSupported()) {
      hls = new Hls({});

      // bind them together
      hls.attachMedia(video);
      // MEDIA_ATTACHED event is fired by hls object once MediaSource is ready
      // hls.on(Hls.Events.MEDIA_ATTACHED, function () {
      //   console.log("video and hls.js are now bound together !");
      // });

      hls.on(Hls.Events.ERROR, function (event, data) {
        console.log(data);
      });

      hls.on(Hls.Events.MANIFEST_PARSED, function (event, data) {
        console.log(event);
        console.log(
          "manifest loaded, found " + data.levels.length + " quality level",
        );
        console.log(data);

        // hls.subtitleDisplay = true;
        // hls.subtitleTrack = 0;
      });

      hls.loadSource(
        `http://10.28.28.6:3000/${data.media.id}/index.m3u8?audio=${audioIndex}&subtitle=${subtitleIndex}`,
      );
    }

    return () => {
      hls?.destroy();
    };
  });

  let timeout = $state<NodeJS.Timeout>();

  function onActivity() {
    isIdle = false;
    setIdleTimeout();
  }

  const setIdleTimeout = debounce(
    () => {
      clearTimeout(timeout);
      timeout = setTimeout(() => {
        isIdle = true;
      }, 2000);
    },
    250,
    true,
  );

  onDestroy(() => {
    clearTimeout(timeout);
  });
</script>

<svelte:window
  onmousemove={() => {
    console.log("onmousemove");
    onActivity();
  }}
/>

<div class="absolute bottom-0 left-0 right-0 top-0 overflow-clip">
  <!-- svelte-ignore a11y_media_has_caption -->
  <video
    bind:this={video}
    id="video"
    class="h-full w-full"
    crossorigin="anonymous"
    playsinline
    disablePictureInPicture
    autoplay
    onloadeddata={() => {
      console.log("onloadeddata");
      loading = false;
    }}
    onloadedmetadata={() => {
      console.log("onloadedmetadata");
      duration = video.duration;
      onActivity();
    }}
    ontimeupdate={() => {
      time = video.currentTime;
    }}
    onplay={() => {
      console.log("onplay");
      videoState = "playing";
    }}
    onpause={() => {
      console.log("onpause");
      videoState = "paused";
    }}
    onplaying={() => {
      console.log("onplaying");
      loading = false;
    }}
    onloadstart={() => {
      console.log("onloadstart");
      loading = true;
    }}
    onwaiting={() => {
      console.log("onwaiting");
      loading = true;
    }}
    onseeking={() => {
      loading = true;
    }}
    onseeked={() => {
      loading = false;
    }}
    onprogress={() => {
      const duration = video.duration;
      if (duration > 0) {
        for (let i = 0; i < video.buffered.length; i++) {
          const bufferedStart = video.buffered.start(
            video.buffered.length - 1 - i,
          );

          if (bufferedStart < video.currentTime) {
            console.log(
              "Bufferd",
              video.buffered.end(video.buffered.length - 1 - i),
            );

            buffered = video.buffered.end(video.buffered.length - 1 - i);

            // document.getElementById("buffered-amount").style.width = `${
            //   (video.buffered.end(video.buffered.length - 1 - i) * 100) /
            //   duration
            // }%`;
            break;
          }
        }
      }
    }}
    onvolumechange={() => {
      volume = video.volume;
    }}
  >
  </video>

  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class={`absolute inset-0 z-40 opacity-0`}
    onclick={() => {
      onActivity();

      if (video.paused) {
        video.play();
      } else {
        video.pause();
      }
    }}
  ></div>

  <div
    class="absolute left-1/2 top-1/2 z-50 -translate-x-1/2 -translate-y-1/2"
  >
    {#if loading}
      <Spinner />
    {/if}
  </div>

  {#if showControls}
    <div
      class="absolute bottom-0 left-0 right-0 z-50 h-12 bg-red-200"
      transition:fade
      onpointerenter={() => {
        console.log("Enter");
        insideControls = true;
      }}
      onpointerleave={() => {
        console.log("Leave");
        insideControls = false;
      }}
    >
      <SeekSlider
        value={time / duration}
        buffered={buffered / duration}
        onSeek={(v) => {
          const t = duration * v;
          video.currentTime = t;
          time = t;
        }}
      />

      <div class="flex">
        {#if videoState == "playing"}
          <Button
            class="rounded-full"
            size="icon"
            variant="ghost"
            onclick={() => {
              video.pause();
            }}
          >
            <Pause />
          </Button>
        {:else}
          <Button
            class="rounded-full"
            size="icon"
            variant="ghost"
            onclick={() => {
              video.play();
            }}
          >
            <Play />
          </Button>
        {/if}

        <p>
          {formatTime(time, shouldDisplayHour(duration))} / {formatTime(
            duration,
            shouldDisplayHour(duration),
          )}
        </p>

        <Button
          size="icon"
          variant="ghost"
          onclick={() => {
            screenfull.toggle();
          }}
        >
          <Fullscreen />
        </Button>

        <DropdownMenu.Root>
          <DropdownMenu.Trigger
            class={cn(
              buttonVariants({ variant: "ghost", size: "icon-lg" }),
              "rounded-full",
            )}
          >
            <EllipsisVertical />
          </DropdownMenu.Trigger>

          <DropdownMenu.Content align="end">
            <DropdownMenu.Group>
              {#each data.media.audioTracks as track}
                <DropdownMenu.Item
                  onSelect={() => {
                    audioIndex = track.index;
                    hls?.loadSource(
                      `${data.apiAddress}/${data.media.id}/index.m3u8?audio=${audioIndex}&subtitle=${subtitleIndex}`,
                    );
                    video.currentTime = time;
                  }}
                >
                  {track.language}
                </DropdownMenu.Item>
              {/each}
            </DropdownMenu.Group>
          </DropdownMenu.Content>
        </DropdownMenu.Root>

        <DropdownMenu.Root>
          <DropdownMenu.Trigger
            class={cn(
              buttonVariants({ variant: "ghost", size: "icon-lg" }),
              "rounded-full",
            )}
          >
            <EllipsisVertical />
          </DropdownMenu.Trigger>

          <DropdownMenu.Content align="end">
            <DropdownMenu.Group>
              <DropdownMenu.Item
                onSelect={() => {
                  subtitleIndex = -1;
                  hls?.loadSource(
                    `${data.apiAddress}/${data.media.id}/index.m3u8?audio=${audioIndex}&subtitle=${subtitleIndex}`,
                  );
                  video.currentTime = time;
                }}
              >
                No subtitles
              </DropdownMenu.Item>
              {#each data.media.subtitles as subtitle}
                <DropdownMenu.Item
                  onSelect={() => {
                    subtitleIndex = subtitle.index;
                    hls?.loadSource(
                      `${data.apiAddress}/${data.media.id}/index.m3u8?audio=${audioIndex}&subtitle=${subtitleIndex}`,
                    );
                    video.currentTime = time;
                  }}
                >
                  {subtitle.title} -
                  {subtitle.language}
                </DropdownMenu.Item>
              {/each}
            </DropdownMenu.Group>
          </DropdownMenu.Content>
        </DropdownMenu.Root>

        <div class="w-32">
          <SeekSlider
            value={volume}
            disableProgress
            onSeek={(v) => {
              // const t = duration * v;
              // video.currentTime = t;
              // time = t;
              video.volume = v;
            }}
          />
        </div>
      </div>
    </div>
  {/if}
</div>
