import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  const collections = await locals.apiClient.getCollections();
  if (!collections.success) {
    throw error(collections.error.code, {
      message: collections.error.message,
    });
  }

  return {
    collections: collections.data.collections,
  };
};
