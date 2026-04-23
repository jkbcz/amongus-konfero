import { inject, type App } from "vue"

export type PlayerState = {
    FinishedStations: number[]
    CodeMask: string
    TotalStations: number
    IsVoting: boolean
}

export type StationState = {
    CurrentCode: string
    CooldownUntil: number
    CooldownDuration: number
    IsVoting: boolean
}

export type ResultState = {
    SolvedTasks: number
    TotalTasks: number

    GameStart: number
    GameDuration: number
    
    IsVoting: boolean
}

export type AdminState = {
    SolvedCodes: number
    IsVoting: boolean
    RequiredCodes: number
    TotalStations: number
    CodeMask: string
    CooldownDuration: number
    GameDuration: number
}

export function useApi(): Api {
    return inject(Api.provideKey)!
}

export class Api {
    playerId: string
    baseUrl: string = import.meta.env.VITE_API_BASE_URL

    static provideKey = Symbol("api-provide-key")

    constructor() {
        const storedPlayerId = localStorage.getItem("player_id")
        if(storedPlayerId) {
            this.playerId = storedPlayerId
            return
        }
        this.playerId = crypto.randomUUID()
        localStorage.setItem("player_id", this.playerId)
    }

    async getPlayerState(): Promise<PlayerState> {
        const result = await fetch(`${this.baseUrl}/api/player_state?player_id=${this.playerId}`)
        const state = await result.json()
        return state
    }

    async getStationState(stationId: number, pass: string): Promise<StationState> {
        const result = await fetch(`${this.baseUrl}/api/station_state?station_id=${stationId}`, {
            headers: {
                'X-Pass': pass
            }
        })
        const state:StationState = await result.json()
        state.CooldownUntil *= 1000
        return state
    }

    async getResultState(pass: string): Promise<ResultState> {
        const result = await fetch(`${this.baseUrl}/api/result_state`, {
            headers: {
                'X-Pass': pass
            }
        })
        const state: ResultState = await result.json()
        return state
    }

    async getAdminState(pass: string): Promise<AdminState> {
        const result = await fetch(`${this.baseUrl}/api/admin_state`, {
            headers: {
                'X-Pass': pass
            }
        })
        const state: AdminState = await result.json()
        return state
    }

    async updateSetting(setting: string, value: any, pass: string) {
        await fetch(`${this.baseUrl}/api/settings?${setting}=${value}`, {
            method: "POST",
            headers: {
                'X-Pass': pass
            }
        })
    }

    async submitCode(code: string): Promise<boolean> {
        const result = await fetch(`${this.baseUrl}/api/submit?player_id=${this.playerId}&code=${code}`, {
            method: "POST"
        })
        if(result.ok) {
            return true
        }
        throw Error(await result.text())
    }


    install(app: App) {
        app.provide(Api.provideKey, this)
    }
}