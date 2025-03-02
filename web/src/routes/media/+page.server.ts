import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  const media = await locals.apiClient.getMedia();
  if (!media.success) {
    throw error(media.error.code, { message: media.error.message });
  }

  return {
    media: media.data.media,
  };
};
