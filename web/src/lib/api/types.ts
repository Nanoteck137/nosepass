// DO NOT EDIT THIS: This file was generated by the Pyrin Typescript Generator
import { z } from "zod";

export const GetSystemInfo = z.object({
  version: z.string(),
});
export type GetSystemInfo = z.infer<typeof GetSystemInfo>;

export const Signup = z.object({
  id: z.string(),
  username: z.string(),
});
export type Signup = z.infer<typeof Signup>;

export const SignupBody = z.object({
  username: z.string(),
  password: z.string(),
  passwordConfirm: z.string(),
});
export type SignupBody = z.infer<typeof SignupBody>;

export const Signin = z.object({
  token: z.string(),
});
export type Signin = z.infer<typeof Signin>;

export const SigninBody = z.object({
  username: z.string(),
  password: z.string(),
});
export type SigninBody = z.infer<typeof SigninBody>;

export const ChangePasswordBody = z.object({
  currentPassword: z.string(),
  newPassword: z.string(),
  newPasswordConfirm: z.string(),
});
export type ChangePasswordBody = z.infer<typeof ChangePasswordBody>;

export const GetMe = z.object({
  id: z.string(),
  username: z.string(),
  role: z.string(),
  displayName: z.string(),
});
export type GetMe = z.infer<typeof GetMe>;

export const UpdateUserSettingsBody = z.object({
  displayName: z.string().nullable().optional(),
});
export type UpdateUserSettingsBody = z.infer<typeof UpdateUserSettingsBody>;

export const CreateApiToken = z.object({
  token: z.string(),
});
export type CreateApiToken = z.infer<typeof CreateApiToken>;

export const CreateApiTokenBody = z.object({
  name: z.string(),
});
export type CreateApiTokenBody = z.infer<typeof CreateApiTokenBody>;

export const ApiToken = z.object({
  id: z.string(),
  name: z.string(),
});
export type ApiToken = z.infer<typeof ApiToken>;

export const GetAllApiTokens = z.object({
  tokens: z.array(ApiToken),
});
export type GetAllApiTokens = z.infer<typeof GetAllApiTokens>;

export const EntryInfo = z.object({
  id: z.string(),
  name: z.string(),
});
export type EntryInfo = z.infer<typeof EntryInfo>;

export const GetEntries = z.object({
  entries: z.array(EntryInfo),
});
export type GetEntries = z.infer<typeof GetEntries>;

export const Entry = z.object({
  id: z.string(),
  name: z.string(),
});
export type Entry = z.infer<typeof Entry>;

export const CreateEntry = z.object({
  id: z.string(),
});
export type CreateEntry = z.infer<typeof CreateEntry>;

export const CreateEntryBody = z.object({
  name: z.string(),
});
export type CreateEntryBody = z.infer<typeof CreateEntryBody>;

export const EditEntryBody = z.object({
  name: z.string().nullable(),
});
export type EditEntryBody = z.infer<typeof EditEntryBody>;

export const Media = z.object({
  id: z.string(),
  path: z.string(),
});
export type Media = z.infer<typeof Media>;

export const GetMedia = z.object({
  media: z.array(Media),
});
export type GetMedia = z.infer<typeof GetMedia>;

export const MediaAudioTrack = z.object({
  index: z.number(),
  language: z.string(),
});
export type MediaAudioTrack = z.infer<typeof MediaAudioTrack>;

export const MediaSubtitle = z.object({
  index: z.number(),
  type: z.string(),
  title: z.string(),
  language: z.string(),
  isDefault: z.boolean(),
});
export type MediaSubtitle = z.infer<typeof MediaSubtitle>;

export const MediaVariant = z.object({
  audio_track: z.number(),
  subtitle: z.number().nullable(),
});
export type MediaVariant = z.infer<typeof MediaVariant>;

export const FullMedia = z.object({
  id: z.string(),
  path: z.string(),
  audioTracks: z.array(MediaAudioTrack),
  subtitles: z.array(MediaSubtitle),
  subVariant: MediaVariant.nullable(),
  dubVariant: MediaVariant.nullable(),
});
export type FullMedia = z.infer<typeof FullMedia>;

export const GetLibraryStatus = z.object({
  syncing: z.boolean(),
});
export type GetLibraryStatus = z.infer<typeof GetLibraryStatus>;

export const Collection = z.object({
  id: z.string(),
  name: z.string(),
});
export type Collection = z.infer<typeof Collection>;

export const GetCollections = z.object({
  collections: z.array(Collection),
});
export type GetCollections = z.infer<typeof GetCollections>;

export const FullCollection = z.object({
  id: z.string(),
  name: z.string(),
  media: z.array(FullMedia),
});
export type FullCollection = z.infer<typeof FullCollection>;

