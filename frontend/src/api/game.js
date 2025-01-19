export class GameWebSocket {
    constructor() {
        this.ws = null;
        this.messageHandlers = new Map();
    }

    connect() {
        // 使用相對路徑，讓 nginx 處理代理
        const wsProtocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${wsProtocol}//${window.location.host}/api/game/ws`;
        this.ws = new WebSocket(wsUrl);
        
        this.ws.onmessage = (event) => {
            const message = JSON.parse(event.data);
            console.log('Received message:', message);
            const handler = this.messageHandlers.get(message.type);
            if (handler) {
                handler(message.payload);
            }
        };

        this.ws.onopen = () => {
            console.log('WebSocket connected');
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
        };

        this.ws.onclose = () => {
            console.log('WebSocket closed');
        };

        return new Promise((resolve, reject) => {
            this.ws.onopen = () => resolve();
            this.ws.onerror = (error) => reject(error);
        });
    }

    on(messageType, handler) {
        this.messageHandlers.set(messageType, handler);
    }

    sendAnswer(color) {
        this.send('answer', color);
    }

    restart() {
        this.send('restart');
    }

    send(type, payload) {
        if (this.ws?.readyState === WebSocket.OPEN) {
            const message = {
                type: type,
                payload: payload
            };
            console.log('Sending message:', message);
            this.ws.send(JSON.stringify(message));
        }
    }

    close() {
        this.ws?.close();
    }
}