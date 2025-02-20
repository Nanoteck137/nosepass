import type { LayoutServerLoad } from "./$types";

export const load: LayoutServerLoad = async ({ locals }) => {
  const user = locals.user;

  return {
    apiAddress: locals.apiAddress,
    userToken: locals.token,
    user,
  };
};
