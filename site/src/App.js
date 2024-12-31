import React, { useState, useEffect, useRef } from "react";

const Game = () => {
	const [players, setPlayers] = useState([]); // List of players
	const [gameState, setGameState] = useState(null); // Game state from server
	const [playerName, setPlayerName] = useState(""); // Current player's name
	const socket = useRef(null); // WebSocket reference
	const canvasRef = useRef(null); // Canvas reference for rendering

	useEffect(() => {
		// Connect to the WebSocket server
		socket.current = new WebSocket("ws://localhost:8080/ws");

		// Handle incoming WebSocket messages
		socket.current.onmessage = (event) => {
			const message = JSON.parse(event.data);

			switch (message.type) {
				case "welcome":
					setPlayerName(message.content); // Set player's name
					break;
				case "players":
					setPlayers(message.content); // Update player list
					break;
				case "start":
					console.log("Game started!");
					setGameState({}); // Reset game state
					break;
				case "gameState":
					console.log(message);
					setGameState(message.content); // Update game state
					break;
				case "gameOver":
					console.log("Game Over");
					setGameState(null);
					break;
				default:
					console.error("Unknown message:", message);
			}
		};

		return () => socket.current.close();
	}, []);

	useEffect(() => {
		// Render the game state on the canvas
		if (gameState && canvasRef.current) {
			const canvas = canvasRef.current;
			const ctx = canvas.getContext("2d");

			// Clear canvas
			ctx.clearRect(0, 0, canvas.width, canvas.height);

			// Draw ball
			ctx.beginPath();
			ctx.arc(gameState.BallX, gameState.BallY, 10, 0, 2 * Math.PI);
			ctx.fillStyle = "red";
			ctx.fill();

			// Draw paddles
			ctx.fillStyle = "blue";
			ctx.fillRect(20, gameState.Paddle1Y, gameState.PaddleWidth, gameState.PaddleHeight);
			ctx.fillRect(
				canvas.width - 40,
				gameState.Paddle2Y,
				gameState.PaddleWidth,
				gameState.PaddleHeight
			);
		}
	}, [gameState]);

	const handleChallenge = (target) => {
		if (socket.current && target !== playerName) {
			socket.current.send(
				JSON.stringify({ type: "challenge", content: target })
			);
		}
	};

	const handleKeyDown = (e) => {
		if (socket.current) {
			if (e.key === "ArrowUp") {
				socket.current.send(JSON.stringify({ type: "move", content: "up" }));
			} else if (e.key === "ArrowDown") {
				socket.current.send(JSON.stringify({ type: "move", content: "down" }));
			}
		}
	};

	return (
		<div onKeyDown={handleKeyDown} tabIndex={0} style={{ outline: "none" }}>
			<h1>Ping Pong Game</h1>
			<p>Your Name: {playerName}</p>
			<div>
				<h2>Players</h2>
				<ul>
					{players.length === 0 ? (
						<p>No players connected</p>
					) : (
						players.filter((player) => player !== playerName).map((player) => (
							<li key={player}>
								{player}{" "}
								<button onClick={() => handleChallenge(player)}>Challenge</button>
							</li>
						))
					)}
				</ul>
			</div>
			{gameState ? (
				<canvas ref={canvasRef} width="800" height="400" style={{ border: "1px solid black" }} />
			) : (
				<p>Waiting for game to start...</p>
			)}
		</div>
	);
};

export default Game;

