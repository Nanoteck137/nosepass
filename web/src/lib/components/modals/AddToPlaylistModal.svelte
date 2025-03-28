<script lang="ts">
  import { getApiClient, handleApiError } from "$lib";
  import type { ApiClient } from "$lib/api/client";
  import type { Playlist, Track } from "$lib/api/types";
  import { Button, Dialog, Input, ScrollArea } from "@nanoteck137/nano-ui";
  import toast from "svelte-5-french-toast";
  import type { ModalProps } from "svelte-modals";

  export type Props = {
    track: Track;
    playlists: Playlist[];
  };

  const { track, playlists, isOpen, close }: Props & ModalProps = $props();
  const apiClient = getApiClient();
</script>

<Dialog.Root
  controlledOpen
  open={isOpen}
  onOpenChange={(v) => {
    if (!v) {
      close(null);
    }
  }}
>
  <Dialog.Content class="flex flex-col gap-4">
    <Dialog.Header>
      <Dialog.Title>Save track to playlist</Dialog.Title>
    </Dialog.Header>

    <ScrollArea class="max-h-36 overflow-y-clip">
      <div class="flex flex-col">
        {#each playlists as playlist, i}
          <Button
            variant="ghost"
            onclick={async () => {
              const res = await apiClient.addItemToPlaylist(playlist.id, {
                trackId: track.id,
              });

              if (!res.success) {
                if (res.error.type === "PLAYLIST_ALREADY_HAS_TRACK") {
                  toast.error("Track already in playlist");
                } else {
                  handleApiError(res.error);
                }

                return;
              }

              toast.success("Track added to playlist");
              close();
            }}
          >
            {playlist.name}
          </Button>
        {/each}
      </div>
    </ScrollArea>

    <Dialog.Footer>
      <Button
        variant="outline"
        onclick={() => {
          close(null);
        }}
      >
        Close
      </Button>
    </Dialog.Footer>
  </Dialog.Content>
</Dialog.Root>
