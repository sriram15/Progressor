// This file is intended to be an abstraction for the frontend to not know about @wailsjs APIs
// In the future, this will help in converting the frontend code to use fetch() request instead of relying on API
// Right now, we are just going to export the same as wails App
import {
  GetAll,
  GetCardById,
  AddCard,
  DeleteCard,
  UpdateCard,
  UpdateCardStatus,
  StartCard,
  StopCard,
  GetActiveTimeEntry,
  GetStats,
  GetDailyTotalMinutes,
  GetTotalExpForUser,
  GetAllSettings,
  CreateSkill,
  GetSkillByID,
  GetSkillsByUserID,
  UpdateSkill,
  DeleteSkill,
  GetUserSkillProgress,
  GetProjects, 
  GetSkillsForProject,
  AddProjectSkill,
  RemoveProjectSkill 
} from "@bindings/github.com/sriram15/progressor-todo-app/progressorapp";


export {
  GetAll,
  GetCardById,
  AddCard,
  DeleteCard,
  UpdateCard,
  UpdateCardStatus,
  StartCard,
  StopCard,
  GetActiveTimeEntry,
  ///----- progressService ----
  GetStats,
  GetDailyTotalMinutes,
  GetTotalExpForUser,
  /// ---- settingservice ---
  GetAllSettings,
  /// ---- skillservice ---
  CreateSkill,
  GetSkillByID,
  GetSkillsByUserID,
  UpdateSkill,
  DeleteSkill,
  GetUserSkillProgress,
  /// ---- projectservice ---
  GetSkillsForProject,
  AddProjectSkill,
  RemoveProjectSkill,
  GetProjects,
};
