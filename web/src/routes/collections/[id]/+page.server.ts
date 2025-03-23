import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params }) => {
  const collection = await locals.apiClient.getCollectionById(params.id);
  if (!collection.success) {
    throw error(collection.error.code, { message: collection.error.message });
  }

  return {
    collection: collection.data,
  };
};
