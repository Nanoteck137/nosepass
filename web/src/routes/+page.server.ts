import { error } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types";

export const load: PageServerLoad = async ({ locals }) => {
  const series = await locals.apiClient.getSeries();
  if (!series.success) {
    throw error(series.error.code, { message: series.error.message });
  }

  return {
    series: series.data.series,
  };
};
