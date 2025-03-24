<script lang="ts">
  import SeekSlider from "$lib/components/SeekSlider.svelte";
  import Spinner from "$lib/components/Spinner.svelte";
  import { cn, formatTime, shouldDisplayHour } from "$lib/utils";
  import { Button, buttonVariants, DropdownMenu } from "@nanoteck137/nano-ui";
  import Hls from "hls.js";
  import { EllipsisVertical, Fullscreen, Pause, Play } from "lucide-svelte";
  import { onDestroy, onMount } from "svelte";
  import debounce from "just-debounce-it";
  import screenfull from "screenfull";
  import { fade } from "svelte/transition";

  const MOVE_AMOUNT = 10;

  type AudioTrack = {
    index: number;
    name: string;
  };

  type Subtitle = {
    index: number;
    title: string;
  };

  type Props = {
    videoUrl: string;
    error?: string;

    audioTracks: AudioTrack[];
    subtitles: Subtitle[];

    onAudioTrackSelected?: (index: number) => void;
    onSubtitleSelected?: (index: number) => void;
  };

  const {
    videoUrl = $bindable(),
    error,
    audioTracks,
    subtitles,
    onAudioTrackSelected,
    onSubtitleSelected,
  }: Props = $props();

  type VideoState = "playing" | "paused";
  let videoState = $state<VideoState>("paused");

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

  let videoPlayerContainer: HTMLDivElement;

  let video = $state<HTMLMediaElement>();
  let hls = $state<Hls>();

  onMount(() => {
    if (!video) return;

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

  $effect(() => {
    hls?.loadSource(videoUrl);
  });
</script>

<svelte:window
  onmousemove={() => {
    console.log("onmousemove");
    onActivity();
  }}
/>

<svelte:document
  onkeydown={(e) => {
    const source = e.target as HTMLElement;
    const exclude = ["input", "textarea"];

    if (exclude.indexOf(source.tagName.toLowerCase()) === -1) {
      console.log(e.key);
      if (e.key === "f") {
        screenfull.toggle(videoPlayerContainer);
        onActivity();
      }

      if (video) {
        if (e.key == " ") {
          if (video.paused) {
            video.play();
          } else {
            video.pause();
          }

          onActivity();
        }

        if (e.key === "ArrowLeft") {
          video.currentTime -= MOVE_AMOUNT;
          time -= MOVE_AMOUNT;
          onActivity();
        }

        if (e.key === "ArrowRight") {
          video.currentTime += MOVE_AMOUNT;
          time += MOVE_AMOUNT;
          onActivity();
        }

        if (e.key === "ArrowUp") {
          video.volume += 0.1;
          onActivity();
        }

        if (e.key === "ArrowDown") {
          video.volume -= 0.1;
          onActivity();
        }
      }
    }
  }}
/>

<div bind:this={videoPlayerContainer} class="relative w-full">
  <!-- svelte-ignore a11y_media_has_caption -->
  <video
    bind:this={video}
    id="video"
    class="h-full w-full"
    crossorigin="anonymous"
    playsinline
    autoplay
    onloadeddata={() => {
      console.log("onloadeddata");
      loading = false;
    }}
    onloadedmetadata={() => {
      if (!video) return;

      console.log("onloadedmetadata");
      duration = video.duration;
      onActivity();
    }}
    ontimeupdate={() => {
      if (!video) return;

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
      if (!video) return;

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
      if (!video) return;

      volume = video.volume;
    }}
  >
  </video>

  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class={`absolute inset-0 z-40 opacity-0 ${showControls ? "" : "cursor-none"}`}
    ondblclick={() => {
      screenfull.toggle(videoPlayerContainer);
    }}
    onclick={() => {
      if (!video) return;

      onActivity();

      if (video.paused) {
        video.play();
      } else {
        video.pause();
      }
    }}
  ></div>

  {#if error}
    <div class="absolute inset-0 bg-black/80"></div>
  {/if}

  <div
    class="absolute left-1/2 top-1/2 z-50 -translate-x-1/2 -translate-y-1/2"
  >
    {#if loading && !error}
      <Spinner />
    {/if}

    {#if error}
      <p>{error}</p>
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
          if (!video) return;

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
              video?.pause();
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
              video?.play();
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
            screenfull.toggle(videoPlayerContainer);
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
              {#each audioTracks as track}
                <DropdownMenu.Item
                  onSelect={() => {
                    onAudioTrackSelected?.(track.index);
                    const t = time;
                    setTimeout(() => {
                      if (video) {
                        video.currentTime = t;
                      }
                    }, 0);
                  }}
                >
                  {track.name}
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
                  onSubtitleSelected?.(-1);
                  const t = time;
                  setTimeout(() => {
                    if (video) {
                      video.currentTime = t;
                    }
                  }, 0);
                }}
              >
                No subtitles
              </DropdownMenu.Item>
              {#each subtitles as subtitle}
                <DropdownMenu.Item
                  onSelect={() => {
                    onSubtitleSelected?.(subtitle.index);
                    const t = time;
                    setTimeout(() => {
                      if (video) {
                        video.currentTime = t;
                      }
                    }, 0);
                  }}
                >
                  {subtitle.title}
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
              if (!video) return;

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
