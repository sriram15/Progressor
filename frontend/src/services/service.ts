// This file is intended to be an abstraction for the frontend to not know about @wailsjs APIs
// In the future, this will help in converting the frontend code to use fetch() request instead of relying on API
// Right now, we are just going to export the same as wails App
import {
  GetAll,
  GetOpenCards,
  GetCardById,
  AddCard,
  DeleteCard,
  UpdateCard,
  UpdateCardStatus,
  StartCard,
  StopCard,
  GetActiveTimeEntry,
} from "@wailsjs/go/service/cardService";

import { GetStats } from "@wailsjs/go/service/progressService";

export {
  GetAll,
  GetOpenCards,
  GetCardById,
  AddCard,
  DeleteCard,
  UpdateCard,
  UpdateCardStatus,
  StartCard,
  StopCard,
  GetStats,
  GetActiveTimeEntry,
};
