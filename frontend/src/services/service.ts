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
} from "@bindings/github.com/sriram15/progressor-todo-app/internal/service/cardservice";

import {
  GetStats,
  GetDailyTotalMinutes,
  GetTotalExpForUser,
} from "@bindings/github.com/sriram15/progressor-todo-app/internal/service/progressservice";

import { GetAllSettings } from "@bindings/github.com/sriram15/progressor-todo-app/internal/service/settingservice";

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
};
