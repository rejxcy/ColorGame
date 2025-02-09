// 遊戲狀態
export interface GameState {
  quiz: string;
  displayColor: string;
  progress: number;
  percentage: number;
  wrongCount: number;
  isFinished: boolean;
  totalQuiz: number;
}

// 玩家排名
export interface PlayerRank {
  id: string;
  name: string;
  score: number;
  wrongCount: number;
  duration: number;
  isFinished: boolean;
}

// 玩家資訊
export interface Player {
  id: string;
  name: string;
  isHost: boolean;
  isReady: boolean;
  score: number;
  game: GameState;
}

// WebSocket 消息類型
export enum MessageType {
  Answer = 'answer',
  Restart = 'restart',
  GameState = 'game_state',
  GameOver = 'game_over',
  Error = 'error',
  JoinRoom = 'join_room',
  LeaveRoom = 'leave_room',
  PlayerList = 'player_list',
  GameStart = 'game_start',
  GameRank = 'game_rank',
  Progress = 'progress',
  Ready = 'ready'
}

// WebSocket 消息
export interface WSMessage {
  type: MessageType;
  payload: any;
} 