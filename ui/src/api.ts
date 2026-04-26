import { inject, type App } from "vue";
import { getAdminPass, getUserPass } from "./pass";

export type PlayerState = {
  PlayerId: number;
  FinishedStations: number[];
  CodeMask: string;
  TotalStations: number;
  VotingState?: {
    Players: boolean[];
    MyVote: number;
  };
  IsDead: boolean;
};

export type StationState = {
  CurrentCode: string;
  CooldownUntil: number;
  CooldownDuration: number;
  IsVoting: boolean;
};

export type ResultState = {
  SolvedTasks: number;
  TotalTasks: number;

  AlivePlayers: number;
  TotalPlayers: number;

  GameStart: number;
  GameDuration: number;

  IsVoting: boolean;
};

export type AdminState = {
  SolvedCodes: number;
  IsVoting: boolean;
  Settings: {
    RequiredCodes: number;
    TotalStations: number;
    TotalPlayers: number;
    CodeMask: string;
    CooldownDuration: number;
    GameDuration: number;
  };
  Players: boolean[];
};

export function useApi(): Api {
  return inject(Api.provideKey)!;
}

export class Api {
  baseUrl: string = import.meta.env.VITE_API_BASE_URL;

  static provideKey = Symbol("api-provide-key");

  async getPlayerState(): Promise<PlayerState> {
    const result = await fetch(`${this.baseUrl}/api/player_state`, {
      headers: {
        "X-Pass": getUserPass(),
      },
    });
    const state = await result.json();
    if(state.VotingState) {
        state.VotingState.Players = decodeDeathList(state.VotingState.Players, state.VotingState.TotalPlayers)
    }
    return state;
  }

  async getStationState(
    stationId: number,
    pass: string,
  ): Promise<StationState> {
    const result = await fetch(
      `${this.baseUrl}/api/station_state?station_id=${stationId}`,
      {
        headers: {
          "X-Pass": getAdminPass(),
        },
      },
    );
    const state: StationState = await result.json();
    state.CooldownUntil *= 1000;
    return state;
  }

  async getResultState(): Promise<ResultState> {
    const result = await fetch(`${this.baseUrl}/api/result_state`, {
      headers: {
        "X-Pass": getAdminPass(),
      },
    });
    const state: ResultState = await result.json();
    return state;
  }

  async getAdminState(): Promise<AdminState> {
    const result = await fetch(`${this.baseUrl}/api/admin_state`, {
      headers: {
        "X-Pass": getAdminPass(),
      },
    });
    const state = await result.json();
    state.Players = decodeDeathList(state.Players, state.Settings.TotalPlayers)
    return state;
  }
  async togglePlayer(playerId: number): Promise<void> {
    await fetch(`${this.baseUrl}/api/toggle_player?player_id=${playerId}`, {
      method: "POST",
      headers: {
        "X-Pass": getAdminPass(),
      },
    });
  }

  async updateSetting(setting: string, value: any) {
    await fetch(`${this.baseUrl}/api/settings?${setting}=${value}`, {
      method: "POST",
      headers: {
        "X-Pass": getAdminPass(),
      },
    });
  }

  async submitCode(code: string): Promise<boolean> {
    const result = await fetch(`${this.baseUrl}/api/submit&code=${code}`, {
      method: "POST",
      headers: {
        "X-Pass": getUserPass(),
      },
    });
    if (result.ok) {
      return true;
    }
    throw Error(await result.text());
  }

  install(app: App) {
    app.provide(Api.provideKey, this);
  }
}


function decodeDeathList(rawPlayers: string, totalPlayers: number) {
    const result: boolean[] = []
    var bytes = Uint8Array.fromBase64(rawPlayers)

    for (let i = 0; i <= totalPlayers; i++) {
        result.push((bytes[Math.floor(i / 8)] & 1 << i % 8) > 0)
    }
    return result
}