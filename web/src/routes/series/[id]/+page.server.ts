import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params }) => {
  const serie = await locals.apiClient.getSerieById(params.id);
  if (!serie.success) {
    throw error(serie.error.code, { message: serie.error.message });
  }

  return {
    serie: serie.data,
  };
};
