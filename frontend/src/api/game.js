export class GameWebSocket {
    constructor() {
        this.ws = null;
        this.messageHandlers = new Set();
    }

    async connect(roomId, playerName, isHost = false) {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
        const params = new URLSearchParams({
            room_id: roomId,
            player_name: playerName,
            is_host: isHost
        })
        
        this.ws = new WebSocket(`${protocol}//${window.location.host}/api/game/ws?${params}`)
        
        this.ws.onmessage = (event) => {
            try {
                const message = JSON.parse(event.data);
                console.log('Received message:', message);
                
                // 處理玩家列表更新
                if (message.type === 'player_list') {
                    const handler = this.messageHandlers.get('player_list');
                    if (handler) {
                        handler(message.payload);
                    }
                }
                
                // 處理其他消息
                const handler = this.messageHandlers.get(message.type);
                if (handler) {
                    handler(message.payload);
                }
            } catch (err) {
                console.error('處理WebSocket消息失敗:', err);
            }
        };

        return new Promise((resolve, reject) => {
            this.ws.onopen = () => {
                console.log('WebSocket connected');
                resolve();
            };
            this.onerror = (error) => {
                console.error('WebSocket error:', error);
                reject(error);
            };
        });
    }

    onMessage(handler) {
        this.messageHandlers.add(handler);
    }

    handleMessage(event) {
        const data = JSON.parse(event.data);
        this.messageHandlers.forEach(handler => handler(data));
    }

    sendAnswer(color) {
        this.send('answer', color);
    }

    restart() {
        this.send('restart');
    }

    send(type, payload = {}) {
        if (this.ws && this.ws.readyState === WebSocket.OPEN) {
            const message = JSON.stringify({ type, payload });
            console.log('Sending message:', message);
            this.ws.send(message);
        } else {
            console.error('WebSocket未連接');
        }
    }

    close() {
        if (this.ws) {
            this.ws.close();
        }
    }
}