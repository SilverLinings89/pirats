import React, { useState, useEffect } from "react";
import GameBoard from "./components/GameBoard";
import ProgrammingPanel from "./components/ProgrammingPanel";
import PlayerList from "./components/PlayerList";
import WebRTCChat from "./components/WebRTCChat";

const App = () => {
	const [socket, setSocket] = useState(null);
	const [gameState, setGameState] = useState(null);
	const [playerId, setPlayerId] = useState("");

	useEffect(() => {
		const ws = new WebSocket("ws://localhost:8080/ws");
		ws.onopen = () => console.log("Connected to game server");
		ws.onmessage = (message) => {
			const data = JSON.parse(message.data);
			if (data.type === "gameState") setGameState(data.state);
			if (data.type === "playerId") setPlayerId(data.playerId);
		};
		setSocket(ws);

		return () => ws.close();
	}, []);

	const submitMoves = (moves) => {
		if (socket) {
			socket.send(JSON.stringify({ type: "submitMoves", moves }));
		}
	};

	return (
		<div className="app">
			<h1>Pirate Rally</h1>
			<PlayerList players={gameState?.players || []} />
			<GameBoard board={gameState?.board} ships={gameState?.ships} />
			<ProgrammingPanel onSubmitMoves={submitMoves} />
			<WebRTCChat playerId={playerId} />
		</div>
	);
};

export default App;

