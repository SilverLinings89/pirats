import React, { useState } from "react";

const ProgrammingPanel = ({ onSubmitMoves }) => {
	const [moves, setMoves] = useState([]);

	const addMove = (move) => {
		if (moves.length < 5) setMoves([...moves, move]);
	};

	const submitMoves = () => {
		onSubmitMoves(moves);
		setMoves([]);
	};

	return (
		<div className="programming-panel">
			<h3>Plan Your Moves</h3>
			<div className="selected-moves">
				{moves.map((move, index) => (
					<span key={index}>{move}</span>
				))}
			</div>
			<div className="move-buttons">
				<button onClick={() => addMove("forward")}>Forward</button>
				<button onClick={() => addMove("left")}>Turn Left</button>
				<button onClick={() => addMove("right")}>Turn Right</button>
			</div>
			<button onClick={submitMoves} disabled={moves.length < 5}>
				Submit Moves
			</button>
		</div>
	);
};

export default ProgrammingPanel;

