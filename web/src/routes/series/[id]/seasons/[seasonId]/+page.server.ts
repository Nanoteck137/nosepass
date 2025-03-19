import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals, params }) => {
  const season = await locals.apiClient.getSeasonById(params.seasonId);
  if (!season.success) {
    throw error(season.error.code, { message: season.error.message });
  }

  return {
    season: season.data,
  };
};
