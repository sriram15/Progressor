import { writable } from "svelte/store";
import { UserSkill, UserSkillProgress } from "@bindings_database";

export const skills = writable<UserSkill[]>([]);
export const userSkillProgress = writable<UserSkillProgress[]>([]);
export const loading = writable(false);
export const error = writable<string | null>(null);
