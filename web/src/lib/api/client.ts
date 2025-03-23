import { z } from "zod";
import * as api from "./types";
import { BaseApiClient, type ExtraOptions } from "./base-client";


export class ApiClient extends BaseApiClient {
  constructor(baseUrl: string) {
    super(baseUrl);
  }
  
  getSystemInfo(options?: ExtraOptions) {
    return this.request("/api/v1/system/info", "GET", api.GetSystemInfo, z.any(), undefined, options)
  }
  
  signup(body: api.SignupBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/signup", "POST", api.Signup, z.any(), body, options)
  }
  
  signin(body: api.SigninBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/signin", "POST", api.Signin, z.any(), body, options)
  }
  
  changePassword(body: api.ChangePasswordBody, options?: ExtraOptions) {
    return this.request("/api/v1/auth/password", "PATCH", z.undefined(), z.any(), body, options)
  }
  
  getMe(options?: ExtraOptions) {
    return this.request("/api/v1/auth/me", "GET", api.GetMe, z.any(), undefined, options)
  }
  
  updateUserSettings(body: api.UpdateUserSettingsBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/settings", "PATCH", z.undefined(), z.any(), body, options)
  }
  
  createApiToken(body: api.CreateApiTokenBody, options?: ExtraOptions) {
    return this.request("/api/v1/user/apitoken", "POST", api.CreateApiToken, z.any(), body, options)
  }
  
  getAllApiTokens(options?: ExtraOptions) {
    return this.request("/api/v1/user/apitoken", "GET", api.GetAllApiTokens, z.any(), undefined, options)
  }
  
  deleteApiToken(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/user/apitoken/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  getEntries(options?: ExtraOptions) {
    return this.request("/api/v1/entries", "GET", api.GetEntries, z.any(), undefined, options)
  }
  
  getEntryById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/entries/${id}`, "GET", api.Entry, z.any(), undefined, options)
  }
  
  createEntry(body: api.CreateEntryBody, options?: ExtraOptions) {
    return this.request("/api/v1/entries", "POST", api.CreateEntry, z.any(), body, options)
  }
  
  editEntry(id: string, body: api.EditEntryBody, options?: ExtraOptions) {
    return this.request(`/api/v1/entries/${id}`, "PATCH", z.undefined(), z.any(), body, options)
  }
  
  deleteEntry(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/entries/${id}`, "DELETE", z.undefined(), z.any(), undefined, options)
  }
  
  getMedia(options?: ExtraOptions) {
    return this.request("/api/v1/media", "GET", api.GetMedia, z.any(), undefined, options)
  }
  
  getMediaById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/media/${id}`, "GET", api.FullMedia, z.any(), undefined, options)
  }
  
  getLibraryStatus(options?: ExtraOptions) {
    return this.request("/api/v1/library", "GET", api.GetLibraryStatus, z.any(), undefined, options)
  }
  
  syncLibrary(options?: ExtraOptions) {
    return this.request("/api/v1/library", "POST", z.undefined(), z.any(), undefined, options)
  }
  
  getCollections(options?: ExtraOptions) {
    return this.request("/api/v1/collections", "GET", api.GetCollections, z.any(), undefined, options)
  }
  
  getCollectionById(id: string, options?: ExtraOptions) {
    return this.request(`/api/v1/collections/${id}`, "GET", api.FullCollection, z.any(), undefined, options)
  }
}
