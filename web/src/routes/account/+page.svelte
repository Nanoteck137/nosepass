<script lang="ts">
  import { Button } from "@nanoteck137/nano-ui";

  let { data } = $props();
</script>

<p>{data.user?.username}</p>

<Button href="/account/password">Change Password</Button>
<Button href="/account/tokens">Api Tokens</Button>
